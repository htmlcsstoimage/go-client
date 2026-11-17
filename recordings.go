package hcti

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Recording identifies a created or deduplicated recording. RecordingURL can be embedded or requested directly; the request waits until the recording is ready.
type Recording struct {
	// ID is the recording identifier.
	ID string `json:"id"`
	// StatusURL is the authenticated status endpoint.
	StatusURL string `json:"status_url"`
	// RecordingURL retrieves the media and waits for processing.
	RecordingURL string `json:"recording_url"`
}

// RecordingError describes a recording attempt. Use status to distinguish retries from terminal failure.
type RecordingError struct {
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	OccurredAt time.Time `json:"occurred_at"`
	Retryable  bool      `json:"retryable"`
	ActionPath *string   `json:"action_path"`
}

// RecordingStatusResponse contains lifecycle information from HTTP 200 or 202.
// Ready, failed, and cancelled are terminal; Error can describe an earlier attempt.
type RecordingStatusResponse struct {
	ID                                    string          `json:"id"`
	RecordingURL                          string          `json:"recording_url"`
	Status                                RecordingStatus `json:"status"`
	StatusEnteredAt                       time.Time       `json:"status_entered_at"`
	AttemptCount                          uint32          `json:"attempt_count"`
	StoredDurationMS                      *int64          `json:"stored_duration_ms"`
	StorageDestinationHCTIStorageDisabled *bool           `json:"storage_destination_hcti_storage_disabled"`
	Error                                 *RecordingError `json:"error"`
}

// CreateRecording creates an HTML/CSS or URL recording. Options are validated by the API.
// See https://docs.htmlcsstoimage.com/getting-started/recordings/.
func (c *Client) CreateRecording(ctx context.Context, request RecordingRequest) (*Recording, error) {
	switch r := request.(type) {
	case *HTMLRecordingRequest:
		if r == nil {
			return nil, fmt.Errorf("hcti: recording request is required")
		}
	case *URLRecordingRequest:
		if r == nil {
			return nil, fmt.Errorf("hcti: recording request is required")
		}
	default:
		return nil, fmt.Errorf("hcti: recording request is required")
	}
	var result Recording
	if err := c.do(ctx, http.MethodPost, "/v1/recording", nil, request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FetchRecordingStatus fetches lifecycle status once. HTTP 202 is a successful status fetch.
func (c *Client) FetchRecordingStatus(ctx context.Context, id string) (*RecordingStatusResponse, error) {
	var result RecordingStatusResponse
	if err := c.do(ctx, http.MethodGet, "/v1/recordings/"+id+"/status", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteRecording deletes a recording and clears its media from the CDN.
func (c *Client) DeleteRecording(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/recordings/"+id, nil, nil, nil)
}

// DeleteRecordingBatch deletes recording IDs in one API call.
// An empty list succeeds without calling the API.
func (c *Client) DeleteRecordingBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return c.do(ctx, http.MethodDelete, "/v1/recordings/batch", nil, struct {
		IDs []string `json:"ids"`
	}{ids}, nil)
}

func (r *Recording) validateResponse() string {
	if r.ID == "" || r.StatusURL == "" || r.RecordingURL == "" {
		return "missing recording ID, status URL, or media URL"
	}
	return ""
}

func (r *RecordingStatusResponse) validateResponse() string {
	if r.ID == "" || r.RecordingURL == "" || r.Status == "" {
		return "missing recording ID, media URL, or status"
	}
	return ""
}
