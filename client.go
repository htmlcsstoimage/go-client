// Package hcti provides a Go client for the HTML/CSS to Image API.
package hcti

import (
	"context"
	"net/http"
	"net/url"

	"github.com/htmlcsstoimage/go-client/internal/httpapi"
)

// Client is safe for concurrent use. Configure it before sharing it between goroutines.
// Requests are never automatically retried: a failed create may already have succeeded.
type Client struct {
	config    httpapi.Config
	transport *httpapi.Client
}

// Option configures a client.
type Option func(*Client)

// WithHTTPClient supplies transport and timeout configuration. The client is copied;
// redirects are disabled to avoid replaying credentials or writes at another endpoint.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) { c.config.SetHTTPClient(client) }
}

// WithBaseURL overrides the API origin, for example for a local test server.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.config.SetBaseURL(baseURL) }
}

// WithUserAgentSuffix appends an integration identifier, such as a provider name
// and version, to the SDK's User-Agent. It does not replace HCTIGo/<version>.
func WithUserAgentSuffix(suffix string) Option {
	return func(c *Client) { c.config.AppendUserAgent(suffix) }
}

// NewClient constructs a client using HTTP Basic authentication.
func NewClient(apiID, apiKey string, options ...Option) *Client {
	c := &Client{config: httpapi.NewConfig(apiID, apiKey, Version())}
	for _, option := range options {
		option(c)
	}
	c.transport = httpapi.New(c.config)
	return c
}

// NewClientFromEnv reads HCTI_API_ID and HCTI_API_KEY.
func NewClientFromEnv(options ...Option) (*Client, error) {
	id, key, err := httpapi.CredentialsFromEnv()
	if err != nil {
		return nil, err
	}
	return NewClient(id, key, options...), nil
}

// Ptr supplies an explicit optional value, including false, zero, or an empty string.
func Ptr[T any](value T) *T { return &value }

// ValidationError identifies a rejected request field.
type ValidationError = httpapi.ValidationError

// APIError describes a non-2xx response, including headers and validation details.
// Error() excludes response contents so credentials are not accidentally logged.
type APIError = httpapi.APIError

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, result any) error {
	var validate func() string
	if value, ok := result.(interface{ validateResponse() string }); ok {
		validate = value.validateResponse
	}
	return c.transport.Do(ctx, method, path, query, body, result, validate)
}

func resourcePath(resource, id string) (string, error) { return httpapi.ResourcePath(resource, id) }
