package hcti

import (
	"net/url"
	"strings"
	"testing"
)

func TestImageURLTransformations(t *testing.T) {
	c := NewClient("id", "secret", WithBaseURL("https://example.com/"))
	px := func(v int) *CropValue { return &CropValue{v, CropPixels} }
	pct := func(v int) *CropValue { return &CropValue{v, CropPercent} }
	cases := []struct {
		name    string
		options RenderImageOptions
		query   string
	}{
		{"none", RenderImageOptions{}, ""},
		{"resize", RenderImageOptions{DPI: Ptr(144), Width: Ptr(600), Height: Ptr(400)}, "dpi=144&height=400&width=600"},
		{"rectangle", RenderImageOptions{Crop: &Crop{Horizontal: &CropSpan{Start: px(0), End: px(600)}, Vertical: &CropSpan{Start: pct(10), Size: pct(80)}}}, "x_1=0px&x_2=600px&y_1=10%25&crop_height=80%25"},
		{"mixed-units", RenderImageOptions{Crop: &Crop{Horizontal: &CropSpan{Start: px(200), End: pct(90)}}}, "x_1=200px&x_2=90%25"},
		{"remaining", RenderImageOptions{Crop: &Crop{Vertical: &CropSpan{Start: px(20)}}}, "y_1=20px"},
		{"center-end", RenderImageOptions{Crop: &Crop{Horizontal: &CropSpan{Size: pct(50), Origin: CropCenter}, Vertical: &CropSpan{Size: px(200), Origin: CropEnd}}}, "x_origin=center&y_origin=end&crop_width=50%25&crop_height=200px"},
		{"aspect-width", RenderImageOptions{Width: Ptr(500), Crop: &Crop{Horizontal: &CropSpan{Size: pct(100)}, AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropWidth, ComputedOrigin: CropCenter}}, "width=500&aspect_ratio=16_9&y_origin=center&crop_width=100%25"},
		{"aspect-height", RenderImageOptions{Crop: &Crop{Vertical: &CropSpan{Start: px(0), Size: px(900)}, AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropHeight, ComputedOrigin: CropEnd}}, "aspect_ratio=16_9&x_origin=end&y_1=0px&crop_height=900px"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, format := range []ImageFormat{"", PNG, JPG, WebP, PDF} {
				options := tc.options
				options.Format = format
				got, err := c.ImageURL("image/id ?#👩🏽‍💻", options)
				if err != nil {
					t.Fatal(err)
				}
				want := "https://example.com/v1/image/" + url.PathEscape("image/id ?#👩🏽‍💻")
				if format != "" {
					want += "." + string(format)
				}
				if tc.query != "" {
					want += "?" + tc.query
				}
				if got != want {
					t.Fatalf("got %s; want %s", got, want)
				}
				template, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", Format: JPG}, options)
				if err != nil {
					t.Fatal(err)
				}
				effective := format
				if effective == "" {
					effective = JPG
				}
				assertSignedQuery(t, template, "secret", effective, tc.query)
				screenshot, err := c.GenerateCreateAndRenderURL(URLImageRequest{URL: "https://example.com", ImageOptions: ImageOptions{Format: JPG}}, options)
				if err != nil {
					t.Fatal(err)
				}
				query := "url=https%3A%2F%2Fexample.com"
				if tc.query != "" {
					query += "&" + tc.query
				}
				assertSignedQuery(t, screenshot, "secret", effective, query)
			}
		})
	}
}

func TestRenderOptionTemplateCollisions(t *testing.T) {
	c := NewClient("id", "secret")
	values := map[string]any{"width": "variable", "dpi": nil, "crop_w": "short", "crop_width": "long", "x_1": "left", "y_origin": "origin", "aspect_ratio": "ratio"}
	options := RenderImageOptions{Width: Ptr(300), DPI: Ptr(72), Crop: &Crop{
		Horizontal:  &CropSpan{Start: &CropValue{0, CropPixels}, Size: &CropValue{50, CropPercent}},
		AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropWidth, ComputedOrigin: CropCenter,
	}}
	got, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: values}, options)
	if err != nil {
		t.Fatal(err)
	}
	want := "aspect_ratio=ratio&crop_w=short&crop_width=long&width=variable&x_1=left&y_origin=origin&__ro_dpi=72&__ro_width=300&__ro_aspect_ratio=16_9&__ro_y_origin=center&__ro_x_1=0px&__ro_crop_width=50%25&__ro_crop_w=50%25"
	assertSignedQuery(t, got, "secret", "", want)
	if values["width"] != "variable" || *options.Width != 300 || options.Format != "" {
		t.Fatal("input mutated")
	}
	// Canonical and legacy size names each reserve their own override independently.
	for _, pair := range []struct{ key, expected string }{{"crop_h", "__ro_crop_h"}, {"crop_height", "__ro_crop_height"}} {
		got, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: map[string]any{pair.key: nil}}, RenderImageOptions{Crop: &Crop{Vertical: &CropSpan{Size: &CropValue{100, CropPixels}}}})
		if err != nil {
			t.Fatal(err)
		}
		assertSignedQuery(t, got, "secret", "", pair.expected+"=100px")
	}
	values["__ro_width"] = "ambiguous"
	if _, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl", TemplateValues: values}, options); err == nil {
		t.Fatal("accepted ambiguous override collision")
	}
}

