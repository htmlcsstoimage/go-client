package management

import hcti "github.com/htmlcsstoimage/go-client"

// OGDefaultImageOptions contains only options supported by HTML/CSS OG configs.
// Nil fields are sent as JSON null, leaving their effective defaults to the API.
type OGDefaultImageOptions struct {
	// CSS CSS injected into the loaded page to override its styles.
	CSS *string `json:"css"`
	// DeviceScale Adjusts the pixel ratio used for the screenshot. Minimum: 0.1. Maximum: 3.
	// HTML and template renders default to 2; URL renders default to 1.
	DeviceScale *float64 `json:"device_scale"`
	// MaxWaitMS Sets a limit on how long to wait before taking the screenshot when the page
	// continues loading irrelevant content. Minimum: 500. Maximum: 10000 and subject to the
	// account plan limit.
	MaxWaitMS *int32 `json:"max_wait_ms"`
	// MSDelay Adds extra time in milliseconds before taking the screenshot so JavaScript can
	// execute. Minimum: 0. Maximum: 10000.
	MSDelay *int32 `json:"ms_delay"`
	// RenderWhenReady Waits until the page signals that the screenshot is ready. The image fails
	// if the readiness signal is never sent.
	RenderWhenReady *bool `json:"render_when_ready"`
	// Selector A CSS selector for an element in the HTML. We’ll crop the image to this specific
	// element.
	Selector *string `json:"selector"`
	// ViewportHeight Sets the height of Chrome's viewport and disables automatic cropping.
	// Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together.
	ViewportHeight *int32 `json:"viewport_height"`
	// ViewportWidth Sets the width of Chrome's viewport and disables automatic cropping.
	// Minimum: 1. Maximum: 6000. Both viewport dimensions must be supplied together.
	ViewportWidth *int32 `json:"viewport_width"`
	// DisableTwemoji Disables the Twemoji fallback and renders emoji using native fonts instead.
	DisableTwemoji *bool `json:"disable_twemoji"`
	// ColorScheme Sets Chrome's preferred color scheme. Valid values: light or dark.
	ColorScheme *hcti.ColorScheme `json:"color_scheme"`
	// Timezone Sets the IANA timezone used by Chrome while rendering. Must be a recognized IANA
	// timezone identifier.
	Timezone *string `json:"timezone"`
	// BlockConsentBanners Attempt to block cookie/consent banners from displaying.
	BlockConsentBanners *bool `json:"block_consent_banners"`
	// IdentifyAsHCTI Identify the top-level page navigation as an HCTI screenshot request using
	// the X-HCTI-SCREENSHOT header.
	IdentifyAsHCTI *bool `json:"identify_as_hcti"`
	// Headers HTTP headers to include on top-level page navigations to the requested URL's
	// origin and any additional_header_origins. Supports up to 20 headers with names up to 512
	// ASCII characters and values up to 8192 UTF-8 bytes. For GET and form-encoded requests,
	// repeat this parameter using the format `headers=name:value`.
	Headers map[string]string `json:"headers"`
	// AdditionalHeaderOrigins Additional exact HTTP or HTTPS origins allowed to receive custom
	// headers. Supports up to 20 unique origins of up to 512 UTF-8 bytes each. Origins must use
	// the format scheme://host[:port] without a path; duplicates are ignored. For GET and form-
	// encoded requests, repeat this parameter for each origin.
	AdditionalHeaderOrigins []string `json:"additional_header_origins"`
	// IncludeHeadersOnSubrequests Include custom headers on subrequests to the requested URL's
	// origin and any additional_header_origins. Defaults to false. Requires at least one header.
	IncludeHeadersOnSubrequests *bool `json:"include_headers_on_subrequests"`
	// ViewportMobile Specifies whether the page uses mobile viewport behavior, including its
	// viewport meta tag.
	ViewportMobile *bool `json:"viewport_mobile"`
	// ViewportLandscape Specifies whether the emulated viewport is in landscape orientation.
	ViewportLandscape *bool `json:"viewport_landscape"`
	// ViewportTouch Specifies whether the emulated viewport supports touch events.
	ViewportTouch *bool `json:"viewport_touch"`
	// MediaType Sets the CSS media type used while rendering the page. Valid values: print or
	// screen.
	MediaType *hcti.MediaType `json:"media_type"`
	// ProxyID Selects the configured proxy id.
	ProxyID *string `json:"proxy_id"`
	// StorageDestinationID Selects the configured storage destination id.
	StorageDestinationID *string `json:"storage_destination_id"`
	// TransparentBackground Specifies whether the image is rendered with a transparent
	// background.
	TransparentBackground *bool `json:"transparent_background"`
}
