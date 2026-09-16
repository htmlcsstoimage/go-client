package management

import "encoding/json"

// AWSS3Connection configures AWSS3 storage.
type AWSS3Connection struct {
	StorageBucket
	// Region AWS region containing the bucket, for example us-east-1.
	Region string `json:"region"`
	// RoleARN ARN of the IAM role HCTI assumes to access the bucket. Configure its trust policy
	// for this organization before saving.
	RoleARN string `json:"role_arn"`
}

func (*AWSS3Connection) storageConnectionRequest() {}

// MarshalJSON includes the aws_s3 provider discriminator.
func (v AWSS3Connection) MarshalJSON() ([]byte, error) {
	type fields AWSS3Connection
	return json.Marshal(struct {
		Provider StorageProvider `json:"provider"`
		fields
	}{AWSS3, fields(v)})
}
