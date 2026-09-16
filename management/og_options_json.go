package management

import (
	"encoding/json"
	"strconv"
)

// UnmarshalJSON accepts numbers and numeric strings as documented by the API.
func (v *OGDefaultImageOptions) UnmarshalJSON(data []byte) error {
	type plain OGDefaultImageOptions
	var decoded plain
	wire := struct {
		*plain
		DeviceScale    *json.Number `json:"device_scale"`
		MaxWaitMS      *json.Number `json:"max_wait_ms"`
		MSDelay        *json.Number `json:"ms_delay"`
		ViewportHeight *json.Number `json:"viewport_height"`
		ViewportWidth  *json.Number `json:"viewport_width"`
	}{plain: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.DeviceScale != nil {
		n, err := strconv.ParseFloat(wire.DeviceScale.String(), 64)
		if err != nil {
			return err
		}
		value := float64(n)
		decoded.DeviceScale = &value
	}
	if wire.MaxWaitMS != nil {
		n, err := strconv.ParseInt(wire.MaxWaitMS.String(), 10, 32)
		if err != nil {
			return err
		}
		value := int32(n)
		decoded.MaxWaitMS = &value
	}
	if wire.MSDelay != nil {
		n, err := strconv.ParseInt(wire.MSDelay.String(), 10, 32)
		if err != nil {
			return err
		}
		value := int32(n)
		decoded.MSDelay = &value
	}
	if wire.ViewportHeight != nil {
		n, err := strconv.ParseInt(wire.ViewportHeight.String(), 10, 32)
		if err != nil {
			return err
		}
		value := int32(n)
		decoded.ViewportHeight = &value
	}
	if wire.ViewportWidth != nil {
		n, err := strconv.ParseInt(wire.ViewportWidth.String(), 10, 32)
		if err != nil {
			return err
		}
		value := int32(n)
		decoded.ViewportWidth = &value
	}
	*v = OGDefaultImageOptions(decoded)
	return nil
}