func TestInvalidRenderOptions(t *testing.T) {
	px := func(v int) *CropValue { return &CropValue{v, CropPixels} }
	invalidCrops := []*Crop{
		{},
		{Horizontal: &CropSpan{}},
		{Horizontal: &CropSpan{Start: px(-1)}},
		{Horizontal: &CropSpan{Size: px(0)}},
		{Horizontal: &CropSpan{Size: &CropValue{1, "em"}}},
		{Horizontal: &CropSpan{Start: &CropValue{0, CropPercent}}},
		{Horizontal: &CropSpan{Size: &CropValue{101, CropPercent}}},
		{Horizontal: &CropSpan{End: px(10)}},
		{Horizontal: &CropSpan{Start: px(10), End: px(10)}},
		{Horizontal: &CropSpan{Start: px(10), End: px(5)}},
		{Horizontal: &CropSpan{Start: px(0), End: px(20), Size: px(10)}},
		{Horizontal: &CropSpan{Start: px(0), Origin: CropCenter}},
		{Horizontal: &CropSpan{Size: px(10), Origin: "invalid"}},
		{Horizontal: &CropSpan{Size: px(10)}, ComputedOrigin: CropCenter},
		{Horizontal: &CropSpan{Size: px(10)}, AspectRatioAxis: CropWidth},
		{Horizontal: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{0, 9}, AspectRatioAxis: CropWidth},
		{Horizontal: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{16, -1}, AspectRatioAxis: CropWidth},
		{Horizontal: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{16, 9}},
		{Horizontal: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropHeight},
		{Vertical: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropWidth},
		{Vertical: &CropSpan{Size: px(10)}, Horizontal: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropHeight},
		{Vertical: &CropSpan{Size: px(10)}, AspectRatio: &AspectRatio{16, 9}, AspectRatioAxis: CropHeight, ComputedOrigin: "invalid"},
	}
	invalid := []RenderImageOptions{{Format: "gif"}, {DPI: Ptr(30)}, {DPI: Ptr(600)}, {Width: Ptr(0)}, {Width: Ptr(5001)}, {Height: Ptr(-1)}, {Height: Ptr(5001)}}
	for _, crop := range invalidCrops {
		invalid = append(invalid, RenderImageOptions{Crop: crop})
	}
	c := NewClient("id", "secret")
	for i, options := range invalid {
		if _, err := c.ImageURL("image", options); err == nil {
			t.Fatalf("ImageURL accepted invalid case %d", i)
		}
		if _, err := c.GenerateTemplatedImageURL(TemplatedImageRequest{TemplateID: "tpl"}, options); err == nil {
			t.Fatalf("template accepted invalid case %d", i)
		}
		if _, err := c.GenerateCreateAndRenderURL(URLImageRequest{URL: "https://example.com"}, options); err == nil {
			t.Fatalf("screenshot accepted invalid case %d", i)
		}
	}
	for _, options := range []RenderImageOptions{{DPI: Ptr(31), Width: Ptr(1), Height: Ptr(1)}, {DPI: Ptr(599), Width: Ptr(5000), Height: Ptr(5000)}} {
		if _, err := c.ImageURL("image", options); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := c.ImageURL("image", RenderImageOptions{}, RenderImageOptions{}); err == nil {
		t.Fatal("accepted multiple option sets")
	}
	for _, id := range []string{"", ".", ".."} {
		if _, err := c.ImageURL(id); err == nil {
			t.Fatal("accepted invalid ID")
		}
	}
}

func TestImageURLArbitraryIDBytes(t *testing.T) {
	c := NewClient("id", "secret")
	for _, tc := range encodingCases {
		id := "image-" + tc.value
		got, err := c.ImageURL(id)
		if err != nil {
			t.Fatal(err)
		}
		u, err := url.Parse(got)
		if err != nil || !strings.HasSuffix(u.Path, id) {
			t.Fatalf("ID bytes changed for %q: %s (%v)", id, got, err)
		}
	}
}
