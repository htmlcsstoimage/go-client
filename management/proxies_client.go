package management

import (
	"context"
	"fmt"
	"github.com/htmlcsstoimage/go-client/internal/httpapi"
	"net/http"
)

// CreateProxy creates a resource.
func (c *Client) CreateProxy(ctx context.Context, request *ProxyRequest) (*Proxy, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path := "/v1/proxies"
	var result Proxy
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProxy reads resource metadata.
func (c *Client) GetProxy(ctx context.Context, id string) (*Proxy, error) {
	path, err := httpapi.ResourcePath("proxies", id)
	if err != nil {
		return nil, err
	}
	var result Proxy
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProxy replaces editable settings; omitted options are not retained.
func (c *Client) UpdateProxy(ctx context.Context, id string, request *ProxyRequest) (*Proxy, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path, err := httpapi.ResourcePath("proxies", id)
	if err != nil {
		return nil, err
	}
	var result Proxy
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteProxy deletes the resource.
func (c *Client) DeleteProxy(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("proxies", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListProxies reads one page of metadata.
func (c *Client) ListProxies(ctx context.Context, options ListOptions) (*Page[Proxy], error) {
	q, err := options.query()
	if err != nil {
		return nil, err
	}
	var result Page[Proxy]
	if err := c.do(ctx, http.MethodGet, "/v1/proxies", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
