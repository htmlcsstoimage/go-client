package hcti

// RequestOverrideAction selects the action for a browser request rule.
type RequestOverrideAction string

const RequestOverrideBlock RequestOverrideAction = "block"

// RequestOverrideResourceType is a browser network resource type.
type RequestOverrideResourceType string

const (
	ResourceBeacon      RequestOverrideResourceType = "beacon"
	ResourceDocument    RequestOverrideResourceType = "document"
	ResourceStylesheet  RequestOverrideResourceType = "stylesheet"
	ResourceImage       RequestOverrideResourceType = "image"
	ResourceImageSet    RequestOverrideResourceType = "image_set"
	ResourceMedia       RequestOverrideResourceType = "media"
	ResourceFont        RequestOverrideResourceType = "font"
	ResourceScript      RequestOverrideResourceType = "script"
	ResourceTextTrack   RequestOverrideResourceType = "text_track"
	ResourceXHR         RequestOverrideResourceType = "xhr"
	ResourceFetch       RequestOverrideResourceType = "fetch"
	ResourceEventSource RequestOverrideResourceType = "event_source"
	ResourceManifest    RequestOverrideResourceType = "manifest"
	ResourcePing        RequestOverrideResourceType = "ping"
	ResourceImg         RequestOverrideResourceType = "img"
	ResourceOther       RequestOverrideResourceType = "other"
)

// RequestOverride blocks a request when its URL and/or resource type matches.
// When both matchers are set, both must match. The API requires at least one.
type RequestOverride struct {
	Action        RequestOverrideAction         `json:"action"`
	URL           *string                       `json:"url,omitempty"`
	ResourceTypes []RequestOverrideResourceType `json:"resource_types,omitempty"`
}
