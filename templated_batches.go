package hcti

import (
	"context"
	"fmt"
	"net/http"
)

// TemplatedBatchImageOptions describes shared defaults or one template variation.
type TemplatedBatchImageOptions struct {
	// TemplateID inherits the default when empty. Supplying an ID resets the inherited version.
	TemplateID string `json:"template_id,omitempty"`
	// TemplateVersion inherits the default version, or uses latest when a TemplateID is supplied.
	TemplateVersion *int64 `json:"template_version,omitempty"`
	// TemplateValues objects merge recursively. Arrays, scalars and nil values replace defaults.
	TemplateValues map[string]any `json:"template_values,omitempty"`
	// Format selects the returned URL format, inheriting the default when empty.
	Format ImageFormat `json:"format,omitempty"`
}

// TemplatedBatchRequest creates images from one or more templates.
type TemplatedBatchRequest struct {
	DefaultOptions *TemplatedBatchImageOptions  `json:"default_options,omitempty"`
	Variations     []TemplatedBatchImageOptions `json:"variations"`
}

// CreateTemplatedImageBatch creates a template batch and preserves variation order.
// Template resolution and merging happen on the API. Empty batches succeed locally.
func (c *Client) CreateTemplatedImageBatch(ctx context.Context, request *TemplatedBatchRequest) (*BatchResult, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: template batch request is required")
	}
	if len(request.Variations) == 0 {
		return &BatchResult{Images: []Image{}}, nil
	}
	var result BatchResult
	if err := c.do(ctx, http.MethodPost, "/v1/image/batch/templated", nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
