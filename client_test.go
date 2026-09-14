package hcti

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAPIErrors(t *testing.T) {
	for _, body := range []string{`{"error":"invalid","message":"details","validation_errors":[{"path":"html","message":"required"}]}`, `<html>edge error</html>`} {
		calls := 0
		c := testClient(t, 429, body, func(*http.Request) { calls++ })
		_, err := c.CreateImage(context.Background(), HTMLImageRequest{HTML: "hello"})
		var apiError *APIError
		if !errors.As(err, &apiError) || apiError.StatusCode != 429 || apiError.Headers.Get("RateLimit") == "" || calls != 1 {
			t.Fatalf("error: %v; calls: %d", err, calls)
		}
		if strings.Contains(err.Error(), "details") {
			t.Fatal("error string includes response content")
		}
		if strings.HasPrefix(body, "{") && len(apiError.ValidationErrors) != 1 {
			t.Fatal(apiError)
		}
	}
}

func TestCancellationAndNoContent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := NewClient("id", "key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) { return nil, r.Context().Err() })}))
	if _, err := c.CreateImage(ctx, URLImageRequest{URL: "https://example.com"}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	c = testClient(t, 204, "", func(r *http.Request) {
		if r.URL.EscapedPath() != "/v1/image/a%2Fb" {
			t.Fatal(r.URL)
		}
	})
	if err := c.DeleteImage(context.Background(), "a/b"); err != nil {
		t.Fatal(err)
	}
}

func TestRedirectIsNotFollowed(t *testing.T) {
	calls := 0
	c := NewClient("id", "key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: 307, Header: http.Header{"Location": {"https://other.example.com"}}, Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
	})}))
	_, err := c.CreateImage(context.Background(), HTMLImageRequest{HTML: "hello"})
	if err == nil || calls != 1 {
		t.Fatalf("redirect followed: %d, %v", calls, err)
	}
}
