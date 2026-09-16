package management

import (
	"context"
	"net/http"
	"testing"
)

type rejectingTransport struct{ t *testing.T }

func (r rejectingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	r.t.Fatal("nil request reached transport")
	return nil, nil
}

func TestRequestInterfacesRequirePointers(t *testing.T) {
	for _, v := range []any{HTMLCSSOGConfigRequest{}, TemplatedOGConfigRequest{}} {
		if _, ok := v.(OGConfigRequest); ok {
			t.Fatalf("%T unexpectedly implements OGConfigRequest", v)
		}
	}
	for _, v := range []any{AWSS3Connection{}, CloudflareR2Connection{}, BackblazeB2Connection{}, DigitalOceanSpacesConnection{}, WasabiConnection{}, GoogleCloudStorageConnection{}, OtherS3CompatibleConnection{}} {
		if _, ok := v.(StorageConnectionRequest); ok {
			t.Fatalf("%T unexpectedly implements StorageConnectionRequest", v)
		}
	}
}

func TestNilRequestsBeforeHTTP(t *testing.T) {
	c := NewClient("id", "key", WithHTTPClient(&http.Client{Transport: rejectingTransport{t}}))
	ctx := context.Background()
	checks := map[string]func() error{
		"create APIKey":             func() error { _, e := c.CreateAPIKey(ctx, nil); return e },
		"update APIKey":             func() error { _, e := c.UpdateAPIKey(ctx, "id", nil); return e },
		"create Proxy":              func() error { _, e := c.CreateProxy(ctx, nil); return e },
		"update Proxy":              func() error { _, e := c.UpdateProxy(ctx, "id", nil); return e },
		"create StorageDestination": func() error { _, e := c.CreateStorageDestination(ctx, nil); return e },
		"update StorageDestination": func() error { _, e := c.UpdateStorageDestination(ctx, "id", nil); return e },
		"image":                     func() error { _, e := c.CreateImageDefinition(ctx, "html_css", nil); return e },
		"template":                  func() error { _, e := c.SaveTemplateDefinition(ctx, "", nil); return e },
	}
	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			if check() == nil {
				t.Fatal("accepted nil request")
			}
		})
	}
	for _, r := range []OGConfigRequest{nil, (*HTMLCSSOGConfigRequest)(nil), (*TemplatedOGConfigRequest)(nil)} {
		if _, err := c.CreateOGConfig(ctx, r); err == nil {
			t.Fatalf("accepted nil OG request %T", r)
		}
		if _, err := c.UpdateOGConfig(ctx, "id", r); err == nil {
			t.Fatalf("accepted nil OG request %T", r)
		}
	}
	for _, conn := range []StorageConnectionRequest{nil, (*AWSS3Connection)(nil), (*CloudflareR2Connection)(nil), (*BackblazeB2Connection)(nil), (*DigitalOceanSpacesConnection)(nil), (*WasabiConnection)(nil), (*GoogleCloudStorageConnection)(nil), (*OtherS3CompatibleConnection)(nil)} {
		r := &StorageDestinationRequest{ConnectionInfo: conn}
		if _, err := c.CreateStorageDestination(ctx, r); err == nil {
			t.Fatalf("accepted nil connection %T", conn)
		}
		if _, err := c.UpdateStorageDestination(ctx, "id", r); err == nil {
			t.Fatalf("accepted nil connection %T", conn)
		}
	}
}
