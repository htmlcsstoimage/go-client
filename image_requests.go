package hcti

// ImageFormat selects PNG, JPG, WebP, or PDF rendering URLs.
type ImageFormat string

const (
	PNG  ImageFormat = "png"
	JPG  ImageFormat = "jpg"
	WebP ImageFormat = "webp"
	PDF  ImageFormat = "pdf"
)

// ImageOptions configures HTML/CSS and URL image creation.
type ImageOptions struct {
	RenderOptions
	// Format selects the extension of the initially returned image URL.
	// It does not restrict the stored image to that format.
	// An empty value uses the API's default image URL.
	Format ImageFormat `json:"format,omitempty"`
	// Selector is a CSS selector identifying the element to capture.
	// The image is cropped to the selected element's dimensions.
	// Nil omits the selector; Ptr("") clears an inherited batch selector.
	Selector *string `json:"selector,omitempty"`
	// MaxRenderOnce requests that the image be rendered and saved only once.
	// Nil uses the API default.
	MaxRenderOnce *bool `json:"max_render_once,omitempty"`
	// DedupeDurationSeconds reuses an identical recently created image without
	// consuming another image credit. Ptr(0) disables deduplication.
	// HTML/CSS defaults vary by plan; URL requests default to 0.
	// Applies only to standard single-image POST creation, not batches or signed URLs.
	// See https://docs.htmlcsstoimage.com/parameters/dedupe_duration_s/.
	DedupeDurationSeconds *int `json:"dedupe_duration_s,omitempty"`
	// PDFOptions configures PDF output, including page dimensions and margins.
	// Nil omits PDF options. Select PDF as the Format to receive a PDF URL.
	PDFOptions *PDFOptions `json:"pdf_options,omitempty"`
}

// ImageRequest accepts one of HTMLImageRequest, URLImageRequest, or TemplatedImageRequest.
type ImageRequest interface{ imageRequest() }

// HTMLImageRequest creates an image from HTML markup and optional CSS.
type HTMLImageRequest struct {
	ImageOptions
	// HTML is the raw HTML markup to render. Required for single-image creation;
	// an empty value is omitted so batch variations can inherit default HTML.
	HTML string `json:"html,omitempty"`
	// CSS supplies styles for the HTML or styles to inject into a URL screenshot.
	// Nil omits the field; Ptr("") sends an explicit empty value.
	CSS *string `json:"css,omitempty"`
	// GoogleFonts lists font families to load, such as GoogleFonts{"Open Sans", "Roboto"}.
	// Use font-family in your CSS to select them. Names are trimmed, deduplicated,
	// and serialized to the API's pipe-delimited string. Nil omits the field; an empty list sends an empty string.
	GoogleFonts GoogleFonts `json:"google_fonts,omitempty"`
}

func (HTMLImageRequest) imageRequest() {}

// URLImageRequest creates a screenshot of a webpage.
type URLImageRequest struct {
	ImageOptions
	// URL is the webpage to capture. Required for single-image creation;
	// an empty value is omitted so batch variations can inherit the default URL.
	URL string `json:"url,omitempty"`
	// CSS supplies styles for the HTML or styles to inject into a URL screenshot.
	// Nil omits the field; Ptr("") sends an explicit empty value.
	CSS *string `json:"css,omitempty"`
	// Headers supplies custom HTTP headers for top-level requests to the URL's
	// origin and AdditionalHeaderOrigins. Subrequests require IncludeHeadersOnSubrequests.
	// Nil inherits batch defaults; an empty map clears inherited headers.
	// See https://docs.htmlcsstoimage.com/parameters/headers/.
	Headers map[string]string `json:"headers,omitempty"`
	// AdditionalHeaderOrigins lists extra exact HTTP or HTTPS origins allowed
	// to receive Headers. Each origin includes its scheme, host, and optional port.
	// Nil inherits batch defaults; an empty slice clears inherited origins.
	AdditionalHeaderOrigins []string `json:"additional_header_origins,omitempty"`
	// IncludeHeadersOnSubrequests also sends Headers on subrequests to allowed origins.
	// Nil uses the API default.
	IncludeHeadersOnSubrequests *bool `json:"include_headers_on_subrequests,omitempty"`
	// IdentifyAsHCTI adds X-HCTI-SCREENSHOT: 1 to the top-level page request.
	// Nil uses the API default.
	IdentifyAsHCTI *bool `json:"identify_as_hcti,omitempty"`
	// FullScreen captures the webpage's full height instead of only its viewport.
	// Nil uses the API default.
	FullScreen *bool `json:"full_screen,omitempty"`
	// BlockConsentBanners attempts to block cookie and consent banners.
	// Nil uses the API default.
	BlockConsentBanners *bool `json:"block_consent_banners,omitempty"`
}

func (URLImageRequest) imageRequest() {}

// TemplatedImageRequest renders a saved template using variable values.
type TemplatedImageRequest struct {
	// TemplateID identifies the template used to render the image. Required.
	TemplateID string `json:"template_id"`
	// TemplateVersion pins a specific template version.
	// Nil uses the template's latest version.
	TemplateVersion *int64 `json:"template_version,omitempty"`
	// TemplateValues maps template variable names to JSON-serializable values.
	// Values can include strings, numbers, booleans, arrays, and objects.
	TemplateValues map[string]any `json:"template_values"`
	// Format selects the extension of the initially returned image URL.
	// It does not restrict the stored image to that format.
	// An empty value uses the API's default image URL.
	Format ImageFormat `json:"format,omitempty"`
}

func (TemplatedImageRequest) imageRequest() {}
