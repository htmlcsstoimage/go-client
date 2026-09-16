package management

import "encoding/json"

// TemplatedOGConfigRequest configures rendering with values extracted from a page.
type TemplatedOGConfigRequest struct {
	OGConfigOptions
	// TemplateID selects the template, including its t- prefix.
	TemplateID string `json:"template_id"`
	// TemplateVersion is optional. Nil follows the latest version.
	TemplateVersion *int64 `json:"template_version"`
	// TemplateValuesMapping contains up to 32 ordered metadata mappings.
	TemplateValuesMapping []OGTemplateValueMapping `json:"template_values_mapping"`
	// Headers are sent to the source origin and AdditionalHeaderOrigins.
	Headers map[string]string `json:"headers"`
	// AdditionalHeaderOrigins lists up to 20 additional exact HTTP(S) origins allowed to receive Headers.
	AdditionalHeaderOrigins []string `json:"additional_header_origins"`
}

func (*TemplatedOGConfigRequest) ogConfigRequest() {}

// MarshalJSON includes the templated discriminator.
func (v TemplatedOGConfigRequest) MarshalJSON() ([]byte, error) {
	type fields TemplatedOGConfigRequest
	return json.Marshal(struct {
		ConfigType OGConfigType `json:"config_type"`
		fields
	}{OGConfigTemplated, fields(v)})
}

// OGMetadataFallback selects a standard page metadata fallback.
type OGMetadataFallback string

const (
	// Titles uses the page's title fallbacks.
	Titles OGMetadataFallback = "titles"
	// Descriptions uses the page's description fallbacks.
	Descriptions OGMetadataFallback = "descriptions"
)

// OGTemplateValueMapping maps either MetaKey or Fallback to a template field.
type OGTemplateValueMapping struct {
	// TemplateKey is the destination template field, up to 128 characters.
	TemplateKey string `json:"template_key"`
	// MetaKey names the source metadata field. Supply exactly one of MetaKey and Fallback.
	MetaKey *string `json:"meta_key"`
	// Fallback selects standard titles or descriptions instead of an explicit MetaKey.
	Fallback *OGMetadataFallback `json:"fallback"`
}
