package hcti

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

// Recompute from the exact outgoing query with the standard library, independently
// of the builder, and check that changing that query invalidates the token.
func assertSignedQuery(t *testing.T, signed, secret string, format ImageFormat, wantQuery string) {
	t.Helper()
	u, err := url.Parse(signed)
	if err != nil {
		t.Fatal(err)
	}
	if u.RawQuery != wantQuery {
		t.Fatal("signed query bytes differ")
	}
	path := u.Path
	if format != "" {
		if !strings.HasSuffix(path, "/"+string(format)) {
			t.Fatal("missing format")
		}
		path = strings.TrimSuffix(path, "/"+string(format))
	}
	token := path[strings.LastIndexByte(path, '/')+1:]
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(u.RawQuery))
	if token != hex.EncodeToString(mac.Sum(nil)) {
		t.Fatal("signature does not match outgoing bytes")
	}
	mac.Reset()
	_, _ = mac.Write([]byte(u.RawQuery + "&tampered=true"))
	if token == hex.EncodeToString(mac.Sum(nil)) {
		t.Fatal("tampered query retained signature")
	}
}

func TestScreenshotSigningUnicodeAndInvalidBytes(t *testing.T) {
	const secret = "secret-é-漢-😀"
	c := NewClient("test-id", secret)
	for _, tc := range encodingCases {
		for _, format := range []ImageFormat{"", PNG, JPG, WebP, PDF} {
			t.Run(tc.name+"/"+string(format), func(t *testing.T) {
				request := URLImageRequest{
					URL: "https://example.com/?value=" + tc.value, CSS: Ptr("/*" + tc.value + "*/"),
					Headers:      map[string]string{"X-Test": tc.value},
					ImageOptions: ImageOptions{Format: format},
				}
				signed, err := c.GenerateCreateAndRenderURL(request)
				if err != nil {
					t.Fatal(err)
				}
				want := url.Values{"url": {request.URL}, "css": {*request.CSS}, "headers": {"X-Test:" + tc.value}}
				assertSignedQuery(t, signed, secret, format, want.Encode())
				u, _ := url.Parse(signed)
				decoded, err := url.ParseQuery(u.RawQuery)
				if err != nil || decoded.Get("url") != request.URL || decoded.Get("css") != *request.CSS {
					t.Fatal("request values did not round trip")
				}
			})
		}
	}
}

func TestTemplateSigningUnicodeValues(t *testing.T) {
	c := NewClient("test-id", "secret")
	for _, tc := range encodingCases {
		t.Run(tc.name, func(t *testing.T) {
			// Template values use JSON semantics: malformed UTF-8 in string values
			// is replaced with U+FFFD by encoding/json. Raw query keys preserve bytes.
			encoded, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatal(err)
			}
			var normalized string
			if err := json.Unmarshal(encoded, &normalized); err != nil {
				t.Fatal(err)
			}
			signed, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{
				TemplateID: "tpl", TemplateVersion: Ptr(int64(42)), Format: PNG,
				TemplateValues: map[string]any{tc.value: tc.value},
			})
			if err != nil {
				t.Fatal(err)
			}
			want := url.Values{tc.value: {normalized}, "template_version": {"42"}}
			assertSignedQuery(t, signed, "secret", PNG, want.Encode())
		})
	}
}

func TestTemplateSigningStructuredValues(t *testing.T) {
	c := NewClient("test-id", "secret")
	values := map[string]any{
		"emoji": "👩🏽‍💻", "count": 42, "enabled": false, "omitted": nil,
		"nested": map[string]any{"title": "漢 & café", "items": []any{"😀", "\xff", 0, false}},
	}
	signed, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: values})
	if err != nil {
		t.Fatal(err)
	}
	nested, err := json.Marshal(values["nested"])
	if err != nil {
		t.Fatal(err)
	}
	want := url.Values{"emoji": {"👩🏽‍💻"}, "count": {"42"}, "enabled": {"false"}, "nested": {string(nested)}}
	assertSignedQuery(t, signed, "secret", "", want.Encode())
}

func FuzzScreenshotSigning(f *testing.F) {
	for _, tc := range encodingCases {
		f.Add(tc.value)
	}
	c := NewClient("test-id", "secret")
	f.Fuzz(func(t *testing.T, value string) {
		request := URLImageRequest{URL: "https://example.com/?q=" + value, CSS: Ptr(value)}
		signed, err := c.GenerateCreateAndRenderURL(request)
		if err != nil {
			t.Fatal(err)
		}
		want := url.Values{"url": {request.URL}}
		want.Set("css", value)
		assertSignedQuery(t, signed, "secret", "", want.Encode())
	})
}
