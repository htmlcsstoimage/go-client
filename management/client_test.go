package management_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	hcti "github.com/htmlcsstoimage/go-client"
	m "github.com/htmlcsstoimage/go-client/management"
)

const keyJSON = `{"id":"key-id","api_id":"auth-id","name":"automation","permissions":[],"enabled":true}`
const proxyJSON = `{"id":"proxy-id","name":"proxy","url":"https://proxy.example","port":"443","bypass_hosts":[],"username":"","enabled":true}`
const storageJSON = `{"id":"storage-id","name":"storage","connection_info":{"provider":"aws_s3","bucket":"images","region":"us-east-1","role_arn":"arn:aws:iam::123:role/hcti"}}`
const ogJSON = `{"id":"og-id","domain_id":"domain-id","config_type":"templated","base_url":"https://example.com","template_id":"t-example","template_version":null,"refresh_interval_s":"86400"}`

func TestResourceRoutes(t *testing.T) {
	ctx := context.Background()
	type action struct {
		name, method, path, response string
		call                         func(*m.Client) error
	}
	actions := []action{
		{"create key", "POST", "/v1/api-keys", strings.TrimSuffix(keyJSON, "}") + `,"api_key":"once-only"}`, func(c *m.Client) error {
			v, e := c.CreateAPIKey(ctx, &m.APIKeyRequest{Permissions: []m.Permission{}})
			if e == nil && v.Secret != "once-only" {
				t.Error("lost create secret")
			}
			return e
		}},
		{"get key", "GET", "/v1/api-keys/key-id", keyJSON, func(c *m.Client) error { _, e := c.GetAPIKey(ctx, "key-id"); return e }},
		{"update key", "POST", "/v1/api-keys/key-id", keyJSON, func(c *m.Client) error {
			_, e := c.UpdateAPIKey(ctx, "key-id", &m.APIKeyRequest{Permissions: []m.Permission{}})
			return e
		}},
		{"delete key", "DELETE", "/v1/api-keys/key-id", "", func(c *m.Client) error { return c.DeleteAPIKey(ctx, "key-id") }},
		{"list keys", "GET", "/v1/api-keys", `{"data":[` + keyJSON + `],"pagination":{"next_page_start":null}}`, func(c *m.Client) error { _, e := c.ListAPIKeys(ctx, m.APIKeyListOptions{}); return e }},
		{"create proxy", "POST", "/v1/proxies", proxyJSON, func(c *m.Client) error { _, e := c.CreateProxy(ctx, &m.ProxyRequest{}); return e }},
		{"get proxy", "GET", "/v1/proxies/proxy-id", proxyJSON, func(c *m.Client) error {
			v, e := c.GetProxy(ctx, "proxy-id")
			if e == nil && (v.Port == nil || *v.Port != 443 || v.Username == nil || *v.Username != "") {
				t.Error("lost numeric port/empty username")
			}
			return e
		}},
		{"update proxy", "POST", "/v1/proxies/proxy-id", proxyJSON, func(c *m.Client) error { _, e := c.UpdateProxy(ctx, "proxy-id", &m.ProxyRequest{}); return e }},
		{"delete proxy", "DELETE", "/v1/proxies/proxy-id", "", func(c *m.Client) error { return c.DeleteProxy(ctx, "proxy-id") }},
		{"list proxies", "GET", "/v1/proxies", `{"data":[]}`, func(c *m.Client) error { _, e := c.ListProxies(ctx, m.ListOptions{}); return e }},
		{"create storage", "POST", "/v1/storage-destinations", storageJSON, func(c *m.Client) error {
			_, e := c.CreateStorageDestination(ctx, &m.StorageDestinationRequest{ConnectionInfo: &m.AWSS3Connection{}})
			return e
		}},
		{"get storage", "GET", "/v1/storage-destinations/storage-id", storageJSON, func(c *m.Client) error { _, e := c.GetStorageDestination(ctx, "storage-id"); return e }},
		{"update storage", "POST", "/v1/storage-destinations/storage-id", storageJSON, func(c *m.Client) error {
			_, e := c.UpdateStorageDestination(ctx, "storage-id", &m.StorageDestinationRequest{ConnectionInfo: &m.AWSS3Connection{}})
			return e
		}},
		{"delete storage", "DELETE", "/v1/storage-destinations/storage-id", "", func(c *m.Client) error { return c.DeleteStorageDestination(ctx, "storage-id") }},
		{"list storage", "GET", "/v1/storage-destinations", `{"data":[]}`, func(c *m.Client) error { _, e := c.ListStorageDestinations(ctx, m.ListOptions{}); return e }},
		{"external ID", "GET", "/v1/storage-destinations/aws-external-id", `{"external_id":"org-external","writer_role_arn":"arn:aws:iam::123456789012:role/hcti-writer"}`, func(c *m.Client) error {
			v, e := c.GetAWSExternalID(ctx)
			if e == nil && (v.ExternalID != "org-external" || v.WriterRoleARN != "arn:aws:iam::123456789012:role/hcti-writer") {
				t.Error("lost AWS trust policy details")
			}
			return e
		}},
		{"create OG", "POST", "/v1/og-configs", ogJSON, func(c *m.Client) error { _, e := c.CreateOGConfig(ctx, &m.TemplatedOGConfigRequest{}); return e }},
		{"get OG", "GET", "/v1/og-configs/og-id", ogJSON, func(c *m.Client) error {
			v, e := c.GetOGConfig(ctx, "og-id")
			if e == nil && (v.TemplateVersion != nil || v.RefreshIntervalSeconds != 86400) {
				t.Error("lost nullable version/refresh interval")
			}
			return e
		}},
		{"update OG", "POST", "/v1/og-configs/og-id", ogJSON, func(c *m.Client) error { _, e := c.UpdateOGConfig(ctx, "og-id", &m.HTMLCSSOGConfigRequest{}); return e }},
		{"delete OG", "DELETE", "/v1/og-configs/og-id", "", func(c *m.Client) error { return c.DeleteOGConfig(ctx, "og-id") }},
		{"list OG", "GET", "/v1/og-configs", `{"data":[]}`, func(c *m.Client) error { _, e := c.ListOGConfigs(ctx, m.ListOptions{}); return e }},
	}
	for _, a := range actions {
		t.Run(a.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != a.method || r.URL.EscapedPath() != a.path {
					t.Errorf("got %s %s", r.Method, r.URL)
				}
				id, key, ok := r.BasicAuth()
				if !ok || id != "id" || key != "key" {
					t.Error("missing authentication")
				}
				if r.Header.Get("User-Agent") != "HCTIGo/"+hcti.Version() {
					t.Error("incorrect user agent")
				}
				if r.Header.Get("Accept") != "application/json" {
					t.Error("incorrect Accept")
				}
				if r.Method == "POST" {
					if r.Header.Get("Content-Type") != "application/json" {
						t.Error("incorrect Content-Type")
					}
					var body any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
				}
				if a.response == "" {
					w.WriteHeader(204)
				} else {
					io.WriteString(w, a.response)
				}
			}))
			defer server.Close()
			if err := a.call(m.NewClient("id", "key", m.WithBaseURL(server.URL))); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Errorf("calls=%d", calls)
			}
		})
	}
}

