package management

import "encoding/json"

// BackblazeB2Connection configures BackblazeB2 storage.
type BackblazeB2Connection struct {
	StorageBucket
	StorageCredentials
	// Region Region code from the bucket's S3 endpoint, for example us-west-004.
	Region string `json:"region"`
}

func (*BackblazeB2Connection) storageConnectionRequest() {}

// MarshalJSON includes the backblaze_b2 provider discriminator.
func (v BackblazeB2Connection) MarshalJSON() ([]byte, error) {
	type fields BackblazeB2Connection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{BackblazeB2, fields(v)})
}
