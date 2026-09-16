package management

import (
	"encoding/json"
	"strconv"
)

// UnmarshalJSON accepts numbers and numeric strings as documented by the API.
func (v *Proxy) UnmarshalJSON(data []byte) error {
	type plain Proxy
	var decoded plain
	wire := struct {
		*plain
		Port *json.Number `json:"port"`
	}{plain: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.Port != nil {
		n, err := strconv.ParseUint(wire.Port.String(), 10, 16)
		if err != nil {
			return err
		}
		value := uint16(n)
		decoded.Port = &value
	}
	*v = Proxy(decoded)
	return nil
}
