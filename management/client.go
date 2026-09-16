package management

import (
	"context"
	"net/http"
	"net/url"

	hcti "github.com/htmlcsstoimage/go-client"
	"github.com/htmlcsstoimage/go-client/internal/httpapi"
)

// Client accesses management endpoints and is safe for concurrent use.
type Client struct {
	config    httpapi.Config
	transport *httpapi.Client
}

// Option configures a management client before use.
type Option func(*Client)

// WithBaseURL overrides the API origin, for example for a local test server.
func WithBaseURL(baseURL string) Option { return func(c *Client) { c.config.SetBaseURL(baseURL) } }

// WithHTTPClient supplies timeout and transport settings. The client is copied;
// redirects are disabled to avoid forwarding credentials or replaying writes.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) { c.config.SetHTTPClient(client) }
}

// WithUserAgentSuffix appends an integration identifier, such as a provider name
// and version, to the SDK's User-Agent. It does not replace HCTIGo/<version>.
func WithUserAgentSuffix(suffix string) Option {
	return func(c *Client) { c.config.AppendUserAgent(suffix) }
}

// NewClient constructs a management client using HTTP Basic authentication.
func NewClient(apiID, apiKey string, options ...Option) *Client {
	c := &Client{config: httpapi.NewConfig(apiID, apiKey, hcti.Version())}
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

// APIError describes an unsuccessful HTTP response, including validation details.
type APIError = hcti.APIError

// ValidationError identifies a rejected request field.
type ValidationError = hcti.ValidationError

// ResponseError indicates malformed successful response data.
type ResponseError = hcti.ResponseError

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, result any) error {
	var validate func() string
	if v, ok := result.(interface{ validateResponse() string }); ok {
		validate = v.validateResponse
	}
	return c.transport.Do(ctx, method, path, query, body, result, validate)
}
