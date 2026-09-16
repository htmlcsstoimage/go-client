package management

import "encoding/json"

// GoogleCloudStorageConnection configures GoogleCloudStorage storage.
type GoogleCloudStorageConnection struct {
	StorageBucket
	StorageCredentials
}

func (*GoogleCloudStorageConnection) storageConnectionRequest() {}

// MarshalJSON includes the google_cloud_storage provider discriminator.
func (v GoogleCloudStorageConnection) MarshalJSON() ([]byte, error) {
	type fields GoogleCloudStorageConnection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{GoogleCloudStorage, fields(v)})
}
