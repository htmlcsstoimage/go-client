package management

import "time"

// APIKeyRequest creates a key or replaces its editable settings. Updates never rotate the secret.
type APIKeyRequest struct {
	// Name is the display name. Nil or blank asks the API to generate a name.
	Name *string `json:"name"`
	// Description is optional; nil clears it on update.
	Description *string `json:"description"`
	// Disabled prevents this key from authenticating.
	Disabled bool `json:"disabled"`
	// AllFuturePermissions grants every current and future permission.
	AllFuturePermissions bool `json:"all_future_permissions"`
	// Permissions is required; use an empty slice for no permissions or when AllFuturePermissions is true.
	Permissions []Permission `json:"permissions"`
}

// APIKey is key metadata; ID is the management identifier, distinct from APIID.
type APIKey struct {
	ID                   string       `json:"id"`
	APIID                string       `json:"api_id"`
	Name                 string       `json:"name"`
	Description          *string      `json:"description"`
	Enabled              bool         `json:"enabled"`
	AllFuturePermissions bool         `json:"all_future_permissions"`
	Permissions          []Permission `json:"permissions"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}

// APIKeyWithSecret is returned only at creation. Store Secret securely; subsequent reads cannot recover it.
type APIKeyWithSecret struct {
	APIKey
	Secret string `json:"api_key"`
}

func (v *APIKey) validateResponse() string {
	if v.ID == "" || v.APIID == "" || v.Name == "" || v.Permissions == nil {
		return "missing API key metadata"
	}
	return ""
}
func (v *APIKeyWithSecret) validateResponse() string {
	if v.Secret == "" {
		return "missing API key secret"
	}
	return v.APIKey.validateResponse()
}

// APIKeyListOptions filters a single page of keys.
type APIKeyListOptions struct {
	ListOptions
	// IncludeDisabled includes disabled keys; false excludes them.
	IncludeDisabled bool
	// WithPermission requires every listed permission.
	WithPermission []Permission
}
