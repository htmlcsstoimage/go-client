package management

import (
	"encoding/json"
	"strconv"

	hcti "github.com/htmlcsstoimage/go-client"
)

// RenderDefinition contains saved rendering inputs. Methods select the fields
// supported by their endpoint and serialize absent nullable inputs as JSON null.
// TemplateID and TemplateVersion select a template through the URL, never the body.
type RenderDefinition struct {
	// RequestOverrides blocks matching browser network requests on paid plans.
	RequestOverrides []hcti.RequestOverride `json:"request_overrides"`
	// HTML HTML to render and take a screenshot of. HTML fragments are rendered in a wrapper document unless a complete HTML document is supplied. Required for HTML image requests.
	HTML *string `json:"html"`
	// Name The name of the template, used to identify it in your account. Maximum: 64 characters.
	Name *string `json:"name"`
	// Description An optional description of the template for your reference. Maximum: 1024 characters.
	Description *string `json:"description"`
	// CSS CSS injected into the loaded URL to override styles on the page.
	CSS *string `json:"css"`
	// DeviceScale Adjusts the pixel ratio used for the screenshot. Minimum: 0.1. Maximum: 3. HTML and template renders default to 2; URL renders default to 1.
	DeviceScale *float64 `json:"device_scale"`
	// GoogleFonts Google fonts to load. Separate multiple fonts with a pipe, such as 'Roboto|OpenSans', and set font-family in the CSS to use them.
	GoogleFonts *string `json:"google_fonts"`
	// MaxWaitMS Sets a limit on how long to wait before taking the screenshot when the page continues loading irrelevant content. Minimum: 500. Maximum: 10000 and subject to the account plan limit.
	MaxWaitMS *int64 `json:"max_wait_ms"`
	// MSDelay Adds extra time in milliseconds before taking the screenshot so JavaScript can execute. Minimum: 0. Maximum: 10000.
	MSDelay *int64 `json:"ms_delay"`
	// RenderWhenReady Waits until the page signals that the screenshot is ready. The image fails if the readiness signal is never sent.
	RenderWhenReady *bool `json:"render_when_ready"`
	// Selector A CSS selector for an element in the HTML. We’ll crop the image to this specific element.
	Selector *string `json:"selector"`
	// ViewportHeight Sets the height of Chrome's viewport and disables automatic cropping. Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together.
	ViewportHeight *int64 `json:"viewport_height"`
	// ViewportWidth Sets the width of Chrome's viewport and disables automatic cropping. Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together.
	ViewportWidth *int64 `json:"viewport_width"`
	// DisableTwemoji: HTML/CSS and template images use Twemoji by default; true disables it. URL images inject Twemoji only with explicit false; nil or true leaves the page unchanged.
	DisableTwemoji *bool `json:"disable_twemoji"`
	// ColorScheme Rendering option.
	ColorScheme *string `json:"color_scheme"`
	// Timezone Sets the IANA timezone used by Chrome while rendering. Must be a recognized IANA timezone identifier.
	Timezone *string `json:"timezone"`
	// ViewportMobile Specifies whether the page uses mobile viewport behavior, including its viewport meta tag.
	ViewportMobile *bool `json:"viewport_mobile"`
	// ViewportLandscape Specifies whether the emulated viewport is in landscape orientation.
	ViewportLandscape *bool `json:"viewport_landscape"`
	// ViewportTouch Specifies whether the emulated viewport supports touch events.
	ViewportTouch *bool `json:"viewport_touch"`
	// MediaType Rendering option.
	MediaType *string `json:"media_type"`
	// ProxyID Specifies which configured organization proxy to use when rendering.
	ProxyID *string `json:"proxy_id"`
	// StorageDestinationID Specifies which configured organization storage destination receives the rendered image.
	StorageDestinationID *string `json:"storage_destination_id"`
	// JumboMaxHeight Maximum height of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_width. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels.
	JumboMaxHeight *int64 `json:"jumbo_max_height"`
	// JumboMaxWidth Maximum width of the rendered image in jumbo mode. Jumbo rendering consumes additional renders and requires jumbo_max_height. Supply both jumbo dimensions. Each must be greater than 0 and no more than 80000; at least one must exceed 8000; total area cannot exceed 400000000 pixels.
	JumboMaxWidth *int64 `json:"jumbo_max_width"`
	// TransparentBackground Specifies whether the image is rendered with a transparent background.
	TransparentBackground *bool `json:"transparent_background"`
	// Metadata Custom key-value metadata stored with the image.
	Metadata map[string]string `json:"metadata"`
	// PDFOptions Rendering option.
	PDFOptions *RenderPDFOptions `json:"pdf_options"`
	// MaxRenderOnce Ensure the image is only ever rendered and saved one time. This is an advanced option not applicable to most requests.
	MaxRenderOnce *bool `json:"max_render_once"`
	// Format Rendering option.
	Format *string `json:"format"`
	// URL Public HTTP or HTTPS URL to capture. Required for URL image requests.
	URL *string `json:"url"`
	// FullScreen Take a screenshot of the entire screen after scrolling down and back to the top.
	FullScreen *bool `json:"full_screen"`
	// BlockConsentBanners Attempt to block cookie/consent banners from displaying.
	BlockConsentBanners *bool `json:"block_consent_banners"`
	// IdentifyAsHCTI Identify the top-level page navigation as an HCTI screenshot request using the X-HCTI-SCREENSHOT header.
	IdentifyAsHCTI *bool `json:"identify_as_hcti"`
	// Headers HTTP headers to include on top-level page navigations to the requested URL's origin and any additional_header_origins. Supports up to 20 headers with names up to 512 ASCII characters and values up to 8192 UTF-8 bytes. For GET and form-encoded requests, repeat this parameter using the format `headers=name:value`.
	Headers map[string]string `json:"headers"`
	// AdditionalHeaderOrigins Additional exact HTTP or HTTPS origins allowed to receive custom headers. Supports up to 20 unique origins of up to 512 UTF-8 bytes each. Origins must use the format scheme://host[:port] without a path; duplicates are ignored. For GET and form-encoded requests, repeat this parameter for each origin.
	AdditionalHeaderOrigins []string `json:"additional_header_origins"`
	// IncludeHeadersOnSubrequests Include custom headers on subrequests to the requested URL's origin and any additional_header_origins. Defaults to false. Requires at least one header.
	IncludeHeadersOnSubrequests *bool `json:"include_headers_on_subrequests"`
	// TemplateValues Values substituted into the template for this render. Must be a non-empty JSON object. Include the values needed by the template.
	TemplateValues json.RawMessage `json:"template_values"`
	// TemplateID Template ID including the t- prefix.
	TemplateID *string `json:"template_id"`
	// TemplateVersion Optional version to pin. Omission selects the latest version when creating the image.
	TemplateVersion *int64 `json:"template_version"`
}

