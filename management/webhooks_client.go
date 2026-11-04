package management

import (
	"context"
	"fmt"
	"github.com/htmlcsstoimage/go-client/internal/httpapi"
	"net/http"
)

// CreateWebhook creates a resource.
func (c *Client) CreateWebhook(ctx context.Context, request *WebhookRequest) (*Webhook, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path := "/v1/webhooks"
	var result Webhook
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWebhook reads resource metadata.
func (c *Client) GetWebhook(ctx context.Context, id string) (*Webhook, error) {
	path, err := httpapi.ResourcePath("webhooks", id)
	if err != nil {
		return nil, err
	}
	var result Webhook
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateWebhook replaces editable settings; omitted options are not retained.
func (c *Client) UpdateWebhook(ctx context.Context, id string, request *WebhookRequest) (*Webhook, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path, err := httpapi.ResourcePath("webhooks", id)
	if err != nil {
		return nil, err
	}
	var result Webhook
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteWebhook deletes the resource.
func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("webhooks", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListWebhooks reads one page of metadata.
func (c *Client) ListWebhooks(ctx context.Context, options ListOptions) (*Page[Webhook], error) {
	q, err := options.query()
	if err != nil {
		return nil, err
	}
	var result Page[Webhook]
	if err := c.do(ctx, http.MethodGet, "/v1/webhooks", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
