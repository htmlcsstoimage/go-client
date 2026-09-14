package hcti

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
)

func BenchmarkGenerateCreateAndRenderURL(b *testing.B) {
	client := NewClient("test-id", "test-secret")
	for _, tc := range []struct {
		name    string
		request URLImageRequest
	}{
		{"minimal", URLImageRequest{URL: "https://example.com"}},
		{"configured", URLImageRequest{
			URL: "https://example.com/?a=1&b=2", CSS: Ptr("body { background: white; }"),
			Headers:                 map[string]string{"Accept-Language": "en-US", "X-Custom": "example"},
			AdditionalHeaderOrigins: []string{"https://cdn.example.com"}, FullScreen: Ptr(true),
			ImageOptions: ImageOptions{Format: PNG, RenderOptions: RenderOptions{
				ViewportWidth: Ptr(1280), ViewportHeight: Ptr(720), DeviceScale: Ptr(1.5),
				TransparentBackground: Ptr(false), MSDelay: Ptr(0), ColorScheme: Ptr(Light),
			}},
		}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := client.GenerateCreateAndRenderURL(tc.request); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGenerateTemplatedImageURL(b *testing.B) {
	client := NewClient("test-id", "test-secret")
	for _, tc := range []struct {
		name    string
		request TemplatedImageRequest
	}{
		{"minimal", TemplatedImageRequest{TemplateID: "tpl", TemplateValues: map[string]any{"title": "Hello world"}}},
		{"scalars", TemplatedImageRequest{TemplateID: "tpl", TemplateVersion: Ptr(int64(42)), Format: PNG, TemplateValues: map[string]any{"title": "Hello 👩🏽‍💻", "count": 42, "enabled": false, "price": 19.95, "description": "A card with <html> & special characters"}}},
		{"structured", TemplatedImageRequest{TemplateID: "tpl", TemplateVersion: Ptr(int64(42)), TemplateValues: map[string]any{"title": "Hello", "items": []any{map[string]any{"label": "One", "count": 1}, map[string]any{"label": "Two", "count": 2}}}}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for _, impl := range []struct {
				name     string
				generate func(TemplatedImageRequest) (string, error)
			}{
				{"buffer", func(r TemplatedImageRequest) (string, error) { return client.GenerateTemplatedImageURL(r) }},
				{"previous", func(r TemplatedImageRequest) (string, error) { return legacyTemplateURL(client, r) }},
			} {
				b.Run(impl.name, func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						if _, err := impl.generate(tc.request); err != nil {
							b.Fatal(err)
						}
					}
				})
			}
		})
	}
}

// Retain the previous implementation only as a benchmark/compatibility baseline.
func legacyTemplateURL(c *Client, r TemplatedImageRequest) (string, error) {
	q := url.Values{}
	if r.TemplateVersion != nil {
		q.Set("template_version", strconv.FormatInt(*r.TemplateVersion, 10))
	}
	for key, value := range r.TemplateValues {
		if value == nil {
			continue
		}
		data, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		text := string(data)
		if len(data) > 0 && data[0] == '"' {
			if err := json.Unmarshal(data, &text); err != nil {
				return "", err
			}
		}
		q.Set(key, text)
	}
	encoded := q.Encode()
	mac := hmac.New(sha256.New, []byte(c.apiKey))
	_, _ = mac.Write([]byte(encoded))
	result := c.baseURL + "/v1/image/" + url.PathEscape(r.TemplateID) + "/" + hex.EncodeToString(mac.Sum(nil))
	if r.Format != "" {
		result += "/" + string(r.Format)
	}
	if encoded != "" {
		result += "?" + encoded
	}
	return result, nil
}