// RenderPDFOptions configures PDF output. Nil fields defer to API defaults.
type RenderPDFOptions struct {
	PageWidth       *string  `json:"page_width"`
	PageHeight      *string  `json:"page_height"`
	Scale           *float64 `json:"scale"`
	PrintBackground *bool    `json:"print_background,omitempty"`
	Margins         []string `json:"margins"`
}

// SavedRender is image metadata or one saved template version.
type SavedRender struct {
	RenderDefinition
	ID                                    string  `json:"id"`
	ImageType                             string  `json:"image_type"`
	TemplateType                          string  `json:"template_type"`
	Version                               int64   `json:"version"`
	CreatedAt                             string  `json:"created_at"`
	UpdatedAt                             string  `json:"updated_at"`
	LastRenderStoredAt                    *string `json:"last_render_stored_at"`
	SavedToStorageDestinationAt           *string `json:"saved_to_storage_destination_at"`
	OGConfigID                            *string `json:"og_config_id"`
	OGConfigContentVersion                *int64  `json:"og_config_content_version"`
	StorageDestinationHCTIStorageDisabled *bool   `json:"storage_destination_hcti_storage_disabled"`
}

// UnmarshalJSON accepts numeric strings without losing int64 precision.
func (v *SavedRender) UnmarshalJSON(data []byte) error {
	type plain SavedRender
	var p plain
	w := struct {
		*plain
		DeviceScale            *json.Number `json:"device_scale"`
		MaxWaitMS              *json.Number `json:"max_wait_ms"`
		MSDelay                *json.Number `json:"ms_delay"`
		ViewportHeight         *json.Number `json:"viewport_height"`
		ViewportWidth          *json.Number `json:"viewport_width"`
		JumboMaxHeight         *json.Number `json:"jumbo_max_height"`
		JumboMaxWidth          *json.Number `json:"jumbo_max_width"`
		TemplateVersion        *json.Number `json:"template_version"`
		Version                *json.Number `json:"version"`
		OGConfigContentVersion *json.Number `json:"og_config_content_version"`
	}{plain: &p}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	if w.DeviceScale != nil {
		n, err := strconv.ParseFloat(w.DeviceScale.String(), 64)
		if err != nil {
			return err
		}
		p.DeviceScale = &n
	}
	if w.MaxWaitMS != nil {
		n, err := strconv.ParseInt(w.MaxWaitMS.String(), 10, 64)
		if err != nil {
			return err
		}
		p.MaxWaitMS = &n
	}
	if w.MSDelay != nil {
		n, err := strconv.ParseInt(w.MSDelay.String(), 10, 64)
		if err != nil {
			return err
		}
		p.MSDelay = &n
	}
	if w.ViewportHeight != nil {
		n, err := strconv.ParseInt(w.ViewportHeight.String(), 10, 64)
		if err != nil {
			return err
		}
		p.ViewportHeight = &n
	}
	if w.ViewportWidth != nil {
		n, err := strconv.ParseInt(w.ViewportWidth.String(), 10, 64)
		if err != nil {
			return err
		}
		p.ViewportWidth = &n
	}
	if w.JumboMaxHeight != nil {
		n, err := strconv.ParseInt(w.JumboMaxHeight.String(), 10, 64)
		if err != nil {
			return err
		}
		p.JumboMaxHeight = &n
	}
	if w.JumboMaxWidth != nil {
		n, err := strconv.ParseInt(w.JumboMaxWidth.String(), 10, 64)
		if err != nil {
			return err
		}
		p.JumboMaxWidth = &n
	}
	if w.TemplateVersion != nil {
		n, err := strconv.ParseInt(w.TemplateVersion.String(), 10, 64)
		if err != nil {
			return err
		}
		p.TemplateVersion = &n
	}
	if w.Version != nil {
		n, err := strconv.ParseInt(w.Version.String(), 10, 64)
		if err != nil {
			return err
		}
		p.Version = n
	}
	if w.OGConfigContentVersion != nil {
		n, err := strconv.ParseInt(w.OGConfigContentVersion.String(), 10, 64)
		if err != nil {
			return err
		}
		p.OGConfigContentVersion = &n
	}
	*v = SavedRender(p)
	return nil
}

