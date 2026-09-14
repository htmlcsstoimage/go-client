package hcti

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestBatchExplicitEmptyOverridesAndDedupe(t *testing.T) {
	defaults := URLImageRequest{
		URL: "https://example.com", CSS: Ptr("body {color:red}"), Headers: map[string]string{"X-Test": "present"},
		AdditionalHeaderOrigins: []string{"https://cdn.example.com"}, FullScreen: Ptr(true),
		ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(60), Selector: Ptr(".card"), RenderOptions: RenderOptions{MSDelay: Ptr(10), ProxyID: Ptr("proxy")}},
	}
	empty := URLImageRequest{
		CSS: Ptr(""), Headers: map[string]string{}, AdditionalHeaderOrigins: []string{}, FullScreen: Ptr(false),
		ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(0), Selector: Ptr(""), RenderOptions: RenderOptions{
			MSDelay: Ptr(0), ProxyID: Ptr(""), StorageDestinationID: Ptr(""), Timezone: Ptr(""), ColorScheme: Ptr(ColorScheme("")), MediaType: Ptr(MediaType("")),
		}},
	}
	html := HTMLImageRequest{CSS: Ptr(""), GoogleFonts: GoogleFonts{}, ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(120)}}
	request := BatchRequest{DefaultOptions: &defaults, Variations: []ImageRequest{empty, &empty, html, &html, URLImageRequest{}, HTMLImageRequest{}}}
	c := testClient(t, 200, `{"images":[]}`, func(r *http.Request) {
		var payload struct {
			Defaults   map[string]any   `json:"default_options"`
			Variations []map[string]any `json:"variations"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Defaults["css"] != *defaults.CSS || payload.Defaults["url"] != defaults.URL {
			t.Fatal(payload.Defaults)
		}
		for _, item := range append(payload.Variations, payload.Defaults) {
			if _, present := item["dedupe_duration_s"]; present {
				t.Fatal("batch includes dedupe")
			}
		}
		for _, item := range payload.Variations[:2] {
			want := map[string]any{"css": "", "headers": map[string]any{}, "additional_header_origins": []any{}, "full_screen": false, "selector": "", "ms_delay": float64(0), "proxy_id": "", "storage_destination_id": "", "timezone": "", "color_scheme": "", "media_type": ""}
			if !reflect.DeepEqual(item, want) {
				t.Fatalf("got %#v, want %#v", item, want)
			}
		}
		for _, item := range payload.Variations[2:4] {
			if !reflect.DeepEqual(item, map[string]any{"css": "", "google_fonts": ""}) {
				t.Fatal(item)
			}
		}
		for _, item := range payload.Variations[4:] {
			if len(item) != 0 {
				t.Fatal("omitted fields unexpectedly sent", item)
			}
		}
	})
	if _, err := c.CreateImageBatch(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	if *defaults.DedupeDurationSeconds != 60 || *empty.DedupeDurationSeconds != 0 || *html.DedupeDurationSeconds != 120 {
		t.Fatal("caller requests were mutated")
	}
	// Single-image requests retain dedupe, including explicit zero.
	for _, item := range []ImageRequest{defaults, &defaults, empty, &empty, html, &html} {
		data, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]any
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		if _, present := fields["dedupe_duration_s"]; !present {
			t.Fatal("single-image dedupe was removed")
		}
	}
}

func TestBatchDedupeDefaultsValueAndPointer(t *testing.T) {
	for _, defaults := range []ImageRequest{
		HTMLImageRequest{ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(10)}},
		&HTMLImageRequest{ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(10)}},
		URLImageRequest{ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(10)}},
		&URLImageRequest{ImageOptions: ImageOptions{DedupeDurationSeconds: Ptr(10)}},
	} {
		data, err := json.Marshal(BatchRequest{DefaultOptions: defaults, Variations: []ImageRequest{defaults}})
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `{"variations":[{}],"default_options":{}}` {
			t.Fatal(string(data))
		}
	}
}
