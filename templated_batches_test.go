package hcti

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestTemplatedBatch(t *testing.T) {
	defaults := &TemplatedBatchImageOptions{TemplateID: "t-card", TemplateVersion: Ptr(int64(3)), Format: WebP, TemplateValues: map[string]any{"brand": map[string]any{"name": "Acme", "color": "red"}}}
	variation := TemplatedBatchImageOptions{TemplateID: "t-other", TemplateValues: map[string]any{"brand": map[string]any{"color": nil}, "tags": []any{}, "active": false}}
	client := testClient(t, 200, `{"images":[{"id":"two","url":"https://hcti.io/v1/image/two"},{"id":"one","url":"https://hcti.io/v1/image/one"}]}`, func(r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/image/batch/templated" {
			t.Fatal(r.URL)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		want := map[string]any{
			"default_options": map[string]any{"template_id": "t-card", "template_version": float64(3), "format": "webp", "template_values": map[string]any{"brand": map[string]any{"name": "Acme", "color": "red"}}},
			"variations":      []any{map[string]any{}, map[string]any{"template_id": "t-other", "template_values": map[string]any{"brand": map[string]any{"color": nil}, "tags": []any{}, "active": false}}},
		}
		if !reflect.DeepEqual(body, want) {
			t.Fatalf("body: %#v", body)
		}
	})
	result, err := client.CreateTemplatedImageBatch(context.Background(), &TemplatedBatchRequest{DefaultOptions: defaults, Variations: []TemplatedBatchImageOptions{{}, variation}})
	if err != nil || len(result.Images) != 2 || result.Images[0].ID != "two" || result.Images[1].ID != "one" {
		t.Fatalf("result: %v, error: %v", result, err)
	}
	if defaults.TemplateValues["brand"].(map[string]any)["color"] != "red" || variation.TemplateVersion != nil {
		t.Fatal("mutated input")
	}
}

func TestTemplatedBatchEmptyAndError(t *testing.T) {
	calls := 0
	client := testClient(t, 400, `{"error":"Bad Request","message":"Invalid template"}`, func(r *http.Request) {
		calls++
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if _, found := body["default_options"]; found {
			t.Fatal("unexpected defaults")
		}
	})
	result, err := client.CreateTemplatedImageBatch(context.Background(), &TemplatedBatchRequest{})
	if err != nil || len(result.Images) != 0 || calls != 0 {
		t.Fatal(result, err, calls)
	}
	if _, err := client.CreateTemplatedImageBatch(context.Background(), nil); err == nil {
		t.Fatal("accepted nil request")
	}
	if _, err := client.CreateTemplatedImageBatch(context.Background(), &TemplatedBatchRequest{Variations: []TemplatedBatchImageOptions{{TemplateID: "t-missing"}}}); err == nil || calls != 1 {
		t.Fatal(err, calls)
	}
}
