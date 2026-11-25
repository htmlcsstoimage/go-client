package hcti

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"reflect"
	"testing"
)

const recordingCreated = `{"id":"recording-id","status_url":"https://hcti.io/v1/recordings/recording-id/status","recording_url":"https://hcti.io/v1/recording/recording-id.mp4"}`

func recordingBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return body
}

func TestCreateRecordingHTML(t *testing.T) {
	request := &HTMLRecordingRequest{
		HTML: "<h1>Hello</h1>", GoogleFonts: GoogleFonts{" Roboto ", "Open Sans", "Open Sans"},
		RecordingOptions: RecordingOptions{DurationMS: 1000, Audio: RecordingAssetAudio{AssetID: "asset_1"}},
	}
	c := testClient(t, 200, recordingCreated, func(r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/recording" {
			t.Fatal(r.Method, r.URL)
		}
		body := recordingBody(t, r)
		want := map[string]any{"html": "<h1>Hello</h1>", "duration_ms": float64(1000), "google_fonts": "Roboto|Open+Sans", "audio": map[string]any{"asset_id": "asset_1"}}
		if !reflect.DeepEqual(body, want) {
			t.Fatalf("payload: %#v", body)
		}
	})
	result, err := c.CreateRecording(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "recording-id" || result.StatusURL != "https://hcti.io/v1/recordings/recording-id/status" || result.RecordingURL != "https://hcti.io/v1/recording/recording-id.mp4" {
		t.Fatal(result)
	}
	if request.GoogleFonts[0] != " Roboto " {
		t.Fatal("mutated input")
	}
}

