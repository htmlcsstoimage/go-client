package management

import (
	"context"
	"net/http"
	"time"
)

// StorageDestinationRequest creates a destination or replaces its editable settings.
type StorageDestinationRequest struct {
	// Name is a display name of 3–255 characters.
	Name string `json:"name"`
	// Disabled prevents this destination from being used for new renders.
	Disabled bool `json:"disabled"`
	// HCTIStorageDisabled stores output only at this destination; rendering then requires authenticated /store requests.
	HCTIStorageDisabled bool `json:"hcti_storage_disabled"`
	// ConnectionInfo selects the storage provider and supplies its connection settings.
	ConnectionInfo StorageConnectionRequest `json:"connection_info"`
}

// StorageDestination describes a storage destination without disclosing credentials.
type StorageDestination struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Enabled             bool   `json:"enabled"`
	HCTIStorageDisabled bool   `json:"hcti_storage_disabled"`
	// Provider is a display label. Use ConnectionInfo.Provider for the machine identifier.
	Provider          string                `json:"provider"`
	ConnectionInfo    StorageConnectionInfo `json:"connection_info"`
	LastTestedAt      *time.Time            `json:"last_tested_at"`
	LastTestSucceeded *bool                 `json:"last_test_succeeded"`
	LastTestError     *string               `json:"last_test_error"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

func (v *StorageDestination) validateResponse() string {
	if v.ID == "" || v.Name == "" {
		return "missing storage destination metadata"
	}
	return v.ConnectionInfo.validateResponse()
}

// AWSExternalID contains the organization's external ID for an AWS role trust policy.
type AWSExternalID struct {
	ExternalID string `json:"external_id"`
}

func (v *AWSExternalID) validateResponse() string {
	if v.ExternalID == "" {
		return "missing AWS external ID"
	}
	return ""
}

// GetAWSExternalID reads the organization's external ID without creating a destination.
func (c *Client) GetAWSExternalID(ctx context.Context) (*AWSExternalID, error) {
	var result AWSExternalID
	if err := c.do(ctx, http.MethodGet, "/v1/storage-destinations/aws-external-id", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
