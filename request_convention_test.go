package hcti

import (
	"context"
	"net/http"
	"testing"
)

func TestImageRequestRequiresPointers(t *testing.T) {
	for _, value := range []any{HTMLImageRequest{}, URLImageRequest{}, TemplatedImageRequest{}} {
		if _, ok := value.(ImageRequest); ok {
			t.Fatalf("%T unexpectedly implements ImageRequest", value)
		}
	}
}

func TestNilRequestPayloads(t *testing.T) {
	c := testClient(t, 200, "", func(*http.Request) { t.Fatal("nil request reached transport") })
	ctx := context.Background()
	checks := map[string]func() error{
		"batch":            func() error { _, e := c.CreateImageBatch(ctx, nil); return e },
		"template":         func() error { _, e := c.CreateTemplate(ctx, nil); return e },
		"template version": func() error { _, e := c.CreateTemplateVersion(ctx, "t-example", nil); return e },
		"signed template":  func() error { _, e := c.GenerateTemplatedImageURL(nil); return e },
		"signed URL":       func() error { _, e := c.GenerateCreateAndRenderURL(nil); return e },
	}
	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			if check() == nil {
				t.Fatal("accepted nil request")
			}
		})
	}
	for _, input := range []ImageRequest{nil, (*HTMLImageRequest)(nil), (*URLImageRequest)(nil)} {
		if _, err := c.CreateImageBatch(ctx, &BatchRequest{Variations: []ImageRequest{input}}); err == nil {
			t.Fatalf("accepted nil variation %T", input)
		}
		if input != nil {
			if _, err := c.CreateImageBatch(ctx, &BatchRequest{Variations: []ImageRequest{&HTMLImageRequest{HTML: "test"}}, DefaultOptions: input}); err == nil {
				t.Fatalf("accepted nil defaults %T", input)
			}
		}
	}
}
