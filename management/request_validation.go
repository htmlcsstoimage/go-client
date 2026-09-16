package management

import "fmt"

func validateOGRequest(request OGConfigRequest) error {
	switch r := request.(type) {
	case *HTMLCSSOGConfigRequest:
		if r != nil {
			return nil
		}
	case *TemplatedOGConfigRequest:
		if r != nil {
			return nil
		}
	}
	return fmt.Errorf("hcti: OG config request is required")
}

func validateStorageConnection(connection StorageConnectionRequest) error {
	switch c := connection.(type) {
	case *AWSS3Connection:
		if c != nil {
			return nil
		}
	case *CloudflareR2Connection:
		if c != nil {
			return nil
		}
	case *BackblazeB2Connection:
		if c != nil {
			return nil
		}
	case *DigitalOceanSpacesConnection:
		if c != nil {
			return nil
		}
	case *WasabiConnection:
		if c != nil {
			return nil
		}
	case *GoogleCloudStorageConnection:
		if c != nil {
			return nil
		}
	case *OtherS3CompatibleConnection:
		if c != nil {
			return nil
		}
	}
	return fmt.Errorf("hcti: storage connection is required")
}
