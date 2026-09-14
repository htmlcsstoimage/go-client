package hcti

import (
	"net/url"
	"strings"
	"testing"
)

func TestTemplateSigningGolden(t *testing.T) {
	c := NewClient("test-id", "test-secret")
	got, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateVersion: Ptr(int64(42)), TemplateValues: map[string]any{"headline": "Hello world"}, Format: PNG})
	want := "https://hcti.io/v1/image/tpl/8033706d06983bca46abc5edc6f7dd041c58f2e9c80bc0c5d4a3d2f2bd5b4040/png?headline=Hello+world&template_version=42"
	if err != nil || got != want {
		t.Fatalf("got %s (%v)", got, err)
	}
}

func TestURLSigningParameters(t *testing.T) {
	c := NewClient("test-id", "test-secret")
	request := URLImageRequest{URL: "https://example.com/?a=1&b=2", Headers: map[string]string{"Z": "last", "A": "first"}, AdditionalHeaderOrigins: []string{"https://cdn.example.com"}, ImageOptions: ImageOptions{Format: WebP, DedupeDurationSeconds: Ptr(42), PDFOptions: &PDFOptions{}, RenderOptions: RenderOptions{TransparentBackground: Ptr(false)}}}
	got, err := c.GenerateCreateAndRenderURL(request)
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if q.Get("url") != request.URL || q.Get("transparent_background") != "false" || q.Has("dedupe_duration_s") || q.Has("pdf_options") || q.Has("format") {
		t.Fatal(q)
	}
	if strings.Join(q["headers"], "|") != "A:first|Z:last" {
		t.Fatal(q)
	}
	again, err := c.GenerateCreateAndRenderURL(request)
	if err != nil || again != got {
		t.Fatal("signing is not deterministic")
	}
	if request.PDFOptions == nil || request.DedupeDurationSeconds == nil {
		t.Fatal("request mutated")
	}
}

func TestInvalidSigningInput(t *testing.T) {
	c := NewClient("id", "key")
	for _, request := range []TemplatedImageRequest{
		{},
		{TemplateID: "tpl", Format: "../invalid"},
		{TemplateID: "tpl", TemplateValues: map[string]any{"bad": make(chan int)}},
		{TemplateID: "tpl", TemplateValues: map[string]any{"template_version": 2}},
	} {
		if _, err := c.GenerateTemplatedImageURL(request); err == nil {
			t.Fatal("expected validation error")
		}
	}
}
