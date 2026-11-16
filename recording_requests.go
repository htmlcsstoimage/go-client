package hcti

// RecordingFormat selects recording format values.
type RecordingFormat string

const (
	RecordingFormatMP4  RecordingFormat = "mp4"
	RecordingFormatWebM RecordingFormat = "webm"
	RecordingFormatGIF  RecordingFormat = "gif"
)

// RecordingDurationMode controls how the requested duration limits the timeline.
// Fixed pads or cuts off; Target adjusts playback speed; AtLeast lets actions finish
// up to 90 seconds; AtMost stops when actions finish or the duration is reached.
type RecordingDurationMode string

const (
	RecordingDurationModeFixed   RecordingDurationMode = "fixed"
	RecordingDurationModeTarget  RecordingDurationMode = "target"
	RecordingDurationModeAtLeast RecordingDurationMode = "at_least"
	RecordingDurationModeAtMost  RecordingDurationMode = "at_most"
)

// RecordingMouseButton selects which mouse button a click action uses.
type RecordingMouseButton string

const (
	RecordingMouseButtonLeft    RecordingMouseButton = "left"
	RecordingMouseButtonRight   RecordingMouseButton = "right"
	RecordingMouseButtonMiddle  RecordingMouseButton = "middle"
	RecordingMouseButtonBack    RecordingMouseButton = "back"
	RecordingMouseButtonForward RecordingMouseButton = "forward"
)

// RecordingStatus is the current lifecycle state. Ready, Failed, and Cancelled are terminal.
// Statuses this client does not recognize decode as RecordingStatusUnknown.
type RecordingStatus string

// RecordingStatusUnknown is an unrecognized status, such as one added by a newer API version.
const RecordingStatusUnknown RecordingStatus = "unknown"

const (
	RecordingStatusProcessing     RecordingStatus = "processing"
	RecordingStatusPendingRetry   RecordingStatus = "pending_retry"
	RecordingStatusReady          RecordingStatus = "ready"
	RecordingStatusFailed         RecordingStatus = "failed"
	RecordingStatusCancelled      RecordingStatus = "cancelled"
	RecordingStatusAwaitingRender RecordingStatus = "awaiting_render"
)

// ScrollEasing selects scroll easing values.
type ScrollEasing string

const (
	ScrollEasingLinear           ScrollEasing = "linear"
	ScrollEasingEaseInCubic      ScrollEasing = "ease_in_cubic"
	ScrollEasingEaseOutCubic     ScrollEasing = "ease_out_cubic"
	ScrollEasingEaseInOutCubic   ScrollEasing = "ease_in_out_cubic"
	ScrollEasingEaseInQuad       ScrollEasing = "ease_in_quad"
	ScrollEasingEaseOutQuad      ScrollEasing = "ease_out_quad"
	ScrollEasingEaseInOutQuad    ScrollEasing = "ease_in_out_quad"
	ScrollEasingEaseInQuart      ScrollEasing = "ease_in_quart"
	ScrollEasingEaseOutQuart     ScrollEasing = "ease_out_quart"
	ScrollEasingEaseInOutQuart   ScrollEasing = "ease_in_out_quart"
	ScrollEasingEaseInQuint      ScrollEasing = "ease_in_quint"
	ScrollEasingEaseOutQuint     ScrollEasing = "ease_out_quint"
	ScrollEasingEaseInOutQuint   ScrollEasing = "ease_in_out_quint"
	ScrollEasingEaseInSine       ScrollEasing = "ease_in_sine"
	ScrollEasingEaseOutSine      ScrollEasing = "ease_out_sine"
	ScrollEasingEaseInOutSine    ScrollEasing = "ease_in_out_sine"
	ScrollEasingEaseInExpo       ScrollEasing = "ease_in_expo"
	ScrollEasingEaseOutExpo      ScrollEasing = "ease_out_expo"
	ScrollEasingEaseInOutExpo    ScrollEasing = "ease_in_out_expo"
	ScrollEasingEaseInCirc       ScrollEasing = "ease_in_circ"
	ScrollEasingEaseOutCirc      ScrollEasing = "ease_out_circ"
	ScrollEasingEaseInOutCirc    ScrollEasing = "ease_in_out_circ"
	ScrollEasingEaseInBack       ScrollEasing = "ease_in_back"
	ScrollEasingEaseOutBack      ScrollEasing = "ease_out_back"
	ScrollEasingEaseInOutBack    ScrollEasing = "ease_in_out_back"
	ScrollEasingEaseInElastic    ScrollEasing = "ease_in_elastic"
	ScrollEasingEaseOutElastic   ScrollEasing = "ease_out_elastic"
	ScrollEasingEaseInOutElastic ScrollEasing = "ease_in_out_elastic"
	ScrollEasingEaseInBounce     ScrollEasing = "ease_in_bounce"
	ScrollEasingEaseOutBounce    ScrollEasing = "ease_out_bounce"
	ScrollEasingEaseInOutBounce  ScrollEasing = "ease_in_out_bounce"
)

