package hcti

// ColorScheme selects Chrome's preferred light or dark appearance.
type ColorScheme string

const (
	Light ColorScheme = "light"
	Dark  ColorScheme = "dark"
)

// MediaType selects which CSS media rules Chrome uses.
type MediaType string

const (
	Screen MediaType = "screen"
	Print  MediaType = "print"
)

// RenderOptions are shared by HTML, URL, and HTML template requests.
// Nil pointers leave the API default intact; Ptr(false) and Ptr(0) send explicit values.
type RenderOptions struct {
	// DeviceScale controls the screenshot pixel ratio, from 0.1 to 3.
	// HTML and template renders default to 2; URL renders default to 1.
	// Nil uses the API default.
	DeviceScale *float64 `json:"device_scale,omitempty"`
	// ViewportHeight sets Chrome's viewport height in pixels, from 1 to 6000.
	// Supply ViewportWidth as well. Setting viewport dimensions disables automatic cropping.
	ViewportHeight *int `json:"viewport_height,omitempty"`
	// ViewportWidth sets Chrome's viewport width in pixels, from 1 to 6000.
	// Supply ViewportHeight as well. Nil leaves the API default unchanged.
	ViewportWidth *int `json:"viewport_width,omitempty"`
	// MaxWaitMS limits how long to wait before taking the screenshot when
	// the page continues loading irrelevant content. Values are milliseconds,
	// from 500 to 10000, and are also subject to the account plan limit.
	MaxWaitMS *int `json:"max_wait_ms,omitempty"`
	// MSDelay adds a delay before taking the screenshot so JavaScript can execute.
	// Values are milliseconds, from 0 to 10000. The API default is 0.
	MSDelay *int `json:"ms_delay,omitempty"`
	// RenderWhenReady waits for JavaScript to call ScreenshotReady().
	// The image fails if the readiness signal is never sent. Nil uses the API default.
	RenderWhenReady *bool `json:"render_when_ready,omitempty"`
	// DisableTwemoji disables the Twemoji fallback and renders emoji with native fonts.
	// Nil uses the API default; Ptr(false) explicitly enables the fallback.
	DisableTwemoji *bool `json:"disable_twemoji,omitempty"`
	// ColorScheme sets Chrome's preferred color scheme to Light or Dark.
	// Nil leaves the API default unchanged; use Ptr(Dark) or Ptr(Light).
	ColorScheme *ColorScheme `json:"color_scheme,omitempty"`
	// Timezone sets Chrome's timezone using an IANA name, such as America/New_York.
	// Nil leaves the API default unchanged; Ptr("") sends an explicit empty value.
	Timezone *string `json:"timezone,omitempty"`
	// ViewportMobile enables mobile viewport behavior, including the page's
	// viewport meta tag. Nil uses the API default.
	ViewportMobile *bool `json:"viewport_mobile,omitempty"`
	// ViewportTouch enables touch interactions in the emulated viewport.
	// Nil uses the API default.
	ViewportTouch *bool `json:"viewport_touch,omitempty"`
	// ViewportLandscape sets the emulated viewport to landscape orientation.
	// Nil uses the API default.
	ViewportLandscape *bool `json:"viewport_landscape,omitempty"`
	// MediaType selects Screen or Print CSS media rules.
	// Nil leaves the API default unchanged; use Ptr(Screen) or Ptr(Print).
	MediaType *MediaType `json:"media_type,omitempty"`
	// ProxyID selects an organization proxy for rendering.
	// Nil omits the proxy selection; Ptr("") sends an explicit empty value.
	// See https://docs.htmlcsstoimage.com/parameters/proxy_id/.
	ProxyID *string `json:"proxy_id,omitempty"`
	// StorageDestinationID selects an organization storage destination for rendered files.
	// Nil omits the destination selection; Ptr("") sends an explicit empty value.
	// See https://docs.htmlcsstoimage.com/parameters/storage_destination_id/.
	StorageDestinationID *string `json:"storage_destination_id,omitempty"`
	// JumboMaxWidth sets the maximum width in pixels for jumbo rendering.
	// Supply JumboMaxHeight as well. Both must be positive and at most 80000;
	// at least one must exceed 8000, and their product must not exceed 400000000.
	// Jumbo rendering consumes additional renders.
	JumboMaxWidth *int `json:"jumbo_max_width,omitempty"`
	// JumboMaxHeight sets the maximum height in pixels for jumbo rendering.
	// Supply JumboMaxWidth as well; its documentation describes the size limits.
	// Jumbo rendering consumes additional renders.
	JumboMaxHeight *int `json:"jumbo_max_height,omitempty"`
	// TransparentBackground requests a transparent image background.
	// Nil uses the API default; Ptr(false) explicitly disables transparency.
	TransparentBackground *bool `json:"transparent_background,omitempty"`
}