func (v RenderDefinition) payload(kind string) map[string]any {
	switch kind {
	case "template":
		return map[string]any{
			"request_overrides":      v.RequestOverrides,
			"html":                   v.HTML,
			"name":                   v.Name,
			"description":            v.Description,
			"css":                    v.CSS,
			"device_scale":           v.DeviceScale,
			"google_fonts":           v.GoogleFonts,
			"max_wait_ms":            v.MaxWaitMS,
			"ms_delay":               v.MSDelay,
			"render_when_ready":      v.RenderWhenReady,
			"selector":               v.Selector,
			"viewport_height":        v.ViewportHeight,
			"viewport_width":         v.ViewportWidth,
			"disable_twemoji":        v.DisableTwemoji,
			"color_scheme":           v.ColorScheme,
			"timezone":               v.Timezone,
			"viewport_mobile":        v.ViewportMobile,
			"viewport_landscape":     v.ViewportLandscape,
			"viewport_touch":         v.ViewportTouch,
			"media_type":             v.MediaType,
			"proxy_id":               v.ProxyID,
			"storage_destination_id": v.StorageDestinationID,
			"jumbo_max_height":       v.JumboMaxHeight,
			"jumbo_max_width":        v.JumboMaxWidth,
			"transparent_background": v.TransparentBackground,
		}
	case "html_css":
		return map[string]any{
			"request_overrides":      v.RequestOverrides,
			"html":                   v.HTML,
			"css":                    v.CSS,
			"device_scale":           v.DeviceScale,
			"google_fonts":           v.GoogleFonts,
			"max_wait_ms":            v.MaxWaitMS,
			"metadata":               v.Metadata,
			"ms_delay":               v.MSDelay,
			"render_when_ready":      v.RenderWhenReady,
			"selector":               v.Selector,
			"viewport_height":        v.ViewportHeight,
			"viewport_width":         v.ViewportWidth,
			"pdf_options":            v.PDFOptions,
			"disable_twemoji":        v.DisableTwemoji,
			"max_render_once":        v.MaxRenderOnce,
			"color_scheme":           v.ColorScheme,
			"timezone":               v.Timezone,
			"viewport_mobile":        v.ViewportMobile,
			"viewport_landscape":     v.ViewportLandscape,
			"viewport_touch":         v.ViewportTouch,
			"media_type":             v.MediaType,
			"proxy_id":               v.ProxyID,
			"storage_destination_id": v.StorageDestinationID,
			"jumbo_max_height":       v.JumboMaxHeight,
			"jumbo_max_width":        v.JumboMaxWidth,
			"transparent_background": v.TransparentBackground,
			"format":                 v.Format,
			"dedupe_duration_s":      0,
		}
	case "url":
		return map[string]any{
			"request_overrides":              v.RequestOverrides,
			"url":                            v.URL,
			"css":                            v.CSS,
			"device_scale":                   v.DeviceScale,
			"full_screen":                    v.FullScreen,
			"max_wait_ms":                    v.MaxWaitMS,
			"metadata":                       v.Metadata,
			"ms_delay":                       v.MSDelay,
			"render_when_ready":              v.RenderWhenReady,
			"selector":                       v.Selector,
			"viewport_height":                v.ViewportHeight,
			"viewport_width":                 v.ViewportWidth,
			"pdf_options":                    v.PDFOptions,
			"disable_twemoji":                v.DisableTwemoji,
			"max_render_once":                v.MaxRenderOnce,
			"color_scheme":                   v.ColorScheme,
			"timezone":                       v.Timezone,
			"block_consent_banners":          v.BlockConsentBanners,
			"identify_as_hcti":               v.IdentifyAsHCTI,
			"headers":                        v.Headers,
			"additional_header_origins":      v.AdditionalHeaderOrigins,
			"include_headers_on_subrequests": v.IncludeHeadersOnSubrequests,
			"viewport_mobile":                v.ViewportMobile,
			"viewport_landscape":             v.ViewportLandscape,
			"viewport_touch":                 v.ViewportTouch,
			"media_type":                     v.MediaType,
			"proxy_id":                       v.ProxyID,
			"storage_destination_id":         v.StorageDestinationID,
			"jumbo_max_height":               v.JumboMaxHeight,
			"jumbo_max_width":                v.JumboMaxWidth,
			"transparent_background":         v.TransparentBackground,
			"format":                         v.Format,
			"dedupe_duration_s":              0,
		}
	case "templated":
		return map[string]any{
			"template_values":   v.TemplateValues,
			"format":            v.Format,
			"dedupe_duration_s": 0,
		}
	}
	panic("unsupported rendering kind")
}
