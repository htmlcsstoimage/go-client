package hcti

import "fmt"

// CropUnit measures a crop boundary or size in pixels or percent of the image.
type CropUnit string

const (
	CropPixels  CropUnit = "px"
	CropPercent CropUnit = "%"
)

// CropValue is a crop position or size. Pixel positions may be zero; pixel sizes
// must be positive. Percentages must be from 1 to 100. URL helpers validate values.
type CropValue struct {
	Value int
	Unit  CropUnit
}

// CropOrigin anchors a size-only span or the calculated axis of an aspect crop.
// The zero value is equivalent to CropStart.
type CropOrigin string

const (
	CropStart  CropOrigin = "start"
	CropCenter CropOrigin = "center"
	CropEnd    CropOrigin = "end"
)

// CropAxis identifies the supplied axis in an aspect-ratio crop.
type CropAxis string

const (
	CropWidth  CropAxis = "width"
	CropHeight CropAxis = "height"
)

// AspectRatio supplies positive width and height components, such as 16 and 9.
type AspectRatio struct{ Width, Height int }

// CropSpan describes one axis. Set Start alone to crop through the remaining
// image; Start and End for boundaries; Start and Size for a fixed size; or Size
// and Origin to anchor a size. End and Size cannot both be set.
type CropSpan struct {
	Start *CropValue
	End   *CropValue
	Size  *CropValue
	// Origin applies only to Size without Start. Empty means CropStart.
	Origin CropOrigin
}

// Crop describes a rectangle or an aspect-ratio crop. A rectangle needs at least
// one span. An aspect crop needs exactly the span named by AspectRatioAxis; the
// other axis is calculated from AspectRatio and anchored by ComputedOrigin.
type Crop struct {
	Horizontal      *CropSpan
	Vertical        *CropSpan
	AspectRatio     *AspectRatio
	AspectRatioAxis CropAxis
	// ComputedOrigin anchors the calculated axis. Empty means CropStart.
	ComputedOrigin CropOrigin
}

func validCropOrigin(origin CropOrigin) bool {
	return origin == "" || origin == CropStart || origin == CropCenter || origin == CropEnd
}

func (v CropValue) validate(size bool) error {
	if v.Unit != CropPixels && v.Unit != CropPercent {
		return fmt.Errorf("hcti: crop unit must be px or %%")
	}
	if v.Unit == CropPercent && (v.Value < 1 || v.Value > 100) {
		return fmt.Errorf("hcti: crop percentage must be from 1 to 100")
	}
	if v.Unit == CropPixels && (v.Value < 0 || (size && v.Value == 0)) {
		return fmt.Errorf("hcti: crop pixel positions must be nonnegative and sizes positive")
	}
	return nil
}

func (s *CropSpan) validate() error {
	if s == nil {
		return nil
	}
	if !validCropOrigin(s.Origin) {
		return fmt.Errorf("hcti: invalid crop origin")
	}
	if s.Start == nil && s.Size == nil {
		return fmt.Errorf("hcti: crop span needs a start or size")
	}
	if s.End != nil && (s.Start == nil || s.Size != nil) {
		return fmt.Errorf("hcti: crop end requires a start and cannot be combined with size")
	}
	if s.Start != nil {
		if err := s.Start.validate(false); err != nil {
			return err
		}
		if s.Origin != "" && s.Origin != CropStart {
			return fmt.Errorf("hcti: crop origin requires size without start")
		}
	}
	if s.End != nil {
		if err := s.End.validate(false); err != nil {
			return err
		}
		if s.Start.Unit == s.End.Unit && s.End.Value <= s.Start.Value {
			return fmt.Errorf("hcti: crop end must be greater than start")
		}
	}
	if s.Size != nil {
		return s.Size.validate(true)
	}
	return nil
}

func (c *Crop) validate() error {
	if !validCropOrigin(c.ComputedOrigin) {
		return fmt.Errorf("hcti: invalid computed crop origin")
	}
	if c.AspectRatio == nil {
		if c.Horizontal == nil && c.Vertical == nil {
			return fmt.Errorf("hcti: crop requires at least one axis")
		}
		if c.AspectRatioAxis != "" || (c.ComputedOrigin != "" && c.ComputedOrigin != CropStart) {
			return fmt.Errorf("hcti: rectangle cannot set aspect-ratio axis or computed origin")
		}
	} else {
		if c.AspectRatio.Width <= 0 || c.AspectRatio.Height <= 0 {
			return fmt.Errorf("hcti: aspect ratio components must be positive")
		}
		switch c.AspectRatioAxis {
		case CropWidth:
			if c.Horizontal == nil || c.Vertical != nil {
				return fmt.Errorf("hcti: width-based aspect crop requires only a horizontal span")
			}
		case CropHeight:
			if c.Vertical == nil || c.Horizontal != nil {
				return fmt.Errorf("hcti: height-based aspect crop requires only a vertical span")
			}
		default:
			return fmt.Errorf("hcti: aspect crop requires a width or height axis")
		}
	}
	if err := c.Horizontal.validate(); err != nil {
		return err
	}
	return c.Vertical.validate()
}
