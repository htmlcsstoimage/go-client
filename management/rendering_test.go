package management

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefinitionWireContracts(t *testing.T) {
	for _, kind := range []string{"html_css", "url", "templated", "template"} {
		t.Run(kind, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if r.Method != "POST" {
					t.Error(r.Method)
				}
				if kind == "templated" {
					if r.URL.Path != "/v1/image/t-test/9007199254740993" {
						t.Error(r.URL.Path)
					}
					if _, ok := body["template_id"]; ok {
						t.Error("template id in body")
					}
					if _, ok := body["template_version"]; ok {
						t.Error("version in body")
					}
					if string(body["template_values"]) != `{"n":9007199254740993,"emoji":"😀","nested":[false,null]}` {
						t.Error(string(body["template_values"]))
					}
				} else {
					for _, k := range []string{"css", "device_scale", "viewport_width", "viewport_height", "selector", "proxy_id", "storage_destination_id", "jumbo_max_width", "jumbo_max_height", "transparent_background"} {
						if string(body[k]) != "null" {
							t.Errorf("%s must be explicit null: %s", k, body[k])
						}
					}
					if string(body["ms_delay"]) != "0" || string(body["render_when_ready"]) != "false" {
						t.Error("lost explicit zero/false")
					}
				}
				if kind != "template" {
					if string(body["dedupe_duration_s"]) != "0" {
						t.Error("dedupe not disabled")
					}
					fmt.Fprint(w, `{"id":"image-1","url":"https://hcti.io/v1/image/image-1"}`)
				} else {
					if _, ok := body["dedupe_duration_s"]; ok {
						t.Error("dedupe on template")
					}
					fmt.Fprint(w, `{"template_id":"t-test","template_version":"9007199254740993"}`)
				}
			}))
			defer server.Close()
			c := NewClient("id", "key", WithBaseURL(server.URL))
			id := "t-test"
			version := int64(9007199254740993)
			zero := int64(0)
			off := false
			request := &RenderDefinition{TemplateID: &id, TemplateVersion: &version, MSDelay: &zero, RenderWhenReady: &off, TemplateValues: json.RawMessage(`{"n":9007199254740993,"emoji":"😀","nested":[false,null]}`)}
			var err error
			if kind == "template" {
				var v *CreatedRender
				v, err = c.SaveTemplateDefinition(context.Background(), "", request)
				if err == nil && v.TemplateVersion != version {
					t.Error("version lost precision")
				}
			} else {
				_, err = c.CreateImageDefinition(context.Background(), kind, request)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestSavedRenderNumericStrings(t *testing.T) {
	var v SavedRender
	err := json.Unmarshal([]byte(`{"id":"image-1","image_type":"templated","template_id":"t-test","template_version":"9007199254740993","version":"9007199254740995","device_scale":"1.5","max_wait_ms":"500","og_config_content_version":"9007199254740997","template_values":{"n":9007199254740999},"pdf_options":{"scale":"0.5"}}`), &v)
	if err != nil {
		t.Fatal(err)
	}
	if v.ID != "image-1" || *v.TemplateVersion != 9007199254740993 || v.Version != 9007199254740995 || *v.DeviceScale != 1.5 || *v.MaxWaitMS != 500 || *v.OGConfigContentVersion != 9007199254740997 || *v.PDFOptions.Scale != 0.5 || string(v.TemplateValues) != `{"n":9007199254740999}` {
		t.Fatalf("incorrect metadata: %+v", v)
	}
	for _, s := range []string{`{"version":"9223372036854775808"}`, `{"template_version":1.2}`, `{"device_scale":"NaN"}`} {
		if err := json.Unmarshal([]byte(s), &v); err == nil {
			t.Errorf("accepted invalid numeric data: %s", s)
		}
	}
}
func TestExactTemplateVersionPagination(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			if r.URL.Query().Get("max_version") != "" {
				t.Error("must start at latest")
			}
			fmt.Fprint(w, `{"data":[{"id":"t-test","version":"9007199254740995","html":"{{x}}","template_type":"html_css","created_at":"now"}],"pagination":{"next_page_start":"9007199254740995"}}`)
		} else {
			if r.URL.Query().Get("max_version") != "9007199254740995" {
				t.Error("wrong exclusive cursor")
			}
			fmt.Fprint(w, `{"data":[{"id":"t-test","version":"9007199254740993","html":"{{x}}","template_type":"html_css","created_at":"now"}],"pagination":{"next_page_start":null}}`)
		}
	}))
	defer server.Close()
	c := NewClient("id", "key", WithBaseURL(server.URL))
	version := int64(9007199254740993)
	v, err := c.GetTemplateDefinition(context.Background(), "t-test", &version)
	if err != nil || v.Version != version || calls != 2 {
		t.Fatalf("lookup failed: %v, calls=%d", err, calls)
	}
}
func TestImageMetadataRejectsMalformedSuccess(t *testing.T) {
	for _, body := range []string{`{}`, `{"id":"different","image_type":"url","created_at":"now","url":"https://example.com"}`, `{"id":"image-1","image_type":"url","created_at":"now"}`, `{"id":"image-1","image_type":"templated","created_at":"now","template_id":"t-test","template_version":1,"template_values":null}`, `{"id":"image-1","image_type":"templated","created_at":"now","template_id":"t-test","template_version":1,"template_values":[]}`} {
		t.Run(body, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/images/image-1" {
					t.Error("wrong metadata endpoint")
				}
				fmt.Fprint(w, body)
			}))
			defer s.Close()
			_, err := NewClient("id", "key", WithBaseURL(s.URL)).GetImageMetadata(context.Background(), "image-1")
			if err == nil {
				t.Fatal("accepted malformed metadata")
			}
			if strings.Contains(err.Error(), body) {
				t.Error("response leaked in error")
			}
		})
	}
}
