package management

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/htmlcsstoimage/go-client/internal/httpapi"
)

// CreatedRender identifies a saved image or newly created template version.
type CreatedRender struct {
	ID              string `json:"id"`
	URL             string `json:"url"`
	TemplateID      string `json:"template_id"`
	TemplateVersion int64  `json:"template_version"`
}

func (v *CreatedRender) validateResponse() string {
	if v.ID != "" && v.URL != "" {
		return ""
	}
	if v.TemplateID != "" && v.TemplateVersion > 0 {
		return ""
	}
	return "missing created rendering identity"
}

// CreateImageDefinition saves an independently owned image without rendering it.
// kind is html_css, url, or templated. Deduplication is always disabled.
// Nullable rendering inputs are serialized explicitly as null.
func (c *Client) CreateImageDefinition(ctx context.Context, kind string, request *RenderDefinition) (*CreatedRender, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path := "/v1/image"
	switch kind {
	case "html_css", "url":
	case "templated":
		if request.TemplateID == nil || !strings.HasPrefix(*request.TemplateID, "t-") {
			return nil, fmt.Errorf("hcti: template_id must have the t- prefix")
		}
		var err error
		path, err = httpapi.ResourcePath("image", *request.TemplateID)
		if err != nil {
			return nil, err
		}
		if request.TemplateVersion != nil {
			if *request.TemplateVersion <= 0 {
				return nil, fmt.Errorf("hcti: template_version must be positive")
			}
			path += "/" + strconv.FormatInt(*request.TemplateVersion, 10)
		}
	default:
		return nil, fmt.Errorf("hcti: invalid image definition kind")
	}
	var result CreatedRender
	if err := c.do(ctx, http.MethodPost, path, nil, request.payload(kind), &result); err != nil {
		return nil, err
	}
	if result.ID == "" || result.URL == "" {
		return nil, &ResponseError{StatusCode: 200, Problem: "missing created image identity"}
	}
	return &result, nil
}

// GetImageMetadata reads a saved image without invoking its rendering URL.
func (c *Client) GetImageMetadata(ctx context.Context, id string) (*SavedRender, error) {
	path, err := httpapi.ResourcePath("images", id)
	if err != nil {
		return nil, err
	}
	var result SavedRender
	if err = c.do(ctx, http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	if result.ID != id || result.CreatedAt == "" || (result.ImageType != "html_css" && result.ImageType != "url" && result.ImageType != "templated") {
		return nil, &ResponseError{StatusCode: 200, Problem: "missing or mismatched image metadata"}
	}
	if (result.ImageType == "html_css" && result.HTML == nil) || (result.ImageType == "url" && result.URL == nil) || (result.ImageType == "templated" && (result.TemplateID == nil || result.TemplateVersion == nil || *result.TemplateVersion <= 0 || len(result.TemplateValues) == 0)) {
		return nil, &ResponseError{StatusCode: 200, Problem: "missing image source"}
	}
	if result.ImageType == "templated" {
		var values map[string]json.RawMessage
		if err := json.Unmarshal(result.TemplateValues, &values); err != nil || len(values) == 0 {
			return nil, &ResponseError{StatusCode: 200, Problem: "invalid stored template values"}
		}
	}
	return &result, nil
}

// DeleteImageDefinition schedules image deletion. A successful asynchronous response is sufficient.
func (c *Client) DeleteImageDefinition(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("image", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// SaveTemplateDefinition creates a template when id is empty, or adds a version
// to an existing template. Null inputs restore API-default rendering behavior.
func (c *Client) SaveTemplateDefinition(ctx context.Context, id string, request *RenderDefinition) (*CreatedRender, error) {
	if request == nil {
		return nil, fmt.Errorf("hcti: request is required")
	}
	path := "/v1/template"
	if id != "" {
		var err error
		path, err = httpapi.ResourcePath("template", id)
		if err != nil {
			return nil, err
		}
	}
	var result CreatedRender
	if err := c.do(ctx, http.MethodPost, path, nil, request.payload("template"), &result); err != nil {
		return nil, err
	}
	if result.TemplateID == "" || result.TemplateVersion <= 0 || (id != "" && result.TemplateID != id) {
		return nil, &ResponseError{StatusCode: 200, Problem: "missing or mismatched template identity"}
	}
	return &result, nil
}

// GetTemplateDefinition reads the latest version when version is nil, otherwise
// searches the exclusive-cursor version listing for the exact requested version.
func (c *Client) GetTemplateDefinition(ctx context.Context, id string, version *int64) (*SavedRender, error) {
	path, err := httpapi.ResourcePath("template", id)
	if err != nil {
		return nil, err
	}
	q := url.Values{"count": {"1"}}
	if version != nil {
		q.Set("count", "100")
	}
	lastCursor := ""
	for {
		var page struct {
			Data       []SavedRender `json:"data"`
			Pagination struct {
				NextPageStart *json.Number `json:"next_page_start"`
			} `json:"pagination"`
		}
		if err := c.do(ctx, http.MethodGet, path, q, nil, &page); err != nil {
			return nil, err
		}
		if page.Data == nil {
			return nil, &ResponseError{StatusCode: 200, Problem: "missing template list"}
		}
		for _, v := range page.Data {
			if v.ID != id || v.Version <= 0 || v.CreatedAt == "" || v.TemplateType == "" || (v.TemplateType == "html_css" && v.HTML == nil) {
				return nil, &ResponseError{StatusCode: 200, Problem: "missing or mismatched template metadata"}
			}
			if version == nil || v.Version == *version {
				return &v, nil
			}
		}
		if version == nil || page.Pagination.NextPageStart == nil {
			return nil, &APIError{StatusCode: 404}
		}
		cursorNumber, err := strconv.ParseInt(page.Pagination.NextPageStart.String(), 10, 64)
		if err != nil {
			return nil, &ResponseError{StatusCode: 200, Problem: "invalid template pagination cursor"}
		}
		cursor := strconv.FormatInt(cursorNumber, 10)
		if cursor == lastCursor || (lastCursor != "" && func() bool { last, _ := strconv.ParseInt(lastCursor, 10, 64); return cursorNumber >= last }()) {
			return nil, &ResponseError{StatusCode: 200, Problem: "template pagination did not advance"}
		}
		lastCursor = cursor
		q.Set("max_version", cursor)
	}
}

// DeleteTemplateDefinition deletes the template and all of its versions.
func (c *Client) DeleteTemplateDefinition(ctx context.Context, id string) error {
	path, err := httpapi.ResourcePath("template", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ImageOperationURL constructs an operation URL after the caller establishes
// the storage mode. It does not request or render the image.
func (c *Client) ImageOperationURL(id, format string, storageOnly bool) (string, error) {
	resource := "image"
	if storageOnly {
		resource = "store"
	}
	path, err := httpapi.ResourcePath(resource, id)
	if err != nil {
		return "", err
	}
	if format != "" {
		switch format {
		case "png", "jpg", "jpeg", "webp", "pdf":
			path += "." + format
		default:
			return "", fmt.Errorf("hcti: invalid image format")
		}
	}
	return c.config.BaseURL + path, nil
}