// ScrollPosition selects scroll position values.
type ScrollPosition string

const (
	ScrollPositionTop    ScrollPosition = "top"
	ScrollPositionBottom ScrollPosition = "bottom"
	ScrollPositionLeft   ScrollPosition = "left"
	ScrollPositionRight  ScrollPosition = "right"
)

// ScrollDistance is the distance of one simple scroll movement. Use an empty Unit to send a plain
// number of CSS pixels, ScrollPixels to send "400px", or ScrollViewportHeight to send "50vh".
type ScrollDistance struct {
	// Value is the distance magnitude. It must be finite and positive.
	Value float64
	// Unit selects ScrollPixels or ScrollViewportHeight. An empty unit sends a JSON number of pixels.
	Unit ScrollDistanceUnit
}

// ScrollDistanceUnit is a unit accepted for simple scroll distances.
type ScrollDistanceUnit string

const (
	ScrollPixels         ScrollDistanceUnit = "px"
	ScrollViewportHeight ScrollDistanceUnit = "vh"
)

// RecordingAction is implemented by the supported timeline actions.
type RecordingAction interface{ recordingAction() }

// RecordingTrigger is implemented by the supported recording triggers.
type RecordingTrigger interface{ recordingTrigger() }

// ScrollTarget is implemented by the supported scroll destinations.
type ScrollTarget interface{ scrollTarget() }

// RecordingScrollAction scrolls through a sequence of timed movements.
type RecordingScrollAction struct {
	// ContainerSelector Scrollable container's CSS selector. Omit or use null to scroll the document.
	ContainerSelector *string `json:"container_selector,omitempty"`
	// Steps Ordered movements. Supply at least one step.
	Steps []ScrollStep `json:"steps"`
	// IfNotFound Fallback when the container is missing. Omitted, null, or empty skips this action.
	IfNotFound []RecordingAction `json:"if_not_found,omitempty"`
}

func (RecordingScrollAction) recordingAction() {}

// RecordingStallAction waits while capture and page activity continue.
type RecordingStallAction struct {
	// DurationMS Wait in milliseconds. Required; at least 1. The recording deadline can end it early.
	DurationMS int `json:"duration_ms"`
}

func (RecordingStallAction) recordingAction() {}

// RecordingLocateSelectorAction finds an element and runs the corresponding branch.
type RecordingLocateSelectorAction struct {
	// Selector CSS selector to find.
	Selector string `json:"selector"`
	// TimeoutMS Maximum wait in milliseconds. Defaults to 0 for a single check.
	TimeoutMS *int `json:"timeout_ms,omitempty"`
	// IfFound Actions when found. Omitted, null, or empty continues with the parent sequence.
	IfFound []RecordingAction `json:"if_found,omitempty"`
	// IfNotFound Actions when missing. Include an error action to fail the recording.
	IfNotFound []RecordingAction `json:"if_not_found,omitempty"`
}

func (RecordingLocateSelectorAction) recordingAction() {}

// RecordingClickAction finds and clicks an element, scrolling it into view if needed.
type RecordingClickAction struct {
	// Selector CSS selector of the element to click.
	Selector string `json:"selector"`
	// TimeoutMS Maximum wait in milliseconds. Defaults to 0 for a single check.
	TimeoutMS *int `json:"timeout_ms,omitempty"`
	// IfNotFound Fallback when missing. Omitted, null, or empty continues; use an error action to
	// fail.
	IfNotFound []RecordingAction `json:"if_not_found,omitempty"`
	// Button Mouse button. Defaults to left.
	Button *RecordingMouseButton `json:"button,omitempty"`
	// Count Event click count, from 1 to 3. Defaults to 1; this does not repeat the action.
	Count *int `json:"count,omitempty"`
	// DelayMS Delay between mouse down and up, in milliseconds. Defaults to 0; maximum 5000.
	DelayMS *int `json:"delay_ms,omitempty"`
	// OffsetX CSS pixels from the element's top-left. Supplying either offset makes the other default
	// to 0.
	OffsetX *float64 `json:"offset_x,omitempty"`
	// OffsetY CSS pixels from the element's top-left. Omit both offsets to use the default clickable
	// point.
	OffsetY *float64 `json:"offset_y,omitempty"`
}

