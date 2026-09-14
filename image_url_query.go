package hcti

import (
	"fmt"
	"strconv"
)

// renderQueryWriter accepts only internal parameter names. Template variable
// collisions use the API's __ro_ prefix, including legacy crop_w/crop_h aliases.
type renderQueryWriter struct {
	b        *signedURLBuilder
	reserved map[string]any
	conflict bool
}

func (w *renderQueryWriter) key(name safeQueryKey) {
	if _, exists := w.reserved[string(name)]; exists {
		name = "__ro_" + name
		if _, exists := w.reserved[string(name)]; exists {
			w.conflict = true
		}
	}
	w.b.safeKey(name)
}

func (w *renderQueryWriter) integer(name safeQueryKey, value *int) {
	if value == nil {
		return
	}
	w.key(name)
	w.b.buf = strconv.AppendInt(w.b.buf, int64(*value), 10)
}

func (w *renderQueryWriter) origin(name safeQueryKey, value CropOrigin) {
	if value == "" || value == CropStart {
		return
	}
	w.key(name)
	w.b.buf = append(w.b.buf, value...)
}

func (w *renderQueryWriter) value(name safeQueryKey, value *CropValue) {
	if value == nil {
		return
	}
	w.key(name)
	w.b.buf = strconv.AppendInt(w.b.buf, int64(value.Value), 10)
	w.b.buf = appendQueryEscaped(w.b.buf, string(value.Unit))
}

func (w *renderQueryWriter) span(s *CropSpan, start, end, size, alias safeQueryKey) {
	if s == nil {
		return
	}
	w.value(start, s.Start)
	w.value(end, s.End)
	if s.Size == nil {
		return
	}
	_, longCollision := w.reserved[string(size)]
	_, shortCollision := w.reserved[string(alias)]
	if longCollision || !shortCollision {
		w.value(size, s.Size)
	}
	if shortCollision {
		w.value(alias, s.Size)
	}
}

func cropSpanOrigin(s *CropSpan) CropOrigin {
	if s != nil && s.Start == nil && s.Size != nil {
		return s.Origin
	}
	return CropStart
}

func (r RenderImageOptions) appendQuery(b *signedURLBuilder, reserved map[string]any) error {
	w := renderQueryWriter{b: b, reserved: reserved}
	w.integer("dpi", r.DPI)
	w.integer("height", r.Height)
	w.integer("width", r.Width)
	if c := r.Crop; c != nil {
		if c.AspectRatio != nil {
			w.key("aspect_ratio")
			b.buf = strconv.AppendInt(b.buf, int64(c.AspectRatio.Width), 10)
			b.buf = append(b.buf, '_')
			b.buf = strconv.AppendInt(b.buf, int64(c.AspectRatio.Height), 10)
		}
		x, y := cropSpanOrigin(c.Horizontal), cropSpanOrigin(c.Vertical)
		if c.AspectRatioAxis == CropWidth {
			y = c.ComputedOrigin
		}
		if c.AspectRatioAxis == CropHeight {
			x = c.ComputedOrigin
		}
		w.origin("x_origin", x)
		w.origin("y_origin", y)
		w.span(c.Horizontal, "x_1", "x_2", "crop_width", "crop_w")
		w.span(c.Vertical, "y_1", "y_2", "crop_height", "crop_h")
	}
	if w.conflict {
		return fmt.Errorf("hcti: template variable conflicts with a reserved render override")
	}
	return nil
}
