package httpapi

import (
	"net/http"
	"strings"
	"time"
)

// Config holds connection settings shared by the public SDK packages.
// Configure it before constructing a transport or sharing a client.
type Config struct {
	APIID, APIKey string
	BaseURL       string
	userAgent     string
	httpClient    *http.Client
}

// NewConfig supplies the common API origin and SDK identification.
func NewConfig(apiID, apiKey, version string) Config {
	return Config{APIID: apiID, APIKey: apiKey, BaseURL: "https://hcti.io", userAgent: "HCTIGo/" + version}
}

// SetBaseURL normalizes the origin for HTTP requests and signed URLs.
func (c *Config) SetBaseURL(baseURL string) { c.BaseURL = strings.TrimRight(baseURL, "/") }

// SetHTTPClient snapshots caller settings. Nil leaves the current settings intact.
func (c *Config) SetHTTPClient(client *http.Client) {
	if client != nil {
		copied := *client
		c.httpClient = &copied
	}
}

// newHTTPClient applies shared timeout and redirect policies without mutating the caller's client.
func (c Config) newHTTPClient() *http.Client {
	client := http.Client{Timeout: 60 * time.Second}
	if c.httpClient != nil {
		client = *c.httpClient
	}
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &client
}

// AppendUserAgent adds an integration identifier while retaining the SDK identifier.
func (c *Config) AppendUserAgent(suffix string) {
	if suffix = strings.TrimSpace(suffix); suffix != "" {
		c.userAgent += " " + suffix
	}
}
