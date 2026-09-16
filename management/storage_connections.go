package management

// StorageProvider identifies a storage connection type.
type StorageProvider string

// Supported storage providers.
const (
	AWSS3              StorageProvider = "aws_s3"
	CloudflareR2       StorageProvider = "cloudflare_r2"
	BackblazeB2        StorageProvider = "backblaze_b2"
	DigitalOceanSpaces StorageProvider = "digitalocean_spaces"
	Wasabi             StorageProvider = "wasabi"
	GoogleCloudStorage StorageProvider = "google_cloud_storage"
	OtherS3Compatible  StorageProvider = "other_s3_compatible"
)

// StorageConnectionRequest accepts a non-nil pointer to a provider-specific connection request.
// Each concrete type supplies its provider discriminator automatically.
type StorageConnectionRequest interface{ storageConnectionRequest() }

// StorageBucket selects an existing bucket; the API does not create buckets.
type StorageBucket struct {
	// Bucket is the existing bucket or Space name.
	Bucket string `json:"bucket"`
	// KeyPrefix prefixes object keys. Nil uses the bucket root.
	KeyPrefix *string `json:"key_prefix"`
}

// StorageCredentials supplies new credentials or explicitly retains the stored secret.
type StorageCredentials struct {
	// AccessKeyID identifies the credentials; Google Cloud Storage requires an HMAC access ID.
	AccessKeyID string `json:"access_key_id"`
	// SecretAccessKey is required on create or when the provider or access key changes.
	// An empty secret is invalid. Nil is allowed only with RetainSecretAccessKey=true.
	SecretAccessKey *string `json:"secret_access_key"`
	// RetainSecretAccessKey=true is update-only and requires unchanged provider and
	// AccessKeyID, existing credentials, and nil SecretAccessKey. Nil or false requires a secret.
	RetainSecretAccessKey *bool `json:"retain_secret_access_key"`
}

// StorageConnectionInfo is the readable connection metadata. Fields apply only to
// their Provider; secrets and retention flags are never included in this type.
type StorageConnectionInfo struct {
	Provider StorageProvider `json:"provider"`
	StorageBucket
	Region                 *string `json:"region"`
	RoleARN                *string `json:"role_arn"`
	AccessKeyID            *string `json:"access_key_id"`
	CloudflareAccountID    *string `json:"cloudflare_account_id"`
	CloudflareJurisdiction *string `json:"cloudflare_jurisdiction"`
	Endpoint               *string `json:"endpoint"`
	ForcePathStyle         *bool   `json:"force_path_style"`
}

func (v *StorageConnectionInfo) validateResponse() string {
	if v.Bucket == "" {
		return "missing storage bucket"
	}
	nonempty := func(s *string) bool { return s != nil && *s != "" }
	switch v.Provider {
	case AWSS3:
		if !nonempty(v.Region) || !nonempty(v.RoleARN) {
			return "missing AWS connection metadata"
		}
	case CloudflareR2:
		if !nonempty(v.CloudflareAccountID) || !nonempty(v.AccessKeyID) {
			return "missing R2 connection metadata"
		}
	case BackblazeB2, DigitalOceanSpaces, Wasabi:
		if !nonempty(v.Region) || !nonempty(v.AccessKeyID) {
			return "missing storage connection metadata"
		}
	case GoogleCloudStorage:
		if !nonempty(v.AccessKeyID) {
			return "missing storage access key ID"
		}
	case OtherS3Compatible:
		if !nonempty(v.Endpoint) || !nonempty(v.AccessKeyID) {
			return "missing S3-compatible connection metadata"
		}
	default:
		return "unknown storage provider"
	}
	return ""
}
