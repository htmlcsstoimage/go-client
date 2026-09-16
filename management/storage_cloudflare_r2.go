package management

import "encoding/json"

// CloudflareR2Connection configures CloudflareR2 storage.
type CloudflareR2Connection struct {
	StorageBucket
	StorageCredentials
	// CloudflareAccountID The 32-character hexadecimal Cloudflare account ID that owns the
	// bucket.
	CloudflareAccountID string `json:"cloudflare_account_id"`
	// CloudflareJurisdiction is eu or fedramp. Nil uses the default jurisdiction.
	CloudflareJurisdiction *string `json:"cloudflare_jurisdiction"`
}

func (*CloudflareR2Connection) storageConnectionRequest() {}

// MarshalJSON includes the cloudflare_r2 provider discriminator.
func (v CloudflareR2Connection) MarshalJSON() ([]byte, error) {
	type fields CloudflareR2Connection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{CloudflareR2, fields(v)})
}
