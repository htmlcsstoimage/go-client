package management

import (
	"context"
	"fmt"
	"github.com/htmlcsstoimage/go-client/internal/httpapi"
	"net/http"
)

// CreateAPIKey creates a resource. The secret is returned only once.
func (c *Client) CreateAPIKey(ctx context.Context, request *APIKeyRequest) (*APIKeyWithSecret, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path := "/v1/api-keys"
	var result APIKeyWithSecret
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAPIKey reads resource metadata.
func (c *Client) GetAPIKey(ctx context.Context, id string) (*APIKey, error) {
	path, err := httpapi.ResourcePath("api-keys", id)
	if err != nil {
		return nil, err
	}
	var result APIKey
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateAPIKey replaces editable settings; omitted options are not retained.
func (c *Client) UpdateAPIKey(ctx context.Context, id string, request *APIKeyRequest) (*APIKey, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path, err := httpapi.ResourcePath("api-keys", id)
	if err != nil {
		return nil, err
	}
	var result APIKey
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAPIKey disables the key. Its metadata remains readable.
func (c *Client) DeleteAPIKey(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("api-keys", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListAPIKeys reads one page of metadata.
func (c *Client) ListAPIKeys(ctx context.Context, options APIKeyListOptions) (*Page[APIKey], error) {
	q, err := options.query()
	if err != nil {
		return nil, err
	}
	if options.IncludeDisabled {
		q.Set("include_disabled", "true")
	}
	for _, p := range options.WithPermission {
		q.Add("with_permission", string(p))
	}
	var result Page[APIKey]
	if err := c.do(ctx, http.MethodGet, "/v1/api-keys", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
