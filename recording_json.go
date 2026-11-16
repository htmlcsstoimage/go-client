package hcti

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// MarshalJSON writes a JSON number for unitless distances and a unit-suffixed string otherwise.
func (d ScrollDistance) MarshalJSON() ([]byte, error) {
	if math.IsNaN(d.Value) || math.IsInf(d.Value, 0) || d.Value <= 0 {
		return nil, fmt.Errorf("hcti: scroll distance must be finite and positive")
	}
	switch d.Unit {
	case "":
		return json.Marshal(d.Value)
	case ScrollPixels, ScrollViewportHeight:
		return json.Marshal(strconv.FormatFloat(d.Value, 'f', -1, 64) + string(d.Unit))
	default:
		return nil, fmt.Errorf("hcti: unsupported scroll distance unit %q", d.Unit)
	}
}

// UnmarshalJSON decodes statuses this client does not recognize as RecordingStatusUnknown.
func (s *RecordingStatus) UnmarshalJSON(data []byte) error {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if value == nil || *value == "" {
		*s = ""
		return nil
	}
	switch status := RecordingStatus(*value); status {
	case RecordingStatusProcessing, RecordingStatusPendingRetry, RecordingStatusReady, RecordingStatusFailed, RecordingStatusCancelled, RecordingStatusAwaitingRender:
		*s = status
	default:
		*s = RecordingStatusUnknown
	}
	return nil
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingScrollAction) MarshalJSON() ([]byte, error) {
	type plain RecordingScrollAction
	var steps *[]ScrollStep
	if r.Steps != nil {
		steps = &r.Steps
	}
	var ifNotFound *[]RecordingAction
	if r.IfNotFound != nil {
		ifNotFound = &r.IfNotFound
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
		Steps      *[]ScrollStep      `json:"steps,omitempty"`
		IfNotFound *[]RecordingAction `json:"if_not_found,omitempty"`
	}{
		Type:       "scroll",
		plain:      plain(r),
		Steps:      steps,
		IfNotFound: ifNotFound,
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingStallAction) MarshalJSON() ([]byte, error) {
	type plain RecordingStallAction
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "stall",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingLocateSelectorAction) MarshalJSON() ([]byte, error) {
	type plain RecordingLocateSelectorAction
	var ifFound *[]RecordingAction
	if r.IfFound != nil {
		ifFound = &r.IfFound
	}
	var ifNotFound *[]RecordingAction
	if r.IfNotFound != nil {
		ifNotFound = &r.IfNotFound
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
		IfFound    *[]RecordingAction `json:"if_found,omitempty"`
		IfNotFound *[]RecordingAction `json:"if_not_found,omitempty"`
	}{
		Type:       "locate_selector",
		plain:      plain(r),
		IfFound:    ifFound,
		IfNotFound: ifNotFound,
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingClickAction) MarshalJSON() ([]byte, error) {
	type plain RecordingClickAction
	var ifNotFound *[]RecordingAction
	if r.IfNotFound != nil {
		ifNotFound = &r.IfNotFound
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
		IfNotFound *[]RecordingAction `json:"if_not_found,omitempty"`
	}{
		Type:       "click",
		plain:      plain(r),
		IfNotFound: ifNotFound,
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingErrorAction) MarshalJSON() ([]byte, error) {
	type plain RecordingErrorAction
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "error",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingSimpleScrollAction) MarshalJSON() ([]byte, error) {
	type plain RecordingSimpleScrollAction
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "simple_scroll",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingExecuteJavaScriptAction) MarshalJSON() ([]byte, error) {
	type plain RecordingExecuteJavaScriptAction
	var ifError *[]RecordingAction
	if r.IfError != nil {
		ifError = &r.IfError
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
		IfError *[]RecordingAction `json:"if_error,omitempty"`
	}{
		Type:    "execute_javascript",
		plain:   plain(r),
		IfError: ifError,
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingImmediateTrigger) MarshalJSON() ([]byte, error) {
	type plain RecordingImmediateTrigger
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "immediate",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingElementVisibleTrigger) MarshalJSON() ([]byte, error) {
	type plain RecordingElementVisibleTrigger
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "element_visible",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r RecordingFunctionCallTrigger) MarshalJSON() ([]byte, error) {
	type plain RecordingFunctionCallTrigger
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "function_call",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r PositionScrollTarget) MarshalJSON() ([]byte, error) {
	type plain PositionScrollTarget
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "position",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r CoordinatesScrollTarget) MarshalJSON() ([]byte, error) {
	type plain CoordinatesScrollTarget
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
	}{
		Type:  "coordinates",
		plain: plain(r),
	})
}

// MarshalJSON adds the API discriminator and preserves empty collections.
func (r SelectorScrollTarget) MarshalJSON() ([]byte, error) {
	type plain SelectorScrollTarget
	var ifNotFound *[]RecordingAction
	if r.IfNotFound != nil {
		ifNotFound = &r.IfNotFound
	}
	return json.Marshal(struct {
		Type string `json:"type"`
		plain
		IfNotFound *[]RecordingAction `json:"if_not_found,omitempty"`
	}{
		Type:       "selector",
		plain:      plain(r),
		IfNotFound: ifNotFound,
	})
}

// MarshalJSON preserves empty collections.
func (r HTMLRecordingRequest) MarshalJSON() ([]byte, error) {
	type plain HTMLRecordingRequest
	var googleFonts *GoogleFonts
	if r.GoogleFonts != nil {
		googleFonts = &r.GoogleFonts
	}
	var requestOverrides *[]RequestOverride
	if r.RequestOverrides != nil {
		requestOverrides = &r.RequestOverrides
	}
	var actions *[]RecordingAction
	if r.Actions != nil {
		actions = &r.Actions
	}
	return json.Marshal(struct {
		plain
		GoogleFonts      *GoogleFonts       `json:"google_fonts,omitempty"`
		RequestOverrides *[]RequestOverride `json:"request_overrides,omitempty"`
		Actions          *[]RecordingAction `json:"actions,omitempty"`
	}{
		plain:            plain(r),
		GoogleFonts:      googleFonts,
		RequestOverrides: requestOverrides,
		Actions:          actions,
	})
}

// MarshalJSON preserves empty collections.
func (r URLRecordingRequest) MarshalJSON() ([]byte, error) {
	type plain URLRecordingRequest
	var headers *map[string]string
	if r.Headers != nil {
		headers = &r.Headers
	}
	var additionalHeaderOrigins *[]string
	if r.AdditionalHeaderOrigins != nil {
		additionalHeaderOrigins = &r.AdditionalHeaderOrigins
	}
	var requestOverrides *[]RequestOverride
	if r.RequestOverrides != nil {
		requestOverrides = &r.RequestOverrides
	}
	var actions *[]RecordingAction
	if r.Actions != nil {
		actions = &r.Actions
	}
	return json.Marshal(struct {
		plain
		Headers                 *map[string]string `json:"headers,omitempty"`
		AdditionalHeaderOrigins *[]string          `json:"additional_header_origins,omitempty"`
		RequestOverrides        *[]RequestOverride `json:"request_overrides,omitempty"`
		Actions                 *[]RecordingAction `json:"actions,omitempty"`
	}{
		plain:                   plain(r),
		Headers:                 headers,
		AdditionalHeaderOrigins: additionalHeaderOrigins,
		RequestOverrides:        requestOverrides,
		Actions:                 actions,
	})
}
