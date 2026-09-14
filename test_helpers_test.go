package hcti

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testClient(t *testing.T, status int, body string, check func(*http.Request)) *Client {
	t.Helper()
	return NewClient("test-id", "test-secret", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		id, key, ok := r.BasicAuth()
		if !ok || id != "test-id" || key != "test-secret" {
			t.Fatal("missing Basic authentication")
		}
		if check != nil {
			check(r)
		}
		return &http.Response{StatusCode: status, Header: http.Header{"Ratelimit": {"remaining=0"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}))
}
