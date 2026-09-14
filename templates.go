package hcti

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TemplateRequest creates an HTML/CSS template or a new version of one.
type TemplateRequest struct {
	RenderOptions
	// HTML contains Handlebars markup. It must be non-empty, include at least
	// one placeholder, and compile as valid Handlebars.
	HTML string `json:"html"`
	// CSS styles the rendered template. Handlebars expressions are not supported
	// here; put dynamic CSS inside HTML instead.
	CSS *string `json:"css,omitempty"`
	// GoogleFonts lists font families to load. Use font-family in CSS to select them.
	// The SDK serializes the list to the API's pipe-delimited format.
	GoogleFonts GoogleFonts `json:"google_fonts,omitempty"`
	// Name identifies the template in your account. Maximum length: 64 characters.
	Name string `json:"name,omitempty"`
	// Description explains the template's purpose. Maximum length: 1024 characters.
	Description string `json:"description,omitempty"`
}

// TemplateVersion is the result of creating a template or adding a version.
type TemplateVersion struct {
	// TemplateID is the stable identifier shared by versions of this template.
	TemplateID string `json:"template_id"`
	// TemplateVersion is the newly created version number. Use it to pin rendering.
	TemplateVersion int64 `json:"template_version"`
}

// Template contains common template metadata and HTML/CSS settings.
// Raw retains the complete response, including editor block fields not yet modeled.
type Template struct {
	RenderOptions
	// ID is the stable template identifier.
	ID string `json:"id"`
	// Version identifies this particular version of the template.
	Version int64 `json:"version"`
	// TemplateType identifies the template variant, such as html_css or blocks.
	TemplateType string `json:"template_type"`
	// Name is the template's display name.
	Name string `json:"name"`
	// Description explains the template's purpose.
	Description string `json:"description"`
	// HTML contains the markup for an HTML/CSS template.
	HTML string `json:"html"`
	// CSS contains the styles for an HTML/CSS template.
	CSS string `json:"css"`
	// GoogleFonts contains decoded font family names.
	GoogleFonts GoogleFonts `json:"google_fonts"`
	// CreatedAt is the creation timestamp returned by the API.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last-update timestamp returned by the API.
	UpdatedAt time.Time `json:"updated_at"`
	// RenderCount is the render count returned by the API.
	RenderCount int64 `json:"render_count"`
	// Raw preserves the complete JSON response, including fields specific to
	// template-editor block variants that are not yet represented by this type.
	Raw json.RawMessage `json:"-"`
}

func (t *Template) UnmarshalJSON(data []byte) error {
	type plain Template
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*t = Template(value)
	t.Raw = append(json.RawMessage(nil), data...)
	return nil
}

// TemplateListOptions uses the template API's numeric version cursor.
type TemplateListOptions struct {
	// Count limits the number of results, from 1 to 100.
	// Zero omits the parameter and uses the API default of 10.
	Count int
	// MaxVersion supplies the version cursor from the previous page's
	// Pagination.NextPageStart. Nil starts from the newest results.
	MaxVersion *int64
}

// TemplatePage is one page of templates or template versions.
type TemplatePage struct {
	// Data contains the templates or template versions on this page.
	Data []Template `json:"data"`
	// Pagination contains the cursor for the next page.
	Pagination struct {
		// NextPageStart is passed as TemplateListOptions.MaxVersion.
		// Nil indicates that no further page is available.
		NextPageStart *int64 `json:"next_page_start"`
	} `json:"pagination"`
}

// CreateTemplate creates an HTML/CSS template and its initial version.
func (c *Client) CreateTemplate(ctx context.Context, request TemplateRequest) (*TemplateVersion, error) {
	return c.writeTemplate(ctx, "/v1/template", request)
}

// CreateTemplateVersion adds a version without replacing the template's stable ID.
func (c *Client) CreateTemplateVersion(ctx context.Context, id string, request TemplateRequest) (*TemplateVersion, error) {
	path, err := resourcePath("template", id)
	if err != nil {
		return nil, err
	}
	return c.writeTemplate(ctx, path, request)
}

func (c *Client) writeTemplate(ctx context.Context, path string, request TemplateRequest) (*TemplateVersion, error) {
	var result TemplateVersion
	if err := c.do(ctx, http.MethodPost, path, nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListTemplates lists templates with their most recent versions.
func (c *Client) ListTemplates(ctx context.Context, options TemplateListOptions) (*TemplatePage, error) {
	return c.listTemplates(ctx, "/v1/template", options)
}

// ListTemplateVersions lists versions of a single template, newest first.
func (c *Client) ListTemplateVersions(ctx context.Context, id string, options TemplateListOptions) (*TemplatePage, error) {
	path, err := resourcePath("template", id)
	if err != nil {
		return nil, err
	}
	return c.listTemplates(ctx, path, options)
}

func (c *Client) listTemplates(ctx context.Context, path string, options TemplateListOptions) (*TemplatePage, error) {
	query := url.Values{}
	if options.Count != 0 {
		query.Set("count", strconv.Itoa(options.Count))
	}
	if options.MaxVersion != nil {
		query.Set("max_version", strconv.FormatInt(*options.MaxVersion, 10))
	}
	var result TemplatePage
	if err := c.do(ctx, http.MethodGet, path, query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTemplate removes the entire template and clears its contents.
// Previously rendered images may remain cached; new renders using the template fail.
func (c *Client) DeleteTemplate(ctx context.Context, id string) error {
	path, err := resourcePath("template", id)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
