package hcti

import (
	"encoding/json"
	"strings"
)

// GoogleFonts is a list of font families, serialized in the API's pipe-delimited format.
type GoogleFonts []string

func (fonts GoogleFonts) MarshalJSON() ([]byte, error) {
	seen := make(map[string]bool)
	values := make([]string, 0, len(fonts))
	for _, font := range fonts {
		font = strings.ReplaceAll(strings.TrimSpace(font), " ", "+")
		if font != "" && !seen[font] {
			values = append(values, font)
			seen[font] = true
		}
	}
	return json.Marshal(strings.Join(values, "|"))
}

func (fonts *GoogleFonts) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*fonts = nil
	if value != "" {
		for _, font := range strings.Split(value, "|") {
			*fonts = append(*fonts, strings.ReplaceAll(font, "+", " "))
		}
	}
	return nil
}
