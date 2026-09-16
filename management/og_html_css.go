package management

import "encoding/json"

// HTMLCSSOGConfigRequest configures page HTML/CSS rendering.
type HTMLCSSOGConfigRequest struct {
	OGConfigOptions
	// DefaultOptions supplies rendering options. Nil clears configured defaults.
	DefaultOptions *OGDefaultImageOptions `json:"default_options"`
	// ExtractValues allows page metadata to override DefaultOptions.
	ExtractValues bool `json:"extract_values"`
}

func (*HTMLCSSOGConfigRequest) ogConfigRequest() {}

// MarshalJSON includes the html_css discriminator.
func (v HTMLCSSOGConfigRequest) MarshalJSON() ([]byte, error) {
	type fields HTMLCSSOGConfigRequest
	return json.Marshal(struct {
		ConfigType OGConfigType `json:"config_type"`
		fields
	}{OGConfigHTMLCSS, fields(v)})
}