func (RecordingClickAction) recordingAction() {}

// RecordingErrorAction fails the recording. It must be the last action in its array.
type RecordingErrorAction struct {
}

func (RecordingErrorAction) recordingAction() {}

// RecordingSimpleScrollAction scrolls down repeatedly until the bottom or recording deadline.
type RecordingSimpleScrollAction struct {
	// ScrollDurationMS Duration of each movement in milliseconds. Defaults to 800.
	ScrollDurationMS *int `json:"scroll_duration_ms,omitempty"`
	// MSBetweenScrolls Pause between movements in milliseconds. Defaults to 400; 0 scrolls
	// continuously.
	MSBetweenScrolls *int `json:"ms_between_scrolls,omitempty"`
	// ScrollDistance Distance of each movement. Defaults to 90vh.
	ScrollDistance *ScrollDistance `json:"scroll_distance,omitempty"`
	// Easing Movement easing. Defaults to ease_in_out_expo.
	Easing *ScrollEasing `json:"easing,omitempty"`
	// StallAtEndMS Time to hold at the bottom in milliseconds, limited by remaining recording time.
	// Defaults to 0.
	StallAtEndMS *int `json:"stall_at_end_ms,omitempty"`
}

func (RecordingSimpleScrollAction) recordingAction() {}

// RecordingExecuteJavaScriptAction executes an async function body in the page while recording.
type RecordingExecuteJavaScriptAction struct {
	// Script Function body supporting await. Its return value is ignored.
	Script string `json:"script"`
	// TimeoutMS Maximum execution time in milliseconds. Zero or omitted uses the page timeout; maximum
	// 90000.
	TimeoutMS *int `json:"timeout_ms,omitempty"`
	// IfError Omitted or null fails on script errors; [] ignores errors. A timeout does not undo side
	// effects.
	IfError []RecordingAction `json:"if_error,omitempty"`
}

func (RecordingExecuteJavaScriptAction) recordingAction() {}

// RecordingImmediateTrigger starts capture immediately after page preparation.
type RecordingImmediateTrigger struct {
}

func (RecordingImmediateTrigger) recordingTrigger() {}

// RecordingElementVisibleTrigger waits for an element to become visible before capture.
type RecordingElementVisibleTrigger struct {
	// Selector CSS selector for the element to wait for.
	Selector string `json:"selector"`
	// TimeoutMS Maximum wait in milliseconds. Defaults to 1000; minimum 100. Plan limits apply.
	TimeoutMS *int `json:"timeout_ms,omitempty"`
}

func (RecordingElementVisibleTrigger) recordingTrigger() {}

// RecordingFunctionCallTrigger waits for the page to call window.RecordingReady().
type RecordingFunctionCallTrigger struct {
	// TimeoutMS Maximum wait in milliseconds. Defaults to 1000; minimum 100. Plan limits apply.
	TimeoutMS *int `json:"timeout_ms,omitempty"`
}

func (RecordingFunctionCallTrigger) recordingTrigger() {}

// PositionScrollTarget scrolls to an edge of the document or selected container.
type PositionScrollTarget struct {
	// Position Edge to scroll to. Defaults to top.
	Position *ScrollPosition `json:"position,omitempty"`
}

func (PositionScrollTarget) scrollTarget() {}

// CoordinatesScrollTarget scrolls to absolute CSS pixel coordinates.
type CoordinatesScrollTarget struct {
	// X Horizontal position. Defaults to 0.
	X *float64 `json:"x,omitempty"`
	// Y Vertical position. Defaults to 0.
	Y *float64 `json:"y,omitempty"`
}

func (CoordinatesScrollTarget) scrollTarget() {}

// SelectorScrollTarget scrolls to a matching element with optional offsets.
type SelectorScrollTarget struct {
	// Selector CSS selector for the destination element.
	Selector string `json:"selector"`
	// OffsetX Horizontal offset from the element. Defaults to 0.
	OffsetX *float64 `json:"offset_x,omitempty"`
	// OffsetY Vertical offset from the element. Defaults to 0.
	OffsetY *float64 `json:"offset_y,omitempty"`
	// IfNotFound Fallback when missing. Omitted, null, or empty skips this movement and continues to
	// the next step.
	IfNotFound []RecordingAction `json:"if_not_found,omitempty"`
}

