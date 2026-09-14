package hcti

import "fmt"

// RenderImageOptions transforms an image when its URL is fetched. Cropping is
// applied before resizing. These settings do not change the stored image.
type RenderImageOptions struct {
	// Format selects PNG, JPG, WebP, or PDF. Empty preserves the request's format.
	Format ImageFormat
	// DPI sets output density, strictly greater than 30 and less than 600.
	DPI *int
	// Height resizes the result to 1–5000 pixels. Nil preserves its height or aspect ratio.
	Height *int
	// Width resizes the result to 1–5000 pixels. Nil preserves its width or aspect ratio.
	Width *int
	// Crop selects a rectangle or an aspect-ratio crop before resizing.
	Crop *Crop
}

// ImageURL builds a retrieval URL for an existing image without an API request.
// Supply at most one set of render options.
func (c *Client) ImageURL(id string, options ...RenderImageOptions) (string, error) {
	path, err := resourcePath("image", id)
	if err != nil {
		return "", err
	}
	render, err := selectRenderOptions(options, "")
	if err != nil {
		return "", err
	}
	var buffer [512]byte
	b := signedURLBuilder{buf: append(buffer[:0], c.baseURL...)}
	b.buf = append(b.buf, path...)
	if render.Format != "" {
		b.buf = append(b.buf, '.')
		b.buf = append(b.buf, render.Format...)
	}
	if err := render.appendQuery(&b, nil); err != nil {
		return "", err
	}
	return string(b.buf), nil
}

func selectRenderOptions(options []RenderImageOptions, fallback ImageFormat) (RenderImageOptions, error) {
	if len(options) > 1 {
		return RenderImageOptions{}, fmt.Errorf("hcti: supply at most one set of render options")
	}
	var result RenderImageOptions
	if len(options) == 1 {
		result = options[0]
	}
	if result.Format == "" {
		result.Format = fallback
	}
	if result.Format != "" && result.Format != PNG && result.Format != JPG && result.Format != WebP && result.Format != PDF {
		return result, fmt.Errorf("hcti: unsupported image format %q", result.Format)
	}
	if result.DPI != nil && (*result.DPI <= 30 || *result.DPI >= 600) {
		return result, fmt.Errorf("hcti: DPI must be greater than 30 and less than 600")
	}
	if result.Width != nil && (*result.Width < 1 || *result.Width > 5000) {
		return result, fmt.Errorf("hcti: width must be from 1 to 5000")
	}
	if result.Height != nil && (*result.Height < 1 || *result.Height > 5000) {
		return result, fmt.Errorf("hcti: height must be from 1 to 5000")
	}
	if result.Crop != nil {
		return result, result.Crop.validate()
	}
	return result, nil
}
