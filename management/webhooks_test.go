package management_test

import (
	"context"
	"encoding/json"
	hcti "github.com/htmlcsstoimage/go-client"
	m "github.com/htmlcsstoimage/go-client/management"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookClient(t *testing.T) {
	const response = `{"id":"hook-1","name":"test","url":"https://example.com/","enabled":true,"events":[],"parameters":[],"config_version":"2"}`
	var calls []string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.RequestURI())
		if id, key, ok := r.BasicAuth(); !ok || id != "id" || key != "key" {
			t.Error("missing authentication")
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "DELETE" {
			w.WriteHeader(204)
			return
		}
		if r.Method == "POST" {
			var req m.WebhookRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Error(err)
			}
			if req.Events == nil || len(req.Parameters) != 1 {
				t.Error("request shape", req)
			}
			p := req.Parameters[0]
			if r.URL.Path == "/v1/webhooks" && (p.Value == nil || *p.Value != "" || p.RetainValue) {
				t.Error("empty value lost")
			}
			if r.URL.Path == "/v1/webhooks/hook-1" && (p.Value != nil || !p.RetainValue) {
				t.Error("retention lost")
			}
		}
		if r.URL.RawQuery != "" {
			_, _ = w.Write([]byte(`{"data":[` + response + `],"pagination":{"next_page_start":"next"}}`))
			return
		}
		if r.Method == "POST" && r.URL.Path == "/v1/webhooks" {
			_, _ = w.Write([]byte(response[:len(response)-1] + `,"signing_secret":"whsec_test"}`))
			return
		}
		_, _ = w.Write([]byte(response))
	}))
	defer s.Close()
	c := m.NewClient("id", "key", m.WithBaseURL(s.URL))
	ctx := context.Background()
	req := &m.WebhookRequest{Name: "test", URL: "https://example.com/", Events: []string{}, Parameters: []m.WebhookParameterRequest{{Location: "header", Key: "Authorization", Secret: true, Value: hcti.Ptr("")}}}
	w, err := c.CreateWebhook(ctx, req)
	if err != nil || w.ConfigVersion != 2 {
		t.Fatalf("create: %v %v", w, err)
	}
	if w.SigningSecret == nil || *w.SigningSecret != "whsec_test" {
		t.Fatal("missing creation secret with signing disabled")
	}
	if got, err := c.GetWebhook(ctx, w.ID); err != nil || got.SigningSecret != nil {
		t.Fatal(err)
	}
	req.Parameters[0].Value = nil
	req.Parameters[0].RetainValue = true
	if got, err := c.UpdateWebhook(ctx, w.ID, req); err != nil || got.SigningSecret != nil {
		t.Fatal(err)
	}
	page, err := c.ListWebhooks(ctx, m.ListOptions{Count: hcti.Ptr(10), PageStart: hcti.Ptr("cursor")})
	if err != nil || len(page.Data) != 1 || *page.Pagination.NextPageStart != "next" {
		t.Fatalf("page: %v %v", page, err)
	}
	if page.Data[0].SigningSecret != nil {
		t.Fatal("list returned secret")
	}
	if err := c.DeleteWebhook(ctx, w.ID); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 5 || calls[3] != "GET /v1/webhooks?count=10&page_start=cursor" {
		t.Fatal(calls)
	}
	if _, err := c.CreateWebhook(ctx, nil); err == nil {
		t.Fatal("nil create accepted")
	}
	if _, err := c.UpdateWebhook(ctx, "hook-1", nil); err == nil {
		t.Fatal("nil update accepted")
	}
	if _, err := c.GetWebhook(ctx, ""); err == nil {
		t.Fatal("empty ID accepted")
	}
	if len(calls) != 5 {
		t.Fatal("invalid input made request")
	}
}
