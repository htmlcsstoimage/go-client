package hcti

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strings"
	"testing"
)

type customTemplateString struct{}

func (customTemplateString) MarshalJSON() ([]byte, error) {
	return []byte(`"custom \ud83d\ude00 <&>"`), nil
}

type namedTemplateString string

func TestTemplateScalarJSONCompatibility(t *testing.T) {
	var nilString *string
	values := []any{
		"", "👩🏽‍💻 <>& \" \\ \u2028\u2029", "bad\xff\xfeutf8", true, false,
		int(-1), int8(-128), int16(-32768), int32(-2147483648), int64(math.MinInt64),
		uint(1), uint8(255), uint16(65535), uint32(math.MaxUint32), uint64(math.MaxUint64), uintptr(42),
		float32(0), float32(1e-6), float32(1e-9), float32(1e21), float32(math.MaxFloat32),
		float64(0), math.Copysign(0, -1), float64(1e-6), float64(1e-9), float64(1e21), math.SmallestNonzeroFloat64, math.MaxFloat64,
		json.Number("1e+3"), namedTemplateString("<&>💻"), customTemplateString{}, Ptr("pointer"), nilString,
		[]byte{0, 255}, []any{1, "👩🏽‍💻"}, map[string]any{"value": "<>&\xff"},
	}
	c := NewClient("id", "secret")
	for _, value := range values {
		want, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		text := string(want)
		if len(want) != 0 && want[0] == '"' {
			if err := json.Unmarshal(want, &text); err != nil {
				t.Fatal(err)
			}
		}
		got, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: map[string]any{"value": value}})
		if err != nil {
			t.Fatalf("%T: %v", value, err)
		}
		assertSignedQuery(t, got, "secret", "", url.Values{"value": {text}}.Encode())
	}
	for _, value := range []any{math.NaN(), math.Inf(1), math.Inf(-1), float32(math.NaN()), float32(math.Inf(1)), json.Number("invalid"), make(chan int)} {
		if _, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: map[string]any{"value": value}}); err == nil {
			t.Fatalf("accepted %T: %v", value, value)
		}
	}
}

func FuzzTemplateSigning(f *testing.F) {
	for _, tc := range encodingCases {
		f.Add(tc.value, tc.value)
	}
	c := NewClient("id", "secret")
	f.Fuzz(func(t *testing.T, key, value string) {
		if key == "template_version" {
			return
		}
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var normalized string
		if err := json.Unmarshal(data, &normalized); err != nil {
			t.Fatal(err)
		}
		got, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: map[string]any{key: value}, TemplateVersion: Ptr(int64(42))})
		if err != nil {
			t.Fatal(err)
		}
		assertSignedQuery(t, got, "secret", "", url.Values{key: {normalized}, "template_version": {"42"}}.Encode())
	})
}

func FuzzTemplateFloatEncoding(f *testing.F) {
	for _, value := range []float64{0, math.Copysign(0, -1), 1e-6, 1e21, math.MaxFloat64, math.SmallestNonzeroFloat64} {
		f.Add(value)
	}
	f.Fuzz(func(t *testing.T, value float64) {
		for _, scalar := range []any{value, float32(value)} {
			b := signedURLBuilder{}
			err := b.appendTemplateValue(scalar)
			want, jsonErr := json.Marshal(scalar)
			if (err == nil) != (jsonErr == nil) {
				t.Fatalf("mismatched errors for %v: %v vs %v", scalar, err, jsonErr)
			}
			if err != nil {
				continue
			}
			if string(b.buf) != url.QueryEscape(string(want)) {
				t.Fatalf("%T %v: got %s want %s", scalar, scalar, b.buf, want)
			}
		}
	})
}

func TestTemplateSigningLargeMapMatchesPrevious(t *testing.T) {
	values := make(map[string]any)
	for i := 0; i < 40; i++ {
		values[fmt.Sprintf("field %02d 💻", i)] = strings.Repeat("<>&👩🏽‍💻", 20)
	}
	values["omitted"] = nil
	c := NewClient("id", "secret")
	for _, format := range []ImageFormat{"", PNG, JPG, WebP, PDF} {
		r := TemplatedImageRequest{TemplateID: "template/id", TemplateVersion: Ptr(int64(42)), TemplateValues: values, Format: format}
		want, err := legacyTemplateURL(c, r)
		if err != nil {
			t.Fatal(err)
		}
		got, err := c.GenerateTemplatedImageURL(r)
		if err != nil || got != want {
			t.Fatalf("large-map signature changed (%v)", err)
		}
	}
}