func (SelectorScrollTarget) scrollTarget() {}

// RecordingOutputOptions sets the format, frame rate, and captured-pixel scaling.
type RecordingOutputOptions struct {
	// Format Media format. Defaults to mp4.
	Format *RecordingFormat `json:"format,omitempty"`
	// FPS Frames per second. Defaults to 15 for GIF and 30 for video. GIF supports 10–50; video
	// supports 10–120.
	FPS *int `json:"fps,omitempty"`
	// Scale Scales captured pixels without changing page layout. Defaults to 1; 0.5 halves both
	// dimensions.
	Scale *float64 `json:"scale,omitempty"`
}

// ScrollStep is one timed movement in a scroll action.
type ScrollStep struct {
	// To Destination for this movement. Select position, coordinates, or selector using its type
	// discriminator.
	To ScrollTarget `json:"to"`
	// DurationMS Movement duration in milliseconds. Required; at least 1. The recording deadline can
	// cut it short.
	DurationMS int `json:"duration_ms"`
	// Easing Movement easing. Defaults to linear.
	Easing *ScrollEasing `json:"easing,omitempty"`
}

// RecordingOptions configures browser settings, duration, timeline, and output. Nil pointers use
// API defaults. See https://docs.htmlcsstoimage.com/getting-started/recordings/.
type RecordingOptions struct {
	// CSS CSS to style HTML content or inject into a URL page after it loads. See
	// https://docs.htmlcsstoimage.com/parameters/css/
	CSS *string `json:"css,omitempty"`
	// DeviceScale Capture pixel density. Defaults to 2 for HTML/CSS and 1 for URL recordings. See
	// https://docs.htmlcsstoimage.com/parameters/device_scale/
	DeviceScale *float64 `json:"device_scale,omitempty"`
	// ViewportHeight Viewport height in CSS pixels. Supply both dimensions to set the recording size.
	// See https://docs.htmlcsstoimage.com/parameters/viewport/
	ViewportHeight *int `json:"viewport_height,omitempty"`
	// ViewportWidth Viewport width in CSS pixels. Supply both dimensions to set the recording size.
	// See https://docs.htmlcsstoimage.com/parameters/viewport/
	ViewportWidth *int `json:"viewport_width,omitempty"`
	// DisableTwemoji Disable Twemoji for HTML/CSS when true. URL recordings inject it only when
	// explicitly false. See https://docs.htmlcsstoimage.com/parameters/disable_twemoji/
	DisableTwemoji *bool `json:"disable_twemoji,omitempty"`
	// DedupeDurationSeconds Reuse an identical recording created within this many seconds. Set to 0 to
	// disable deduplication. See https://docs.htmlcsstoimage.com/parameters/dedupe_duration_s/
	DedupeDurationSeconds *int `json:"dedupe_duration_s,omitempty"`
	// ColorScheme Chrome's preferred light or dark color scheme while recording. See
	// https://docs.htmlcsstoimage.com/parameters/color_scheme/
	ColorScheme *ColorScheme `json:"color_scheme,omitempty"`
	// Timezone Browser timezone as an IANA identifier, such as 'America/New_York'. See
	// https://docs.htmlcsstoimage.com/parameters/timezone/
	Timezone *string `json:"timezone,omitempty"`
	// ViewportMobile Enable mobile viewport behavior, including the page's viewport meta tag. See
	// https://docs.htmlcsstoimage.com/parameters/viewport/
	ViewportMobile *bool `json:"viewport_mobile,omitempty"`
	// ViewportLandscape Emulate landscape orientation. See
	// https://docs.htmlcsstoimage.com/parameters/viewport/
	ViewportLandscape *bool `json:"viewport_landscape,omitempty"`
	// ViewportTouch Emulate a viewport that supports touch events. See
	// https://docs.htmlcsstoimage.com/parameters/viewport/
	ViewportTouch *bool `json:"viewport_touch,omitempty"`
	// MediaType CSS media type used while recording. See
	// https://docs.htmlcsstoimage.com/parameters/media_type/
	MediaType *MediaType `json:"media_type,omitempty"`
	// ProxyID Organization proxy to use while recording. See
	// https://docs.htmlcsstoimage.com/guides/advanced/proxies/
	ProxyID *string `json:"proxy_id,omitempty"`
	// StorageDestinationID Storage destination to receive the recording. See
	// https://docs.htmlcsstoimage.com/guides/advanced/storage-destinations/
	StorageDestinationID *string `json:"storage_destination_id,omitempty"`
	// RequestOverrides Block matching browser requests while recording. Available on paid plans. See
	// https://docs.htmlcsstoimage.com/parameters/request_overrides/
	RequestOverrides []RequestOverride `json:"request_overrides,omitempty"`
	// DurationMS Requested duration in milliseconds: 500–30000, or up to 90000 in target mode.
	// duration_mode controls how actions fit into this duration. See
	// https://docs.htmlcsstoimage.com/parameters/duration_ms/
	DurationMS int `json:"duration_ms"`
	// Trigger When recording and actions start after page preparation. Defaults to immediate. See
	// https://docs.htmlcsstoimage.com/parameters/trigger/
	Trigger RecordingTrigger `json:"trigger,omitempty"`
	// Actions Ordered interactions. Leave null or empty to record the page without interactions. See
	// https://docs.htmlcsstoimage.com/parameters/actions/
	Actions []RecordingAction `json:"actions,omitempty"`
	// Output Output format, frame rate, and scale. Defaults to MP4 at 30 fps without audio. See
	// https://docs.htmlcsstoimage.com/parameters/output/
	Output *RecordingOutputOptions `json:"output,omitempty"`
	// DurationMode How duration_ms limits the timeline. Defaults to at_most. See
	// https://docs.htmlcsstoimage.com/parameters/duration_mode/
	DurationMode *RecordingDurationMode `json:"duration_mode,omitempty"`
	// DeferRender Wait to render until media or storage is requested. Defaults to false; does not
	// reset deduplicated recordings. See https://docs.htmlcsstoimage.com/parameters/defer_render/
	DeferRender *bool `json:"defer_render,omitempty"`
	// SkipScrollPreload Skip the preparation scroll that loads lazy content. Scrolling actions still
	// run. Defaults to false. See https://docs.htmlcsstoimage.com/parameters/skip_scroll_preload/
	SkipScrollPreload *bool `json:"skip_scroll_preload,omitempty"`
}

