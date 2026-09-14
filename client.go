// Package hcti provides a Go client for the HTML/CSS to Image API.
package hcti

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client is safe for concurrent use. Configure it before sharing it between goroutines.
// Requests are never automatically retried: a failed create may already have succeeded.
type Client struct {
	apiID, apiKey string
	baseURL       string
	httpClient    *http.Client
}

// Option configures a client.
type Option func(*Client)

// WithHTTPClient supplies transport and timeout configuration. The client is copied;
// redirects are disabled to avoid replaying credentials or writes at another endpoint.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			copied := *client
			c.httpClient = &copied
		}
	}
}

// WithBaseURL overrides the API origin, for example for a local test server.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// NewClient constructs a client using HTTP Basic authentication.
func NewClient(apiID, apiKey string, options ...Option) *Client {
	c := &Client{apiID: apiID, apiKey: apiKey, baseURL: "https://hcti.io", httpClient: &http.Client{Timeout: 60 * time.Second}}
	for _, option := range options {
		option(c)
	}
	c.httpClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return c
}

// NewClientFromEnv reads HCTI_API_ID and HCTI_API_KEY.
func NewClientFromEnv(options ...Option) (*Client, error) {
	id, key := os.Getenv("HCTI_API_ID"), os.Getenv("HCTI_API_KEY")
	if id == "" || key == "" {
		return nil, fmt.Errorf("hcti: HCTI_API_ID and HCTI_API_KEY are required")
	}
	return NewClient(id, key, options...), nil
}

// Ptr supplies an explicit optional value, including false, zero, or an empty string.
func Ptr[T any](value T) *T { return &value }

// ValidationError identifies a rejected request field.
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// APIError is returned for non-2xx responses. Headers include rate-limit metadata.
// The response body is not included in Error(), to avoid accidentally logging secrets.
type APIError struct {
	StatusCode       int
	Code             string
	Message          string
	ValidationErrors []ValidationError
	Headers          http.Header
}

func (e *APIError) Error() string {
	return fmt.Sprintf("hcti: API request failed (HTTP %d)", e.StatusCode)
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, result any) error {
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("hcti: encode request: %w", err)
		}
	}
	endpoint := c.baseURL + path
	if len(query) != 0 {
		endpoint += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("hcti: create request: %w", err)
	}
	req.SetBasicAuth(c.apiID, c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("hcti: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var payload struct {
			Error                  string            `json:"error"`
			Message                string            `json:"message"`
			ValidationErrors       []ValidationError `json:"validation_errors"`
			LegacyValidationErrors []ValidationError `json:"validationErrors"`
		}
		// Edge errors may contain HTML instead of an API error document.
		_ = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload)
		if payload.ValidationErrors == nil {
			payload.ValidationErrors = payload.LegacyValidationErrors
		}
		return &APIError{StatusCode: resp.StatusCode, Code: payload.Error, Message: payload.Message, ValidationErrors: payload.ValidationErrors, Headers: resp.Header.Clone()}
	}
	if result == nil {
		_, err = io.Copy(io.Discard, resp.Body)
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result); err != nil {
		return &ResponseError{StatusCode: resp.StatusCode, Problem: "expected a JSON response"}
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return &ResponseError{StatusCode: resp.StatusCode, Problem: "unexpected content after JSON response"}
	}
	if value, ok := result.(interface{ validateResponse() string }); ok {
		if problem := value.validateResponse(); problem != "" {
			return &ResponseError{StatusCode: resp.StatusCode, Problem: problem}
		}
	}
	return nil
}

func resourcePath(resource, id string) (string, error) {
	if id == "" || id == "." || id == ".." {
		return "", fmt.Errorf("hcti: a resource ID is required")
	}
	return "/v1/" + resource + "/" + url.PathEscape(id), nil
}