func TestPaginationAndEscaping(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			q := r.URL.Query()
			if q.Get("count") != "2" || q.Get("page_start") != "cursor/+?💜" || q.Get("include_disabled") != "true" || !reflect.DeepEqual(q["with_permission"], []string{"images:read", "proxies:read"}) {
				t.Errorf("query=%v", q)
			}
			io.WriteString(w, `{"data":[],"pagination":{"next_page_start":"next/+💜"}}`)
		} else {
			if r.URL.EscapedPath() != "/v1/proxies/a%2Fb%3F%23" {
				t.Errorf("path=%s", r.URL.EscapedPath())
			}
			io.WriteString(w, proxyJSON)
		}
	}))
	defer server.Close()
	c := m.NewClient("id", "key", m.WithBaseURL(server.URL))
	page, err := c.ListAPIKeys(context.Background(), m.APIKeyListOptions{ListOptions: m.ListOptions{Count: hcti.Ptr(2), PageStart: hcti.Ptr("cursor/+?💜")}, IncludeDisabled: true, WithPermission: []m.Permission{m.PermissionImagesRead, m.PermissionProxiesRead}})
	if err != nil {
		t.Fatal(err)
	}
	if page.Pagination.NextPageStart == nil || *page.Pagination.NextPageStart != "next/+💜" {
		t.Fatal("lost cursor")
	}
	if _, err := c.GetProxy(context.Background(), "a/b?#"); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"", ".", ".."} {
		if _, err := c.GetProxy(context.Background(), id); err == nil {
			t.Error("accepted invalid ID")
		}
	}
	for _, count := range []int{0, 101} {
		if _, err := c.ListProxies(context.Background(), m.ListOptions{Count: &count}); err == nil {
			t.Error("accepted invalid count")
		}
	}
	if calls != 2 {
		t.Errorf("unexpected extra requests: %d", calls)
	}
}

