package management

import (
	"context"
	"github.com/htmlcsstoimage/go-client/internal/httpapi"
	"net/http"
)

// CreateOGConfig creates a resource.
func (c *Client) CreateOGConfig(ctx context.Context, request OGConfigRequest) (*OGConfig, error) {
	if err := validateOGRequest(request); err != nil {
		return nil, err
	}
	path := "/v1/og-configs"
	var result OGConfig
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOGConfig reads resource metadata.
func (c *Client) GetOGConfig(ctx context.Context, id string) (*OGConfig, error) {
	path, err := httpapi.ResourcePath("og-configs", id)
	if err != nil {
		return nil, err
	}
	var result OGConfig
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOGConfig replaces editable settings; omitted options are not retained.
func (c *Client) UpdateOGConfig(ctx context.Context, id string, request OGConfigRequest) (*OGConfig, error) {
	if err := validateOGRequest(request); err != nil {
		return nil, err
	}
	path, err := httpapi.ResourcePath("og-configs", id)
	if err != nil {
		return nil, err
	}
	var result OGConfig
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteOGConfig deletes the resource.
func (c *Client) DeleteOGConfig(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("og-configs", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListOGConfigs reads one page of metadata.
func (c *Client) ListOGConfigs(ctx context.Context, options ListOptions) (*Page[OGConfig], error) {
	q, err := options.query()
	if err != nil {
		return nil, err
	}
	var result Page[OGConfig]
	if err := c.do(ctx, http.MethodGet, "/v1/og-configs", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
