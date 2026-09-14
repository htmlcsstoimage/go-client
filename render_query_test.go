package hcti

import (
	"encoding/json"
	"math"
	"net/url"
	"reflect"
	"sort"
	"testing"
)

func TestRenderQueryMatchesWireFields(t *testing.T) {
	request := URLImageRequest{
		URL:                         "https://example.com/?q=hello world&x=✓",
		CSS:                         Ptr("body::before { content: \"<&>\"; }"),
		Headers:                     map[string]string{"Z-Test": "last", "A-Test": "first"},
		AdditionalHeaderOrigins:     []string{"https://cdn.example.com", "https://assets.example.com"},
		IncludeHeadersOnSubrequests: Ptr(false), IdentifyAsHCTI: Ptr(true),
		FullScreen: Ptr(false), BlockConsentBanners: Ptr(true),
		ImageOptions: ImageOptions{
			Format: PDF, Selector: Ptr("#main"), MaxRenderOnce: Ptr(false),
			DedupeDurationSeconds: Ptr(0), PDFOptions: &PDFOptions{Scale: Ptr(1.0)},
			RenderOptions: RenderOptions{
				DeviceScale: Ptr(1.5), ViewportHeight: Ptr(720), ViewportWidth: Ptr(1280),
				MaxWaitMS: Ptr(500), MSDelay: Ptr(0), RenderWhenReady: Ptr(false),
				DisableTwemoji: Ptr(true), ColorScheme: Ptr(Dark), Timezone: Ptr("America/New_York"),
				ViewportMobile: Ptr(false), ViewportTouch: Ptr(true), ViewportLandscape: Ptr(false),
				MediaType: Ptr(Print), ProxyID: Ptr("proxy-id"), StorageDestinationID: Ptr("storage-id"),
				JumboMaxWidth: Ptr(9000), JumboMaxHeight: Ptr(1000), TransparentBackground: Ptr(false),
			},
		},
	}
	for _, request := range []URLImageRequest{request, {URL: "https://example.com"}} {
		got, err := renderQueryForTest(request)
		if err != nil {
			t.Fatal(err)
		}
		// The HTTP wire representation is an independent compatibility reference.
		data, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		want := url.Values{}
		for key, raw := range fields {
			switch key {
			case "format", "pdf_options", "dedupe_duration_s":
				continue
			case "headers":
				keys := make([]string, 0, len(request.Headers))
				for name := range request.Headers {
					keys = append(keys, name)
				}
				sort.Strings(keys)
				for _, name := range keys {
					want.Add(key, name+":"+request.Headers[name])
				}
			case "additional_header_origins":
				want[key] = request.AdditionalHeaderOrigins
			default:
				var value string
				if raw[0] == '"' {
					if err := json.Unmarshal(raw, &value); err != nil {
						t.Fatal(err)
					}
				} else {
					value = string(raw)
				}
				want.Set(key, value)
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}
		b, err := newSignedURLBuilder(nil, "https://hcti.io", "/test", request.Format)
		if err != nil {
			t.Fatal(err)
		}
		if err := request.appendRenderQuery(&b); err != nil {
			t.Fatal(err)
		}
		if encoded := string(b.buf[b.queryStart:]); encoded != want.Encode() {
			t.Fatalf("query order/encoding changed: got %s; want %s", encoded, want.Encode())
		}
	}
}

func TestRenderQueryDeviceScale(t *testing.T) {
	for _, value := range []float64{0, math.Copysign(0, -1), 0.1, 1.5, 3, 1e-9, 1e21} {
		request := URLImageRequest{URL: "https://example.com", ImageOptions: ImageOptions{RenderOptions: RenderOptions{DeviceScale: &value}}}
		query, err := renderQueryForTest(request)
		if err != nil {
			t.Fatal(err)
		}
		want, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if query.Get("device_scale") != string(want) {
			t.Fatalf("got %s; want %s", query.Get("device_scale"), want)
		}
	}
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		request := URLImageRequest{URL: "https://example.com", ImageOptions: ImageOptions{RenderOptions: RenderOptions{DeviceScale: &value}}}
		if _, err := renderQueryForTest(request); err == nil {
			t.Fatal("accepted nonfinite device scale")
		}
	}
}

// Keep map decoding in tests only, as an independent oracle for buffer output.
func renderQueryForTest(request URLImageRequest) (url.Values, error) {
	b, err := newSignedURLBuilder(nil, "https://hcti.io", "/test", request.Format)
	if err != nil {
		return nil, err
	}
	if err := request.appendRenderQuery(&b); err != nil {
		return nil, err
	}
	return url.ParseQuery(string(b.buf[b.queryStart:]))
}
