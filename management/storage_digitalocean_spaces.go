package management

import "encoding/json"

// DigitalOceanSpacesConnection configures DigitalOceanSpaces storage.
type DigitalOceanSpacesConnection struct {
	StorageBucket
	StorageCredentials
	// Region Region containing the Space, for example nyc3.
	Region string `json:"region"`
}

func (*DigitalOceanSpacesConnection) storageConnectionRequest() {}

// MarshalJSON includes the digitalocean_spaces provider discriminator.
func (v DigitalOceanSpacesConnection) MarshalJSON() ([]byte, error) {
	type fields DigitalOceanSpacesConnection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{DigitalOceanSpaces, fields(v)})
}
