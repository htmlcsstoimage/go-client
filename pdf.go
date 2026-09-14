package hcti

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// PDFLength expresses a PDF dimension. The zero value means zero pixels.
type PDFLength struct {
	// Value is the dimension magnitude. It must be finite and nonnegative.
	Value float64
	// Unit selects Pixels, Inches, Centimeters, or Millimeters.
	// An empty unit means Pixels.
	Unit PDFUnit
}

// PDFUnit is a unit accepted for PDF dimensions and margins.
type PDFUnit string

const (
	Pixels      PDFUnit = "px"
	Inches      PDFUnit = "in"
	Centimeters PDFUnit = "cm"
	Millimeters PDFUnit = "mm"
)

func (length PDFLength) MarshalJSON() ([]byte, error) {
	unit := length.Unit
	if unit == "" {
		unit = Pixels
	}
	if unit != Pixels && unit != Inches && unit != Centimeters && unit != Millimeters {
		return nil, fmt.Errorf("hcti: unsupported PDF unit %q", unit)
	}
	if math.IsNaN(length.Value) || math.IsInf(length.Value, 0) || length.Value < 0 {
		return nil, fmt.Errorf("hcti: PDF dimensions must be finite and nonnegative")
	}
	return json.Marshal(strconv.FormatFloat(length.Value, 'f', -1, 64) + string(unit))
}

// PDFMargins serializes in CSS order: top, right, bottom, left.
type PDFMargins struct {
	// Top sets the top margin. The zero value is zero pixels.
	Top PDFLength
	// Right sets the right margin. The zero value is zero pixels.
	Right PDFLength
	// Bottom sets the bottom margin. The zero value is zero pixels.
	Bottom PDFLength
	// Left sets the left margin. The zero value is zero pixels.
	Left PDFLength
}

func (m PDFMargins) MarshalJSON() ([]byte, error) {
	return json.Marshal([4]PDFLength{m.Top, m.Right, m.Bottom, m.Left})
}

// PDFOptions controls PDF page layout and printing.
type PDFOptions struct {
	// PrintBackground includes background graphics in PDF output.
	// Nil uses the API default; Ptr(false) explicitly excludes them.
	PrintBackground *bool `json:"print_background,omitempty"`
	// Scale sets the PDF output scale factor. Nil uses the API default.
	Scale *float64 `json:"scale,omitempty"`
	// Margins sets the top, right, bottom, and left PDF margins.
	// Nil uses the API defaults; a non-nil zero value sends four zero-pixel margins.
	Margins *PDFMargins `json:"margins,omitempty"`
	// PageHeight sets the PDF page height with explicit units.
	// Nil uses the API default.
	PageHeight *PDFLength `json:"page_height,omitempty"`
	// PageWidth sets the PDF page width with explicit units.
	// Nil uses the API default.
	PageWidth *PDFLength `json:"page_width,omitempty"`
}
