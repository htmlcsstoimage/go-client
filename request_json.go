package hcti

import "encoding/json"

// MarshalJSON preserves empty collections, which clear inherited batch defaults.
func (r HTMLImageRequest) MarshalJSON() ([]byte, error) {
	type plain HTMLImageRequest
	var overrides *[]RequestOverride
	if r.RequestOverrides != nil {
		overrides = &r.RequestOverrides
	}
	var fonts *GoogleFonts
	if r.GoogleFonts != nil {
		fonts = &r.GoogleFonts
	}
	return json.Marshal(struct {
		plain
		GoogleFonts      *GoogleFonts       `json:"google_fonts,omitempty"`
		RequestOverrides *[]RequestOverride `json:"request_overrides,omitempty"`
	}{plain(r), fonts, overrides})
}

// MarshalJSON preserves empty collections, which clear inherited batch defaults.
func (r URLImageRequest) MarshalJSON() ([]byte, error) {
	type plain URLImageRequest
	var headers *map[string]string
	var origins *[]string
	var overrides *[]RequestOverride
	if r.Headers != nil {
		headers = &r.Headers
	}
	if r.AdditionalHeaderOrigins != nil {
		origins = &r.AdditionalHeaderOrigins
	}
	if r.RequestOverrides != nil {
		overrides = &r.RequestOverrides
	}
	return json.Marshal(struct {
		plain
		Headers          *map[string]string `json:"headers,omitempty"`
		Origins          *[]string          `json:"additional_header_origins,omitempty"`
		RequestOverrides *[]RequestOverride `json:"request_overrides,omitempty"`
	}{plain(r), headers, origins, overrides})
}

// MarshalJSON excludes deduplication from batch defaults and variations without
// changing the caller's requests. The batch endpoint does not support dedupe.
func (r BatchRequest) MarshalJSON() ([]byte, error) {
	variations := make([]ImageRequest, len(r.Variations))
	for i, request := range r.Variations {
		variations[i] = withoutDedupe(request)
	}
	return json.Marshal(struct {
		Variations []ImageRequest `json:"variations"`
		Defaults   ImageRequest   `json:"default_options,omitempty"`
	}{variations, withoutDedupe(r.DefaultOptions)})
}

func withoutDedupe(request ImageRequest) ImageRequest {
	switch r := request.(type) {
	case *HTMLImageRequest:
		if r != nil {
			copy := *r
			copy.DedupeDurationSeconds = nil
			return &copy
		}
	case *URLImageRequest:
		if r != nil {
			copy := *r
			copy.DedupeDurationSeconds = nil
			return &copy
		}
	}
	return request
}
