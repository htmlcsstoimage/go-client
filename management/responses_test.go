package management_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	m "github.com/htmlcsstoimage/go-client/management"
)

func TestMalformedResourceResponses(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name, body string
		call       func(*m.Client) error
	}{
		{"key missing permissions", `{"id":"id","api_id":"auth","name":"key"}`, func(c *m.Client) error { _, e := c.GetAPIKey(ctx, "id"); return e }},
		{"create key missing secret", keyJSON, func(c *m.Client) error { _, e := c.CreateAPIKey(ctx, &m.APIKeyRequest{}); return e }},
		{"page missing data", `{"pagination":{}}`, func(c *m.Client) error { _, e := c.ListProxies(ctx, m.ListOptions{}); return e }},
		{"page null item", `{"data":[null]}`, func(c *m.Client) error { _, e := c.ListProxies(ctx, m.ListOptions{}); return e }},
		{"page malformed item", `{"data":[{"name":"proxy"}]}`, func(c *m.Client) error { _, e := c.ListProxies(ctx, m.ListOptions{}); return e }},
		{"unknown storage", strings.Replace(storageJSON, "aws_s3", "unrecognized", 1), func(c *m.Client) error { _, e := c.GetStorageDestination(ctx, "id"); return e }},
		{"storage missing role", strings.Replace(storageJSON, `"role_arn":"arn:aws:iam::123:role/hcti"`, `"role_arn":null`, 1), func(c *m.Client) error { _, e := c.GetStorageDestination(ctx, "id"); return e }},
		{"unknown OG", strings.Replace(ogJSON, `"config_type":"templated"`, `"config_type":"unrecognized"`, 1), func(c *m.Client) error { _, e := c.GetOGConfig(ctx, "id"); return e }},
		{"missing OG selector", strings.Replace(ogJSON, `"template_id":"t-example"`, `"template_id":null`, 1), func(c *m.Client) error { _, e := c.GetOGConfig(ctx, "id"); return e }},
		{"missing writer role ARN", `{"external_id":"org-external"}`, func(c *m.Client) error { _, e := c.GetAWSExternalID(ctx); return e }},
		{"missing external ID", `{}`, func(c *m.Client) error { _, e := c.GetAWSExternalID(ctx); return e }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, tc.body) }))
			defer server.Close()
			err := tc.call(m.NewClient("id", "key", m.WithBaseURL(server.URL)))
			var responseErr *m.ResponseError
			if !errors.As(err, &responseErr) || responseErr.StatusCode != 200 {
				t.Fatalf("expected response error, got %v", err)
			}
		})
	}
}

func TestStorageReadVariantsAndHTMLCSSOG(t *testing.T) {
	for _, connection := range []string{
		`{"provider":"aws_s3","bucket":"images","role_arn":"role","region":"us-east-1"}`,
		`{"provider":"cloudflare_r2","bucket":"images","cloudflare_account_id":"account","access_key_id":"access"}`,
		`{"provider":"backblaze_b2","bucket":"images","region":"us-west-004","access_key_id":"access"}`,
		`{"provider":"digitalocean_spaces","bucket":"images","region":"nyc3","access_key_id":"access"}`,
		`{"provider":"wasabi","bucket":"images","region":"us-east-1","access_key_id":"access"}`,
		`{"provider":"google_cloud_storage","bucket":"images","access_key_id":"GOOG-access"}`,
		`{"provider":"other_s3_compatible","bucket":"images","endpoint":"https://s3.example","region":null,"force_path_style":false,"access_key_id":"access"}`,
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/v1/og-configs/id" {
				io.WriteString(w, `{"id":"id","domain_id":"domain","config_type":"html_css","base_url":"https://example.com","default_options":{"headers":{"Authorization":"Bearer secret"}},"extract_values":true}`)
				return
			}
			io.WriteString(w, `{"id":"id","name":"storage","connection_info":`+connection+`,"last_test_succeeded":false,"last_test_error":"connection failed","last_tested_at":"2026-09-15T12:00:00Z"}`)
		}))
		c := m.NewClient("id", "key", m.WithBaseURL(server.URL))
		storage, err := c.GetStorageDestination(context.Background(), "id")
		if err != nil {
			server.Close()
			t.Fatal(err)
		}
		if storage.LastTestSucceeded == nil || *storage.LastTestSucceeded || storage.LastTestError == nil || storage.LastTestedAt == nil {
			t.Error("connection-test status lost")
		}
		og, err := c.GetOGConfig(context.Background(), "id")
		if err != nil {
			server.Close()
			t.Fatal(err)
		}
		if !og.ExtractValues || og.DefaultOptions.Headers["Authorization"] != "Bearer secret" {
			t.Error("readback lost OG extraction or headers")
		}
		server.Close()
	}
}
