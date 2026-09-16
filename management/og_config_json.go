package management

import (
	"encoding/json"
	"strconv"
)

// UnmarshalJSON accepts numbers and numeric strings as documented by the API.
func (v *OGConfig) UnmarshalJSON(data []byte) error {
	type plain OGConfig
	var decoded plain
	wire := struct {
		*plain
		TemplateVersion        *json.Number `json:"template_version"`
		RefreshIntervalSeconds *json.Number `json:"refresh_interval_s"`
	}{plain: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.TemplateVersion != nil {
		n, err := strconv.ParseInt(wire.TemplateVersion.String(), 10, 64)
		if err != nil {
			return err
		}
		value := int64(n)
		decoded.TemplateVersion = &value
	}
	if wire.RefreshIntervalSeconds != nil {
		n, err := strconv.ParseUint(wire.RefreshIntervalSeconds.String(), 10, 32)
		if err != nil {
			return err
		}
		value := uint32(n)
		decoded.RefreshIntervalSeconds = value
	}
	*v = OGConfig(decoded)
	return nil
}
