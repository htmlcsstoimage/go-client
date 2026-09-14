package hcti

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateImageSerialization(t *testing.T) {
	client := testClient(t, 200, `{"id":"image-id","url":"https://hcti.io/v1/image/image-id"}`, func(r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/image" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["google_fonts"] != "Open+Sans|Roboto" {
			t.Fatalf("fonts: %v", payload["google_fonts"])
		}
		if payload["transparent_background"] != false || payload["dedupe_duration_s"] != float64(0) {
			t.Fatal("explicit zero values were omitted")
		}
		if _, exists := payload["device_scale"]; exists {
			t.Fatal("omitted value was sent")
		}
		pdf := payload["pdf_options"].(map[string]any)
		if pdf["page_width"] != "8.5in" || pdf["print_background"] != false {
			t.Fatalf("PDF: %v", pdf)
		}
		margins := pdf["margins"].([]any)
		for i, want := range []string{"1cm", "2mm", "3px", "4in"} {
			if margins[i] != want {
				t.Fatalf("margins: %v", margins)
			}
		}
	})
	request := HTMLImageRequest{HTML: "<h1>Hello</h1>", GoogleFonts: GoogleFonts{" Open Sans ", "Roboto", "Open Sans", " "}, ImageOptions: ImageOptions{
		RenderOptions: RenderOptions{TransparentBackground: Ptr(false)}, DedupeDurationSeconds: Ptr(0),
		PDFOptions: &PDFOptions{PrintBackground: Ptr(false), PageWidth: &PDFLength{8.5, Inches}, Margins: &PDFMargins{PDFLength{1, Centimeters}, PDFLength{2, Millimeters}, PDFLength{3, Pixels}, PDFLength{4, Inches}}},
	}}
	image, err := client.CreateImage(context.Background(), request)
	if err != nil || image.ID != "image-id" {
		t.Fatalf("image: %v, error: %v", image, err)
	}
}

func TestBatch(t *testing.T) {
	c := testClient(t, 200, `{"images":[{"id":"one","url":"https://example.com/one"}]}`, func(r *http.Request) {
		if r.URL.Path != "/v1/image/batch" {
			t.Fatal(r.URL)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["default_options"].(map[string]any)["html"] != "<h1>Hello</h1>" {
			t.Fatal(body)
		}
	})
	result, err := c.CreateImageBatch(context.Background(), BatchRequest{Variations: []ImageRequest{HTMLImageRequest{CSS: Ptr("h1{color:red}")}}, DefaultOptions: HTMLImageRequest{HTML: "<h1>Hello</h1>"}})
	if err != nil || len(result.Images) != 1 {
		t.Fatalf("%v %v", result, err)
	}
	if _, err := c.CreateImageBatch(context.Background(), BatchRequest{Variations: []ImageRequest{TemplatedImageRequest{TemplateID: "tpl"}}}); err == nil {
		t.Fatal("accepted unsupported batch request")
	}
	var nilRequest *HTMLImageRequest
	if _, err := c.CreateImage(context.Background(), nilRequest); err == nil {
		t.Fatal("accepted nil request")
	}
}
