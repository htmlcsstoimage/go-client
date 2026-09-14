package hcti

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

// Image identifies a created image and its rendering URL.
type Image struct {
	// ID is the image identifier used for deletion and rendering.
	ID string `json:"id"`
	// URL is the render URL returned by the API.
	URL string `json:"url"`
}

// CreateImage creates an image definition and returns its render URL.
func (c *Client) CreateImage(ctx context.Context, request ImageRequest) (*Image, error) {
	if request == nil || (reflect.ValueOf(request).Kind() == reflect.Ptr && reflect.ValueOf(request).IsNil()) {
		return nil, fmt.Errorf("hcti: image request is required")
	}
	var result Image
	if err := c.do(ctx, http.MethodPost, "/v1/image", nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchRequest creates HTML or URL variations with optional shared defaults.
// Templated requests are not supported by the batch endpoint.
type BatchRequest struct {
	// Variations contains HTML or URL requests. Omitted fields inherit DefaultOptions.
	// An empty list returns an empty result without calling the API.
	Variations []ImageRequest `json:"variations"`
	// DefaultOptions supplies shared HTML or URL settings. Nil omits defaults.
	DefaultOptions ImageRequest `json:"default_options,omitempty"`
}

// BatchResult contains the images created by a batch request.
type BatchResult struct {
	// Images contains the created image identifiers and rendering URLs.
	Images []Image `json:"images"`
}

// CreateImageBatch creates HTML or URL variations with optional shared defaults.
// Templated requests are rejected.
func (c *Client) CreateImageBatch(ctx context.Context, request BatchRequest) (*BatchResult, error) {
	if len(request.Variations) == 0 {
		return &BatchResult{Images: []Image{}}, nil
	}
	requests := append([]ImageRequest{}, request.Variations...)
	if request.DefaultOptions != nil {
		requests = append(requests, request.DefaultOptions)
	}
	for _, item := range requests {
		switch v := item.(type) {
		case HTMLImageRequest, URLImageRequest:
		case *HTMLImageRequest:
			if v == nil {
				return nil, fmt.Errorf("hcti: batch request must not be nil")
			}
		case *URLImageRequest:
			if v == nil {
				return nil, fmt.Errorf("hcti: batch request must not be nil")
			}
		default:
			return nil, fmt.Errorf("hcti: batches support only HTML and URL requests")
		}
	}
	var result BatchResult
	if err := c.do(ctx, http.MethodPost, "/v1/image/batch", nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteImage deletes the image identified by id.
func (c *Client) DeleteImage(ctx context.Context, id string) error {
	path, err := resourcePath("image", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// DeleteImageBatch deletes the supplied image IDs in one request.
// An empty list succeeds without calling the API.
func (c *Client) DeleteImageBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return c.do(ctx, http.MethodDelete, "/v1/image/batch", nil, struct {
		IDs []string `json:"ids"`
	}{ids}, nil)
}
