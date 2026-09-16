package httpapi

import (
	"net/http"
	"testing"
	"time"
)

func TestConfigDefaultsAndCallerIsolation(t *testing.T) {
	config := NewConfig("id", "secret", "0.2.0")
	if config.BaseURL != "https://hcti.io" || config.userAgent != "HCTIGo/0.2.0" {
		t.Fatalf("incorrect defaults: origin=%s, agent=%s", config.BaseURL, config.userAgent)
	}
	if client := config.newHTTPClient(); client.Timeout != 60*time.Second {
		t.Fatalf("default timeout=%s", client.Timeout)
	}
	config.SetBaseURL("https://example.com///")
	if config.BaseURL != "https://example.com" {
		t.Fatal("origin wasn't normalized")
	}

	callerRedirect := func(*http.Request, []*http.Request) error { return nil }
	caller := &http.Client{Timeout: 3 * time.Second, CheckRedirect: callerRedirect}
	config.SetHTTPClient(caller)
	config.SetHTTPClient(nil) // Does not discard a previously supplied client.
	caller.Timeout = 9 * time.Second
	first, second := New(config), New(config)
	if first.httpClient.Timeout != 3*time.Second || second.httpClient.Timeout != 3*time.Second {
		t.Fatal("configuration did not snapshot caller settings")
	}
	if err := first.httpClient.CheckRedirect(nil, nil); err != http.ErrUseLastResponse {
		t.Fatal("redirects enabled")
	}
	if err := caller.CheckRedirect(nil, nil); err != nil {
		t.Fatal("caller redirect policy changed")
	}
	first.httpClient.Timeout = time.Second
	if second.httpClient.Timeout != 3*time.Second {
		t.Fatal("clients share mutable HTTP configuration")
	}
}

func TestUserAgentSuffix(t *testing.T) {
	config := NewConfig("id", "key", "0.2.0")
	config.AppendUserAgent(" HCTITerraform/dev ")
	config.AppendUserAgent("  ")
	if config.userAgent != "HCTIGo/0.2.0 HCTITerraform/dev" {
		t.Fatalf("unexpected user agent: %s", config.userAgent)
	}
}