func TestErrorsAndNoRetry(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
	}{
		{"validation", `{"error":"invalid","message":"secret-message","validation_errors":[{"path":"password","message":"invalid secret"}]}`, 400},
		{"rate limit", `{"error":"limited"}`, 429},
		{"server error", `bad gateway`, 502},
		{"missing resource", `{}`, 200},
		{"malformed", `{"password":"secret-message"`, 200},
		{"trailing", proxyJSON + ` {}`, 200},
		{"overflow", strings.Replace(proxyJSON, `"443"`, `"65536"`, 1), 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Retry-After", "7")
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			}))
			defer server.Close()
			_, err := m.NewClient("id", "key", m.WithBaseURL(server.URL)).GetProxy(context.Background(), "id")
			if err == nil {
				t.Fatal("expected error")
			}
			if strings.Contains(err.Error(), "secret") {
				t.Fatal("error leaks response")
			}
			if tc.status == 200 {
				var e *m.ResponseError
				if !errors.As(err, &e) {
					t.Fatalf("%T", err)
				}
			} else {
				var e *hcti.APIError
				if !errors.As(err, &e) || e.StatusCode != tc.status || e.Headers.Get("Retry-After") != "7" {
					t.Fatalf("wrong API error: %#v", err)
				}
			}
			if calls != 1 {
				t.Fatalf("retried %d times", calls)
			}
		})
	}
}

func TestEnvRedirectAndCancellation(t *testing.T) {
	t.Setenv("HCTI_API_ID", "")
	t.Setenv("HCTI_API_KEY", "")
	if _, err := m.NewClientFromEnv(); err == nil {
		t.Fatal("accepted missing credentials")
	}
	t.Setenv("HCTI_API_ID", "env-id")
	t.Setenv("HCTI_API_KEY", "env-key")
	redirects := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, key, _ := r.BasicAuth()
		if id != "env-id" || key != "env-key" {
			t.Error("wrong env credentials")
		}
		if r.URL.Path == "/target" {
			redirects++
			io.WriteString(w, proxyJSON)
			return
		}
		http.Redirect(w, r, "/target", 302)
	}))
	defer server.Close()
	custom := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { t.Error("caller redirect callback invoked"); return nil }}
	c, err := m.NewClientFromEnv(m.WithBaseURL(server.URL), m.WithHTTPClient(custom))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetProxy(context.Background(), "id")
	var apiErr *m.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 302 {
		t.Fatalf("redirect: %v", err)
	}
	if redirects != 0 || custom.CheckRedirect == nil {
		t.Fatal("redirect followed or caller client changed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = c.GetProxy(ctx, "id")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
}
