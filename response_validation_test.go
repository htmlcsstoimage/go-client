package hcti

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestRejectMalformedSuccessResponses(t *testing.T) {
	calls := []struct {
		name    string
		call    func(*Client) error
		invalid []string
		valid   string
	}{
		{"image", func(c *Client) error {
			_, err := c.CreateImage(context.Background(), HTMLImageRequest{HTML: "hello"})
			return err
		},
			[]string{`{}`, `null`, `[]`, `{"id":"image"}`, `{"url":"https://hcti.io/image"}`, `{"id":"image","url":"/relative"}`, `{"id":"image","url":"javascript:secret"}`, `{"id":"image","url":12}`, `{"id":" ","url":"https://hcti.io/image"}`},
			`{"id":"image","url":"https://hcti.io/image","future_field":true}`},
		{"batch", func(c *Client) error {
			_, err := c.CreateImageBatch(context.Background(), BatchRequest{Variations: []ImageRequest{HTMLImageRequest{HTML: "hello"}}})
			return err
		},
			[]string{`{}`, `null`, `{"images":null}`, `{"images":{}}`, `{"images":[null]}`, `{"images":[{}]}`, `{"images":[{"id":"ok","url":"https://hcti.io/ok"},{"id":"bad"}]}`},
			`{"images":[{"id":"ok","url":"https://hcti.io/ok"}]}`},
		{"template", func(c *Client) error {
			_, err := c.CreateTemplate(context.Background(), TemplateRequest{HTML: "{{title}}"})
			return err
		},
			[]string{`{}`, `null`, `{"template_id":"tpl"}`, `{"template_version":42}`, `{"template_id":"tpl","template_version":0}`, `{"template_id":"tpl","template_version":-1}`},
			`{"template_id":"tpl","template_version":42}`},
		{"version", func(c *Client) error {
			_, err := c.CreateTemplateVersion(context.Background(), "tpl", TemplateRequest{HTML: "{{title}}"})
			return err
		},
			[]string{`{}`, `null`, `{"template_id":"tpl","template_version":"42"}`},
			`{"template_id":"tpl","template_version":42}`},
		{"list", func(c *Client) error {
			_, err := c.ListTemplates(context.Background(), TemplateListOptions{})
			return err
		},
			[]string{`{}`, `null`, `{"data":null}`, `{"data":[null]}`, `{"data":[{}]}`, `{"data":[{"id":"tpl"}]}`},
			`{"data":[],"pagination":{"next_page_start":null}}`},
	}
	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) {
			for _, body := range append(tc.invalid, "", "secret invalid JSON", tc.valid+` {"secret":true}`, tc.valid+" trailing-secret") {
				c := testClient(t, http.StatusOK, body, nil)
				err := tc.call(c)
				var responseErr *ResponseError
				if !errors.As(err, &responseErr) || responseErr.StatusCode != 200 {
					t.Fatalf("body %q: expected ResponseError, got %v", body, err)
				}
				if strings.Contains(err.Error(), "secret") {
					t.Fatal("response body leaked into error")
				}
			}
			if err := tc.call(testClient(t, http.StatusNoContent, "", nil)); err == nil {
				t.Fatal("accepted unexpected 204")
			}
			for _, status := range []int{200, 201} {
				if err := tc.call(testClient(t, status, tc.valid+"\n ", nil)); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestDeleteAcceptsNoContent(t *testing.T) {
	c := testClient(t, http.StatusNoContent, "", nil)
	if err := c.DeleteImage(context.Background(), "image"); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteTemplate(context.Background(), "template"); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteImageBatch(context.Background(), []string{"image"}); err != nil {
		t.Fatal(err)
	}
}
