package hcti

import (
	"context"
	"net/http"
	"testing"
)

func TestTemplateLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name, method, path, response string
		call                         func(*Client) error
	}{
		{"create", "POST", "/v1/template", `{"template_id":"tpl","template_version":42}`, func(c *Client) error {
			v, e := c.CreateTemplate(context.Background(), TemplateRequest{HTML: "{{title}}"})
			if e == nil && v.TemplateVersion != 42 {
				t.Fatal(v)
			}
			return e
		}},
		{"version", "POST", "/v1/template/tpl", `{"template_id":"tpl","template_version":43}`, func(c *Client) error {
			_, e := c.CreateTemplateVersion(context.Background(), "tpl", TemplateRequest{HTML: "{{title}}"})
			return e
		}},
		{"delete", "DELETE", "/v1/template/tpl", "", func(c *Client) error { return c.DeleteTemplate(context.Background(), "tpl") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := testClient(t, 200, tc.response, func(r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path {
					t.Fatalf("%s %s", r.Method, r.URL)
				}
			})
			if err := tc.call(c); err != nil {
				t.Fatal(err)
			}
		})
	}
	c := testClient(t, 200, `{"data":[{"id":"tpl","version":42,"template_type":"html_css","google_fonts":"Open+Sans|Roboto"}],"pagination":{"next_page_start":41}}`, func(r *http.Request) {
		if r.URL.Path != "/v1/template/tpl" || r.URL.RawQuery != "count=5&max_version=99" {
			t.Fatal(r.URL)
		}
	})
	page, err := c.ListTemplateVersions(context.Background(), "tpl", TemplateListOptions{Count: 5, MaxVersion: Ptr(int64(99))})
	if err != nil {
		t.Fatal(err)
	}
	if *page.Pagination.NextPageStart != 41 || page.Data[0].GoogleFonts[0] != "Open Sans" || len(page.Data[0].Raw) == 0 {
		t.Fatal(page)
	}
}