func TestRecordingOptionsAndTimeline(t *testing.T) {
	request := &URLRecordingRequest{
		URL: "https://example.com", Headers: map[string]string{"Authorization": "Bearer test"},
		AdditionalHeaderOrigins:     []string{"https://assets.example.com"},
		IncludeHeadersOnSubrequests: Ptr(false), IdentifyAsHCTI: Ptr(true), BlockConsentBanners: Ptr(true),
		RecordingOptions: RecordingOptions{
			DurationMS: 5000, CSS: Ptr("body { color: red; }"), DeviceScale: Ptr(float64(2)),
			ViewportWidth: Ptr(1280), ViewportHeight: Ptr(720), DisableTwemoji: Ptr(false),
			DedupeDurationSeconds: Ptr(0), ColorScheme: Ptr(Dark), Timezone: Ptr("America/New_York"),
			ViewportMobile: Ptr(false), ViewportLandscape: Ptr(true), ViewportTouch: Ptr(false),
			MediaType: Ptr(Screen), ProxyID: Ptr("proxy"), StorageDestinationID: Ptr("store"),
			RequestOverrides: []RequestOverride{{Action: RequestOverrideBlock, ResourceTypes: []RequestOverrideResourceType{ResourceScript}}},
			Trigger:          RecordingElementVisibleTrigger{Selector: "#ready"},
			Output:           &RecordingOutputOptions{Format: Ptr(RecordingFormatWebM), FPS: Ptr(60), Scale: Ptr(0.5)},
			Audio: RecordingURLAudio{
				URL: "https://cdn.example.com/audio/intro.mp3",
				RecordingAudioOptions: RecordingAudioOptions{
					SourceStartMS: Ptr(0), DelayMS: Ptr(500), WhenAudioEnds: Ptr(RecordingAudioEndTrimVideo), ErrorOnSourceInvalid: Ptr(false),
				},
			},
			DurationMode: Ptr(RecordingDurationModeAtLeast), DeferRender: Ptr(false), SkipScrollPreload: Ptr(true),
			Actions: []RecordingAction{
				RecordingLocateSelectorAction{Selector: "#modal", IfNotFound: []RecordingAction{}, IfFound: []RecordingAction{
					RecordingClickAction{Selector: "#close", Button: Ptr(RecordingMouseButtonRight), Count: Ptr(2), OffsetX: Ptr(float64(0)), IfNotFound: []RecordingAction{RecordingErrorAction{}}},
				}},
				RecordingExecuteJavaScriptAction{Script: "await Promise.resolve();"},
				RecordingExecuteJavaScriptAction{Script: "throw Error();", IfError: []RecordingAction{}},
				RecordingScrollAction{Steps: []ScrollStep{
					{To: PositionScrollTarget{Position: Ptr(ScrollPositionBottom)}, DurationMS: 800, Easing: Ptr(ScrollEasingEaseInOutCubic)},
					{To: CoordinatesScrollTarget{X: Ptr(float64(0)), Y: Ptr(float64(200))}, DurationMS: 600},
					{To: SelectorScrollTarget{Selector: "#footer", OffsetY: Ptr(float64(-40)), IfNotFound: []RecordingAction{}}, DurationMS: 500},
				}, IfNotFound: []RecordingAction{RecordingStallAction{DurationMS: 250}}},
				RecordingSimpleScrollAction{ScrollDistance: &ScrollDistance{Value: 90, Unit: ScrollViewportHeight}, MSBetweenScrolls: Ptr(0)},
				RecordingSimpleScrollAction{ScrollDistance: &ScrollDistance{Value: 900}},
				RecordingSimpleScrollAction{ScrollDistance: &ScrollDistance{Value: 400, Unit: ScrollPixels}},
				RecordingSimpleScrollAction{},
			},
		},
	}
	c := testClient(t, 200, recordingCreated, func(r *http.Request) {
		body := recordingBody(t, r)
		for key, want := range map[string]any{
			"viewport_width": float64(1280), "viewport_height": float64(720), "device_scale": float64(2),
			"dedupe_duration_s": float64(0), "disable_twemoji": false, "color_scheme": "dark",
			"timezone": "America/New_York", "viewport_mobile": false, "viewport_landscape": true,
			"viewport_touch": false, "media_type": "screen", "proxy_id": "proxy", "storage_destination_id": "store",
			"include_headers_on_subrequests": false, "identify_as_hcti": true, "block_consent_banners": true,
			"duration_mode": "at_least", "defer_render": false, "skip_scroll_preload": true,
		} {
			if body[key] != want {
				t.Errorf("%s: %#v", key, body[key])
			}
		}
		if !reflect.DeepEqual(body["output"], map[string]any{"format": "webm", "fps": float64(60), "scale": 0.5}) {
			t.Fatal(body["output"])
		}
		if !reflect.DeepEqual(body["audio"], map[string]any{
			"url": "https://cdn.example.com/audio/intro.mp3", "source_start_ms": float64(0), "delay_ms": float64(500),
			"when_audio_ends": "trim_video", "error_on_source_invalid": false,
		}) {
			t.Fatal(body["audio"])
		}
		if body["trigger"].(map[string]any)["type"] != "element_visible" {
			t.Fatal(body["trigger"])
		}
		actions := body["actions"].([]any)
		locate := actions[0].(map[string]any)
		if !reflect.DeepEqual(locate["if_not_found"], []any{}) {
			t.Fatal(locate)
		}
		click := locate["if_found"].([]any)[0].(map[string]any)
		if click["type"] != "click" || click["if_not_found"].([]any)[0].(map[string]any)["type"] != "error" {
			t.Fatal(click)
		}
		if _, ok := actions[1].(map[string]any)["if_error"]; ok {
			t.Fatal("omitted if_error was sent")
		}
		if !reflect.DeepEqual(actions[2].(map[string]any)["if_error"], []any{}) {
			t.Fatal("empty if_error was lost")
		}
		scroll := actions[3].(map[string]any)
		steps := scroll["steps"].([]any)
		for i, want := range []struct {
			target   string
			duration float64
		}{{"position", 800}, {"coordinates", 600}, {"selector", 500}} {
			step := steps[i].(map[string]any)
			if step["to"].(map[string]any)["type"] != want.target || step["duration_ms"] != want.duration {
				t.Fatal(steps)
			}
		}
		if !reflect.DeepEqual(scroll["if_not_found"], []any{map[string]any{"type": "stall", "duration_ms": float64(250)}}) {
			t.Fatal(scroll["if_not_found"])
		}
		if actions[4].(map[string]any)["scroll_distance"] != "90vh" || actions[5].(map[string]any)["scroll_distance"] != float64(900) || actions[6].(map[string]any)["scroll_distance"] != "400px" {
			t.Fatal(actions)
		}
		if _, ok := actions[7].(map[string]any)["scroll_distance"]; ok {
			t.Fatal("omitted scroll_distance was sent")
		}
	})
	if _, err := c.CreateRecording(context.Background(), request); err != nil {
		t.Fatal(err)
	}
}

