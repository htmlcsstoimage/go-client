package management

import "encoding/json"

// WasabiConnection configures Wasabi storage.
type WasabiConnection struct {
	StorageBucket
	StorageCredentials
	// Region Region containing the bucket, for example us-east-1.
	Region string `json:"region"`
}

func (*WasabiConnection) storageConnectionRequest() {}

// MarshalJSON includes the wasabi provider discriminator.
func (v WasabiConnection) MarshalJSON() ([]byte, error) {
	type fields WasabiConnection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{Wasabi, fields(v)})
}