// RecordingRequest accepts *HTMLRecordingRequest or *URLRecordingRequest.
type RecordingRequest interface{ recordingRequest() }

// HTMLRecordingRequest records an HTML document or fragment.
type HTMLRecordingRequest struct {
	RecordingOptions
	// GoogleFonts Google Fonts to load. Set font-family in your CSS to use them. See
	// https://docs.htmlcsstoimage.com/parameters/google_fonts/
	GoogleFonts GoogleFonts `json:"google_fonts,omitempty"`
	// HTML HTML document or fragment to record. Fragments are rendered in a wrapper document. See
	// https://docs.htmlcsstoimage.com/parameters/html/
	HTML string `json:"html"`
}

func (*HTMLRecordingRequest) recordingRequest() {}

// URLRecordingRequest records a public HTTP or HTTPS webpage.
type URLRecordingRequest struct {
	RecordingOptions
	// URL Public HTTP or HTTPS URL to record. See https://docs.htmlcsstoimage.com/parameters/url/
	URL string `json:"url"`
	// BlockConsentBanners Attempt to hide cookie and consent banners before recording. See
	// https://docs.htmlcsstoimage.com/guides/advanced/blocking-cookie-banners/
	BlockConsentBanners *bool `json:"block_consent_banners,omitempty"`
	// IdentifyAsHCTI Add X-HCTI-SCREENSHOT: 1 to the top-level navigation. See
	// https://docs.htmlcsstoimage.com/parameters/identify_as_hcti/
	IdentifyAsHCTI *bool `json:"identify_as_hcti,omitempty"`
	// Headers Custom headers for top-level navigations to the URL's origin and
	// additional_header_origins. See https://docs.htmlcsstoimage.com/parameters/headers/
	Headers map[string]string `json:"headers,omitempty"`
	// AdditionalHeaderOrigins Additional exact HTTP or HTTPS origins allowed to receive custom
	// headers. See https://docs.htmlcsstoimage.com/parameters/headers/#additional-header-origins
	AdditionalHeaderOrigins []string `json:"additional_header_origins,omitempty"`
	// IncludeHeadersOnSubrequests Send custom headers with subrequests to allowed origins. Defaults to
	// false. See https://docs.htmlcsstoimage.com/parameters/headers/#include-headers-on-subrequests
	IncludeHeadersOnSubrequests *bool `json:"include_headers_on_subrequests,omitempty"`
}

func (*URLRecordingRequest) recordingRequest() {}
