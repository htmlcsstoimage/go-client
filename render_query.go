package hcti

import (
	"slices"
)

// appendRenderQuery writes fields in alphabetical order to preserve existing URLs
// without collecting or sorting parameter keys. Format, PDF options, and dedupe
// duration are excluded from the query.
func (r URLImageRequest) appendRenderQuery(b *signedURLBuilder) error {
	for _, origin := range r.AdditionalHeaderOrigins {
		b.safeStringValue("additional_header_origins", origin)
	}
	b.optionalSafeBool("block_consent_banners", r.BlockConsentBanners)
	appendOptionalString(b, "color_scheme", r.ColorScheme)
	appendOptionalString(b, "css", r.CSS)
	if r.DeviceScale != nil {
		b.safeKey("device_scale")
		if err := b.appendJSONFloat(*r.DeviceScale, 64); err != nil {
			return err
		}
	}
	b.optionalSafeBool("disable_twemoji", r.DisableTwemoji)
	b.optionalSafeBool("full_screen", r.FullScreen)
	// Header names are dynamic; sort only these to keep repeated values stable.
	var keyBuffer [16]string
	keys := keyBuffer[:0]
	for name := range r.Headers {
		keys = append(keys, name)
	}
	slices.Sort(keys)
	for _, name := range keys {
		b.safeKey("headers")
		b.buf = appendQueryEscaped(b.buf, name)
		b.buf = append(b.buf, "%3A"...)
		b.buf = appendQueryEscaped(b.buf, r.Headers[name])
	}
	b.optionalSafeBool("identify_as_hcti", r.IdentifyAsHCTI)
	b.optionalSafeBool("include_headers_on_subrequests", r.IncludeHeadersOnSubrequests)
	b.optionalSafeInt("jumbo_max_height", r.JumboMaxHeight)
	b.optionalSafeInt("jumbo_max_width", r.JumboMaxWidth)
	b.optionalSafeBool("max_render_once", r.MaxRenderOnce)
	b.optionalSafeInt("max_wait_ms", r.MaxWaitMS)
	appendOptionalString(b, "media_type", r.MediaType)
	b.optionalSafeInt("ms_delay", r.MSDelay)
	appendOptionalString(b, "proxy_id", r.ProxyID)
	b.optionalSafeBool("render_when_ready", r.RenderWhenReady)
	appendOptionalString(b, "selector", r.Selector)
	appendOptionalString(b, "storage_destination_id", r.StorageDestinationID)
	appendOptionalString(b, "timezone", r.Timezone)
	b.optionalSafeBool("transparent_background", r.TransparentBackground)
	b.safeStringValue("url", r.URL)
	b.optionalSafeInt("viewport_height", r.ViewportHeight)
	b.optionalSafeBool("viewport_landscape", r.ViewportLandscape)
	b.optionalSafeBool("viewport_mobile", r.ViewportMobile)
	b.optionalSafeBool("viewport_touch", r.ViewportTouch)
	b.optionalSafeInt("viewport_width", r.ViewportWidth)
	return nil
}
