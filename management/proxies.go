package management

import "time"

// ProxyRequest creates a proxy or replaces all its editable settings.
type ProxyRequest struct {
	// Name is a display name of 3–500 characters.
	Name string `json:"name"`
	// URL is an HTTP(S) origin without a port, path, or embedded credentials.
	URL string `json:"url"`
	// Port is 1–65535. Nil uses the scheme's default port.
	Port *uint16 `json:"port"`
	// Disabled prevents the proxy from being used for rendering.
	Disabled bool `json:"disabled"`
	// Authentication configures credentials. Nil removes authentication on update.
	Authentication *ProxyAuthentication `json:"authentication"`
	// BypassHosts lists up to 100 hosts that should connect directly.
	BypassHosts []string `json:"bypass_hosts"`
}

// ProxyAuthentication supplies replacement credentials or explicitly retains the existing password.
type ProxyAuthentication struct {
	// Username is required; an empty username is valid and whitespace is preserved.
	Username string `json:"username"`
	// Password replaces the stored password. An explicit empty string is valid.
	// Nil requires RetainPassword=true; passwords are never returned on reads.
	Password *string `json:"password"`
	// RetainPassword=true is update-only and requires unchanged Username and nil Password.
	// Nil or false requires a supplied Password, even when the proxy is disabled.
	RetainPassword *bool `json:"retain_password"`
}

// Proxy describes a proxy without disclosing its password.
type Proxy struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Port        *uint16  `json:"port"`
	BypassHosts []string `json:"bypass_hosts"`
	// Username is nil without authentication; a pointer to an empty string is a valid username.
	Username  *string   `json:"username"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (v *Proxy) validateResponse() string {
	if v.ID == "" || v.Name == "" || v.URL == "" || v.BypassHosts == nil {
		return "missing proxy metadata"
	}
	return ""
}
