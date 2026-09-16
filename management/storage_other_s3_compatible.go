package management

import "encoding/json"

// OtherS3CompatibleConnection configures OtherS3Compatible storage.
type OtherS3CompatibleConnection struct {
	StorageBucket
	StorageCredentials
	// Endpoint Public HTTPS endpoint of the S3-compatible service. Do not include a path, query
	// string, fragment, or credentials.
	Endpoint string `json:"endpoint"`
	// Region Signing region expected by the service. Omit or set to null to use us-east-1.
	Region *string `json:"region"`
	// ForcePathStyle Whether to put the bucket name in the URL path instead of the hostname.
	// Omitted or null defaults to true.
	ForcePathStyle *bool `json:"force_path_style"`
}

func (*OtherS3CompatibleConnection) storageConnectionRequest() {}

// MarshalJSON includes the other_s3_compatible provider discriminator.
func (v OtherS3CompatibleConnection) MarshalJSON() ([]byte, error) {
	type fields OtherS3CompatibleConnection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{OtherS3Compatible, fields(v)})
}
