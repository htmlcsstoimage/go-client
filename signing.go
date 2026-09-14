package hcti

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"unicode/utf8"
)

// GenerateTemplatedImageURL signs template values without making an API request.
// Anyone with the resulting URL can render it. Do not put confidential values in URLs.
// Optional render options crop or resize the result and override the request format.
func (c *Client) GenerateTemplatedImageURL(request TemplatedImageRequest, options ...RenderImageOptions) (string, error) {
	path, err := resourcePath("image", request.TemplateID)
	if err != nil {
		return "", err
	}
	render, err := selectRenderOptions(options, request.Format)
	if err != nil {
		return "", err
	}
	var buffer [512]byte
	builder, err := newSignedURLBuilder(buffer[:0], c.baseURL, path, render.Format)
	if err != nil {
		return "", err
	}
	var keyBuffer [16]string
	keys := keyBuffer[:0]
	for key, value := range request.TemplateValues {
		if key == "template_version" {
			return "", fmt.Errorf("hcti: template_version is a reserved query parameter")
		}
		if value != nil {
			keys = append(keys, key)
		}
	}
	if request.TemplateVersion != nil {
		keys = append(keys, "template_version")
	}
	slices.Sort(keys)
	for _, key := range keys {
		if key == "template_version" {
			builder.safeKey("template_version")
			builder.buf = strconv.AppendInt(builder.buf, *request.TemplateVersion, 10)
		} else {
			builder.key(key)
			if err := builder.appendTemplateValue(request.TemplateValues[key]); err != nil {
				return "", err
			}
		}
	}
	if err := render.appendQuery(&builder, request.TemplateValues); err != nil {
		return "", err
	}
	return builder.finish(c.apiKey), nil
}

// GenerateCreateAndRenderURL signs a URL screenshot request without calling the API.
// PDF options and dedupe duration are not supported by this signed-URL helper.
// Optional render options crop or resize the result and override the request format.
func (c *Client) GenerateCreateAndRenderURL(request URLImageRequest, options ...RenderImageOptions) (string, error) {
	if request.URL == "" {
		return "", fmt.Errorf("hcti: URL is required")
	}
	render, err := selectRenderOptions(options, request.Format)
	if err != nil {
		return "", err
	}
	var buffer [1024]byte
	builder, err := newSignedURLBuilder(buffer[:0], c.baseURL, "/v1/image/create-and-render/"+url.PathEscape(c.apiID), render.Format)
	if err != nil {
		return "", err
	}
	if err := request.appendRenderQuery(&builder); err != nil {
		return "", err
	}
	if err := render.appendQuery(&builder, nil); err != nil {
		return "", err
	}
	return builder.finish(c.apiKey), nil
}

// appendTemplateValue avoids JSON serialization for ordinary scalar values.
// Structured values and custom marshalers retain encoding/json semantics.
func (b *signedURLBuilder) appendTemplateValue(value any) error {
	switch v := value.(type) {
	case string:
		b.appendTemplateString(v)
		return nil
	case bool:
		b.buf = strconv.AppendBool(b.buf, v)
		return nil
	case int, int8, int16, int32, int64:
		b.buf = strconv.AppendInt(b.buf, reflect.ValueOf(v).Int(), 10)
		return nil
	case uint, uint8, uint16, uint32, uint64, uintptr:
		b.buf = strconv.AppendUint(b.buf, reflect.ValueOf(v).Uint(), 10)
		return nil
	case float32:
		return b.appendJSONFloat(float64(v), 32)
	case float64:
		return b.appendJSONFloat(v, 64)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("hcti: encode template value: %w", err)
	}
	if len(data) > 0 && data[0] == '"' {
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		b.buf = appendQueryEscaped(b.buf, text)
	} else {
		b.buf = appendQueryEscaped(b.buf, string(data))
	}
	return nil
}

func (b *signedURLBuilder) appendTemplateString(value string) {
	// JSON replaces each malformed UTF-8 byte with U+FFFD; preserve that behavior
	// without marshaling and decoding ordinary strings.
	if utf8.ValidString(value) {
		b.buf = appendQueryEscaped(b.buf, value)
		return
	}
	for len(value) != 0 {
		r, size := utf8.DecodeRuneInString(value)
		if r == utf8.RuneError && size == 1 {
			b.buf = append(b.buf, "%EF%BF%BD"...)
		} else {
			b.buf = appendQueryEscaped(b.buf, value[:size])
		}
		value = value[size:]
	}
}
