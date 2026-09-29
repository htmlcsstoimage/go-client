package hcti

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRequestOverridesSerializeAsStringsAndClearBatchDefaults(t *testing.T) {
	rules := []RequestOverride{{Action: RequestOverrideBlock, URL: Ptr("*.js"), ResourceTypes: []RequestOverrideResourceType{ResourceScript, ResourceFetch}}}
	for _, request := range []ImageRequest{
		&HTMLImageRequest{HTML: "<p>test</p>", ImageOptions: ImageOptions{RenderOptions: RenderOptions{RequestOverrides: rules}}},
		&URLImageRequest{URL: "https://example.com", ImageOptions: ImageOptions{RenderOptions: RenderOptions{RequestOverrides: rules}}},
	} {
		data, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), `"request_overrides":[{"action":"block","url":"*.js","resource_types":["script","fetch"]}]`) {
			t.Fatalf("incorrect override JSON: %s", data)
		}
	}
	data, err := json.Marshal(&URLImageRequest{ImageOptions: ImageOptions{RenderOptions: RenderOptions{RequestOverrides: []RequestOverride{}}}})
	if err != nil || !strings.Contains(string(data), `"request_overrides":[]`) {
		t.Fatalf("empty overrides must clear batch defaults: %s, %v", data, err)
	}
	client := NewClient("id", "key")
	signed, err := client.GenerateCreateAndRenderURL(&URLImageRequest{URL: "https://example.com", ImageOptions: ImageOptions{RenderOptions: RenderOptions{RequestOverrides: rules}}})
	if err != nil || strings.Contains(signed, "request_overrides") {
		t.Fatalf("signed URL included POST-only overrides: %s, %v", signed, err)
	}
}
