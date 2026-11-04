package management

import (
	"encoding/json"
	"strconv"
	"time"
)

// WebhookRequest creates a webhook or replaces its complete configuration.
type WebhookRequest struct {
	// Name is a display name, trimmed by the API.
	Name string `json:"name"`
	// URL is an HTTPS endpoint without credentials, query, or fragment.
	URL string `json:"url"`
	// Disabled prevents deliveries. Enabling queues an automatic test event.
	Disabled bool `json:"disabled"`
	// Events selects subscriptions. An empty slice receives only test events.
	Events []string `json:"events"`
	// Parameters replaces all headers and query parameters. Nil removes them.
	Parameters []WebhookParameterRequest `json:"parameters"`
	// SigningEnabled signs deliveries using the secret returned at creation.
	SigningEnabled bool `json:"signing_enabled"`
}

// WebhookParameterRequest supplies a delivery parameter or retains a stored secret.
type WebhookParameterRequest struct {
	// Location is header or query_parameter.
	Location string `json:"location"`
	// Key identifies the parameter. Header names match case-insensitively.
	Key string `json:"key"`
	// Secret hides the value from API reads and delivery history.
	Secret bool `json:"secret"`
	// Value is required unless RetainValue is true. An empty string is valid.
	Value *string `json:"value"`
	// RetainValue keeps an existing secret on update, with nil Value and unchanged URL.
	RetainValue bool `json:"retain_value"`
}

// WebhookParameter describes a delivery parameter. Secret values are never returned.
type WebhookParameter struct {
	Location string  `json:"location"`
	Key      string  `json:"key"`
	Secret   bool    `json:"secret"`
	Value    *string `json:"value"`
}

// Webhook is a webhook configuration. Its signing secret is returned only on REST creation.
type Webhook struct {
	// SigningSecret is returned only by CreateWebhook, even when signing is disabled.
	// Store it securely; updates, reads, and lists do not return it.
	SigningSecret  *string            `json:"signing_secret"`
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	URL            string             `json:"url"`
	Enabled        bool               `json:"enabled"`
	Events         []string           `json:"events"`
	Parameters     []WebhookParameter `json:"parameters"`
	SigningEnabled bool               `json:"signing_enabled"`
	ConfigVersion  int32              `json:"config_version"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	// TestEventID is populated only when this write queued an automatic test.
	TestEventID *string `json:"test_event_id"`
}

// UnmarshalJSON accepts numeric strings as well as numbers for config_version.
func (v *Webhook) UnmarshalJSON(data []byte) error {
	type plain Webhook
	var decoded plain
	wire := struct {
		*plain
		ConfigVersion json.Number `json:"config_version"`
	}{plain: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.ConfigVersion != "" {
		n, err := strconv.ParseInt(wire.ConfigVersion.String(), 10, 32)
		if err != nil {
			return err
		}
		decoded.ConfigVersion = int32(n)
	}
	*v = Webhook(decoded)
	return nil
}

func (v *Webhook) validateResponse() string {
	if v.ID == "" || v.Name == "" || v.URL == "" || v.Events == nil || v.Parameters == nil {
		return "missing webhook metadata"
	}
	return ""
}
