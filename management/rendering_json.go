package management

import (
	"encoding/json"
	"strconv"
)

// UnmarshalJSON accepts the API's number-or-string version representation.
func (v *CreatedRender) UnmarshalJSON(data []byte) error {
	type plain CreatedRender
	var p plain
	w := struct {
		*plain
		Version *json.Number `json:"template_version"`
	}{plain: &p}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	if w.Version != nil {
		n, err := strconv.ParseInt(w.Version.String(), 10, 64)
		if err != nil {
			return err
		}
		p.TemplateVersion = n
	}
	*v = CreatedRender(p)
	return nil
}

// UnmarshalJSON accepts number-or-string PDF scale values.
func (v *RenderPDFOptions) UnmarshalJSON(data []byte) error {
	type plain RenderPDFOptions
	var p plain
	w := struct {
		*plain
		Scale *json.Number `json:"scale"`
	}{plain: &p}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	if w.Scale != nil {
		n, err := strconv.ParseFloat(w.Scale.String(), 64)
		if err != nil {
			return err
		}
		p.Scale = &n
	}
	*v = RenderPDFOptions(p)
	return nil
}
