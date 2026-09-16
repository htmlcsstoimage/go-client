package management

import (
	"context"
	"fmt"
	"github.com/htmlcsstoimage/go-client/internal/httpapi"
	"net/http"
)

// CreateStorageDestination creates a resource.
func (c *Client) CreateStorageDestination(ctx context.Context, request *StorageDestinationRequest) (*StorageDestination, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	if err := validateStorageConnection(request.ConnectionInfo); err != nil {
		return nil, err
	}
	path := "/v1/storage-destinations"
	var result StorageDestination
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStorageDestination reads resource metadata.
func (c *Client) GetStorageDestination(ctx context.Context, id string) (*StorageDestination, error) {
	path, err := httpapi.ResourcePath("storage-destinations", id)
	if err != nil {
		return nil, err
	}
	var result StorageDestination
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateStorageDestination replaces editable settings; omitted options are not retained.
func (c *Client) UpdateStorageDestination(ctx context.Context, id string, request *StorageDestinationRequest) (*StorageDestination, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	if err := validateStorageConnection(request.ConnectionInfo); err != nil {
		return nil, err
	}
	path, err := httpapi.ResourcePath("storage-destinations", id)
	if err != nil {
		return nil, err
	}
	var result StorageDestination
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteStorageDestination deletes the resource.
func (c *Client) DeleteStorageDestination(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("storage-destinations", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListStorageDestinations reads one page of metadata.
func (c *Client) ListStorageDestinations(ctx context.Context, options ListOptions) (*Page[StorageDestination], error) {
	q, err := options.query()
	if err != nil {
		return nil, err
	}
	var result Page[StorageDestination]
	if err := c.do(ctx, http.MethodGet, "/v1/storage-destinations", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