func TestRecordingDefaultsAndTriggers(t *testing.T) {
	for _, test := range []struct {
		trigger RecordingTrigger
		want    map[string]any
	}{
		{nil, nil},
		{RecordingImmediateTrigger{}, map[string]any{"type": "immediate"}},
		{RecordingFunctionCallTrigger{}, map[string]any{"type": "function_call"}},
		{RecordingFunctionCallTrigger{TimeoutMS: Ptr(1000)}, map[string]any{"type": "function_call", "timeout_ms": float64(1000)}},
		{RecordingElementVisibleTrigger{Selector: "#ready"}, map[string]any{"type": "element_visible", "selector": "#ready"}},
		{RecordingElementVisibleTrigger{Selector: "#ready", TimeoutMS: Ptr(2500)}, map[string]any{"type": "element_visible", "selector": "#ready", "timeout_ms": float64(2500)}},
	} {
		c := testClient(t, 200, recordingCreated, func(r *http.Request) {
			body := recordingBody(t, r)
			if body["duration_ms"] != float64(0) {
				t.Fatal(body)
			}
			if _, ok := body["output"]; ok {
				t.Fatal("added output defaults")
			}
			trigger, ok := body["trigger"]
			if test.want == nil {
				if ok {
					t.Fatal(body)
				}
				return
			}
			if !reflect.DeepEqual(trigger, test.want) {
				t.Fatalf("trigger: %#v", trigger)
			}
		})
		if _, err := c.CreateRecording(context.Background(), &URLRecordingRequest{URL: "", RecordingOptions: RecordingOptions{DurationMS: 0, Trigger: test.trigger}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestScrollDistanceValidation(t *testing.T) {
	for _, distance := range []ScrollDistance{{}, {Value: -1}, {Value: math.NaN()}, {Value: math.Inf(1), Unit: ScrollPixels}, {Value: 50, Unit: "%"}} {
		if _, err := json.Marshal(distance); err == nil {
			t.Fatalf("accepted %#v", distance)
		}
	}
	data, err := json.Marshal(ScrollDistance{Value: 12.5, Unit: ScrollViewportHeight})
	if err != nil || string(data) != `"12.5vh"` {
		t.Fatal(string(data), err)
	}
}

func TestFetchRecordingStatus(t *testing.T) {
	for _, status := range []RecordingStatus{RecordingStatusAwaitingRender, RecordingStatusProcessing, RecordingStatusPendingRetry, RecordingStatusReady, RecordingStatusFailed, RecordingStatusCancelled} {
		t.Run(string(status), func(t *testing.T) {
			httpStatus := 202
			if status == RecordingStatusReady || status == RecordingStatusFailed || status == RecordingStatusCancelled {
				httpStatus = 200
			}
			body := `{"id":"recording-id","recording_url":"https://hcti.io/v1/recording/recording-id.mp4","status":"` + string(status) + `","status_entered_at":"2026-11-16T10:00:00Z","attempt_count":2,"stored_duration_ms":null,"storage_destination_hcti_storage_disabled":true,"error":{"code":"action_error","message":"Missing element","occurred_at":"2026-11-16T09:59:00Z","retryable":true,"action_path":"actions[1]"}}`
			c := testClient(t, httpStatus, body, func(r *http.Request) {
				if r.Method != "GET" || r.URL.Path != "/v1/recordings/recording-id/status" {
					t.Fatal(r.Method, r.URL)
				}
			})
			result, err := c.FetchRecordingStatus(context.Background(), "recording-id")
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != status || result.AttemptCount != 2 || result.StoredDurationMS != nil || !result.Error.Retryable || *result.Error.ActionPath != "actions[1]" || result.StatusEnteredAt.IsZero() {
				t.Fatal(result)
			}
		})
	}
}

func TestDeleteRecordings(t *testing.T) {
	c := testClient(t, 202, "", func(r *http.Request) {
		if r.Method != "DELETE" || r.URL.Path != "/v1/recordings/recording-id" {
			t.Fatal(r.Method, r.URL)
		}
	})
	if err := c.DeleteRecording(context.Background(), "recording-id"); err != nil {
		t.Fatal(err)
	}
	c = testClient(t, 202, "", func(r *http.Request) {
		if r.Method != "DELETE" || r.URL.Path != "/v1/recordings/batch" {
			t.Fatal(r.Method, r.URL)
		}
		if body := recordingBody(t, r); !reflect.DeepEqual(body["ids"], []any{"one", "two"}) {
			t.Fatal(body)
		}
	})
	if err := c.DeleteRecordingBatch(context.Background(), []string{"one", "two"}); err != nil {
		t.Fatal(err)
	}
}

func TestEmptyBatchDeletesSkipRequest(t *testing.T) {
	c := testClient(t, 500, "", func(r *http.Request) { t.Fatal("unexpected request", r.Method, r.URL) })
	for _, ids := range [][]string{nil, {}} {
		if err := c.DeleteRecordingBatch(context.Background(), ids); err != nil {
			t.Fatal(err)
		}
		if err := c.DeleteImageBatch(context.Background(), ids); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRecordingStatusUnknown(t *testing.T) {
	for raw, want := range map[string]RecordingStatus{
		`"ready"`:              RecordingStatusReady,
		`"awaiting_render"`:    RecordingStatusAwaitingRender,
		`"unknown"`:            RecordingStatusUnknown,
		`"some_future_status"`: RecordingStatusUnknown,
		`""`:                   "",
		`null`:                 "",
	} {
		var status RecordingStatus
		if err := json.Unmarshal([]byte(raw), &status); err != nil || status != want {
			t.Fatalf("%s: %q, %v", raw, status, err)
		}
	}
	var status RecordingStatus
	if err := json.Unmarshal([]byte(`1`), &status); err == nil {
		t.Fatal("accepted non-string status")
	}
	c := testClient(t, 202, `{"id":"recording-id","recording_url":"https://hcti.io/v1/recording/recording-id.mp4","status":"queued_elsewhere","status_entered_at":"2026-11-16T10:00:00Z","attempt_count":1}`, nil)
	result, err := c.FetchRecordingStatus(context.Background(), "recording-id")
	if err != nil || result.Status != RecordingStatusUnknown {
		t.Fatal(result, err)
	}
}

func TestRecordingErrorsAndCancellation(t *testing.T) {
	operations := map[string]func(*Client, context.Context) error{
		"create": func(c *Client, ctx context.Context) error {
			_, err := c.CreateRecording(ctx, &URLRecordingRequest{URL: "https://example.com", RecordingOptions: RecordingOptions{DurationMS: 1000}})
			return err
		},
		"status": func(c *Client, ctx context.Context) error {
			_, err := c.FetchRecordingStatus(ctx, "recording-id")
			return err
		},
		"delete": func(c *Client, ctx context.Context) error { return c.DeleteRecording(ctx, "recording-id") },
		"batch": func(c *Client, ctx context.Context) error {
			return c.DeleteRecordingBatch(ctx, []string{"recording-id"})
		},
	}
	for name, call := range operations {
		t.Run(name, func(t *testing.T) {
			for _, body := range []string{`{"error":"Forbidden","message":"Missing permission"}`, "Unavailable"} {
				c := testClient(t, 403, body, nil)
				var apiErr *APIError
				if err := call(c, context.Background()); !errors.As(err, &apiErr) || apiErr.StatusCode != 403 {
					t.Fatal(err)
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			c := NewClient("id", "key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) { return nil, r.Context().Err() })}))
			if err := call(c, ctx); !errors.Is(err, context.Canceled) {
				t.Fatal(err)
			}
			transportErr := errors.New("network failure")
			c = NewClient("id", "key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, transportErr })}))
			if err := call(c, context.Background()); !errors.Is(err, transportErr) {
				t.Fatal(err)
			}
		})
	}
	c := testClient(t, 200, "{}", nil)
	var malformed *ResponseError
	if _, err := c.CreateRecording(context.Background(), &URLRecordingRequest{}); !errors.As(err, &malformed) {
		t.Fatal(err)
	}
}
