package management

import "time"

// OGConfigType identifies HTML/CSS or template rendering.
type OGConfigType string

const (
	// OGConfigHTMLCSS renders page HTML/CSS.
	OGConfigHTMLCSS OGConfigType = "html_css"
	// OGConfigTemplated renders a template populated from page metadata.
	OGConfigTemplated OGConfigType = "templated"
)

// OptimizationMode controls sizing for social platforms.
type OptimizationMode string

const (
	// NoOptimization preserves original dimensions.
	NoOptimization OptimizationMode = "no_optimization"
	// PostProcess adapts dimensions after rendering.
	PostProcess OptimizationMode = "post_process"
	// SetViewport sizes the viewport before rendering.
	SetViewport OptimizationMode = "set_viewport"
)

// OGConfigRequest accepts a non-nil *HTMLCSSOGConfigRequest or *TemplatedOGConfigRequest.
// Concrete request types supply their config_type automatically.
type OGConfigRequest interface{ ogConfigRequest() }

// OGConfigOptions contains the editable settings shared by both config types.
type OGConfigOptions struct {
	// Name is a display name of up to 255 characters.
	Name string `json:"name"`
	// Description is optional; nil clears it on update.
	Description *string `json:"description"`
	// BaseURL is the source website's HTTPS origin, without a path or credentials.
	BaseURL string `json:"base_url"`
	// OptimizationMode is omitted when nil, leaving the default to the API.
	OptimizationMode *OptimizationMode `json:"optimization_mode,omitempty"`
	// Disabled prevents serving Open Graph images using this configuration.
	Disabled bool `json:"disabled"`
	// RefreshIntervalSeconds is 1800–31536000, subject to the plan minimum.
	// Nil leaves the default to the API; no default is inserted by the SDK.
	RefreshIntervalSeconds *uint32 `json:"refresh_interval_s,omitempty"`
}

// OGConfig is readable configuration metadata. ConfigType determines which
// rendering fields apply: DefaultOptions/ExtractValues for HTML/CSS, or
// TemplateID/TemplateVersion/TemplateValuesMapping/Headers/AdditionalHeaderOrigins for templates.
type OGConfig struct {
	ID                     string                 `json:"id"`
	DomainID               string                 `json:"domain_id"`
	ConfigType             OGConfigType           `json:"config_type"`
	Name                   *string                `json:"name"`
	Description            *string                `json:"description"`
	BaseURL                string                 `json:"base_url"`
	OptimizationMode       OptimizationMode       `json:"optimization_mode"`
	Enabled                bool                   `json:"enabled"`
	RefreshIntervalSeconds uint32                 `json:"refresh_interval_s"`
	DefaultOptions         *OGDefaultImageOptions `json:"default_options"`
	ExtractValues          bool                   `json:"extract_values"`
	TemplateID             *string                `json:"template_id"`
	// TemplateVersion is nil when the configuration follows the latest template version.
	TemplateVersion         *int64                   `json:"template_version"`
	TemplateValuesMapping   []OGTemplateValueMapping `json:"template_values_mapping"`
	Headers                 map[string]string        `json:"headers"`
	AdditionalHeaderOrigins []string                 `json:"additional_header_origins"`
	CreatedAt               time.Time                `json:"created_at"`
	UpdatedAt               time.Time                `json:"updated_at"`
}

func (v *OGConfig) validateResponse() string {
	if v.ID == "" || v.DomainID == "" || v.BaseURL == "" {
		return "missing OG config metadata"
	}
	switch v.ConfigType {
	case OGConfigHTMLCSS:
	case OGConfigTemplated:
		if v.TemplateID == nil || *v.TemplateID == "" {
			return "missing OG template ID"
		}
	default:
		return "unknown OG config type"
	}
	return ""
}
