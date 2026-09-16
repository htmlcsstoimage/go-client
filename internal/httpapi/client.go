// Package httpapi implements shared SDK HTTP behavior. It is not a public API.
package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Client holds immutable connection settings and can be shared by goroutines.
type Client struct {
	config     Config
	httpClient *http.Client
}

func New(config Config) *Client {
	return &Client{config: config, httpClient: config.newHTTPClient()}
}

func CredentialsFromEnv() (string, string, error) {
	id, key := os.Getenv("HCTI_API_ID"), os.Getenv("HCTI_API_KEY")
	if id == "" || key == "" {
		return "", "", fmt.Errorf("hcti: HCTI_API_ID and HCTI_API_KEY are required")
	}
	return id, key, nil
}

func ResourcePath(resource, id string) (string, error) {
	if id == "" || id == "." || id == ".." {
		return "", fmt.Errorf("hcti: a resource ID is required")
	}
	return "/v1/" + resource + "/" + url.PathEscape(id), nil
}

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

// ResponseError indicates an invalid successful API response. The response body
// is deliberately excluded from the error so it is safe to log.
type ResponseError struct {
	StatusCode int
	Problem    string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("hcti: invalid API response (HTTP %d): %s", e.StatusCode, e.Problem)
}

func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body, result any, validate func() string) error {
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("hcti: encode request: %w", err)
		}
	}
	endpoint := c.config.BaseURL + path
	if len(query) != 0 {
		endpoint += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("hcti: create request: %w", err)
	}
	req.SetBasicAuth(c.config.APIID, c.config.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.config.userAgent)
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
	if validate != nil {
		if problem := validate(); problem != "" {
			return &ResponseError{StatusCode: resp.StatusCode, Problem: problem}
		}
	}
	return nil
}
