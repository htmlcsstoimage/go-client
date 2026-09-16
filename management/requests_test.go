package management_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	hcti "github.com/htmlcsstoimage/go-client"
	m "github.com/htmlcsstoimage/go-client/management"
)

func jsonObject(t *testing.T, v any) map[string]json.RawMessage {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	return fields
}
func wantJSON(t *testing.T, fields map[string]json.RawMessage, key, want string) {
	t.Helper()
	if string(fields[key]) != want {
		t.Errorf("%s=%s, want %s", key, fields[key], want)
	}
}

func TestExplicitSecretRetention(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		secret                 *string
		retain                 *bool
		wantSecret, wantRetain string
	}{
		{"unspecified", nil, nil, "null", "null"},
		{"retain", nil, hcti.Ptr(true), "null", "true"},
		{"replace", hcti.Ptr("replacement"), nil, `"replacement"`, "null"},
		{"explicit false", hcti.Ptr("replacement"), hcti.Ptr(false), `"replacement"`, "false"},
		{"empty value", hcti.Ptr(""), nil, `""`, "null"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// The SDK preserves caller intent; the API validates operation-specific rules.
			proxy := jsonObject(t, m.ProxyAuthentication{Username: "", Password: tc.secret, RetainPassword: tc.retain})
			wantJSON(t, proxy, "username", `""`)
			wantJSON(t, proxy, "password", tc.wantSecret)
			wantJSON(t, proxy, "retain_password", tc.wantRetain)
			storage := jsonObject(t, m.StorageCredentials{AccessKeyID: "access", SecretAccessKey: tc.secret, RetainSecretAccessKey: tc.retain})
			wantJSON(t, storage, "secret_access_key", tc.wantSecret)
			wantJSON(t, storage, "retain_secret_access_key", tc.wantRetain)
		})
	}
	fields := jsonObject(t, &m.ProxyRequest{})
	for _, k := range []string{"port", "authentication", "bypass_hosts"} {
		wantJSON(t, fields, k, "null")
	}
	wantJSON(t, fields, "disabled", "false")
}

func TestStorageVariants(t *testing.T) {
	bucket := m.StorageBucket{Bucket: "images", KeyPrefix: hcti.Ptr("folder/💜")}
	creds := m.StorageCredentials{AccessKeyID: "access", SecretAccessKey: hcti.Ptr("secret")}
	for _, tc := range []struct {
		provider string
		value    m.StorageConnectionRequest
		extra    map[string]string
	}{
		{"aws_s3", &m.AWSS3Connection{StorageBucket: bucket, Region: "us-east-1", RoleARN: "role"}, map[string]string{"region": `"us-east-1"`, "role_arn": `"role"`}},
		{"cloudflare_r2", &m.CloudflareR2Connection{StorageBucket: bucket, StorageCredentials: creds, CloudflareAccountID: "account"}, map[string]string{"cloudflare_account_id": `"account"`, "cloudflare_jurisdiction": "null"}},
		{"backblaze_b2", &m.BackblazeB2Connection{StorageBucket: bucket, StorageCredentials: creds, Region: "us-west-004"}, map[string]string{"region": `"us-west-004"`}},
		{"digitalocean_spaces", &m.DigitalOceanSpacesConnection{StorageBucket: bucket, StorageCredentials: creds, Region: "nyc3"}, map[string]string{"region": `"nyc3"`}},
		{"wasabi", &m.WasabiConnection{StorageBucket: bucket, StorageCredentials: creds, Region: "us-east-1"}, map[string]string{"region": `"us-east-1"`}},
		{"google_cloud_storage", &m.GoogleCloudStorageConnection{StorageBucket: bucket, StorageCredentials: creds}, nil},
		{"other_s3_compatible", &m.OtherS3CompatibleConnection{StorageBucket: bucket, StorageCredentials: creds, Endpoint: "https://s3.example"}, map[string]string{"endpoint": `"https://s3.example"`, "region": "null", "force_path_style": "null"}},
	} {
		t.Run(tc.provider, func(t *testing.T) {
			// Verify embedding into the outer request doesn't lose fields to MarshalJSON promotion.
			outer := jsonObject(t, &m.StorageDestinationRequest{Name: "storage", ConnectionInfo: tc.value, HCTIStorageDisabled: true})
			wantJSON(t, outer, "name", `"storage"`)
			wantJSON(t, outer, "hcti_storage_disabled", "true")
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(outer["connection_info"], &fields); err != nil {
				t.Fatal(err)
			}
			wantJSON(t, fields, "provider", `"`+tc.provider+`"`)
			wantJSON(t, fields, "bucket", `"images"`)
			wantJSON(t, fields, "key_prefix", `"folder/💜"`)
			for k, v := range tc.extra {
				wantJSON(t, fields, k, v)
			}
			if tc.provider == "aws_s3" {
				if _, ok := fields["secret_access_key"]; ok {
					t.Error("AWS role connection contains secret")
				}
			} else {
				wantJSON(t, fields, "secret_access_key", `"secret"`)
				wantJSON(t, fields, "retain_secret_access_key", "null")
			}
			// Even if a server accidentally returns request-only fields, the read model cannot retain them.
			var info m.StorageConnectionInfo
			if err := json.Unmarshal(outer["connection_info"], &info); err != nil {
				t.Fatal(err)
			}
			output := jsonObject(t, info)
			for _, key := range []string{"secret_access_key", "retain_secret_access_key"} {
				if _, ok := output[key]; ok {
					t.Errorf("read model retained %s", key)
				}
			}
		})
	}
}

func TestOGVariantsAndDefaults(t *testing.T) {
	common := m.OGConfigOptions{Name: "social", BaseURL: "https://example.com"}
	html := jsonObject(t, &m.HTMLCSSOGConfigRequest{OGConfigOptions: common, DefaultOptions: &m.OGDefaultImageOptions{MSDelay: hcti.Ptr(int32(0)), TransparentBackground: hcti.Ptr(false)}})
	wantJSON(t, html, "config_type", `"html_css"`)
	wantJSON(t, html, "name", `"social"`)
	wantJSON(t, html, "description", "null")
	wantJSON(t, html, "extract_values", "false")
	for _, k := range []string{"optimization_mode", "refresh_interval_s", "template_id"} {
		if _, ok := html[k]; ok {
			t.Errorf("unexpected %s", k)
		}
	}
	var opts map[string]json.RawMessage
	if err := json.Unmarshal(html["default_options"], &opts); err != nil {
		t.Fatal(err)
	}
	wantJSON(t, opts, "device_scale", "null")
	wantJSON(t, opts, "ms_delay", "0")
	wantJSON(t, opts, "transparent_background", "false")
	for _, k := range []string{"html", "google_fonts", "pdf_options", "jumbo", "format", "dedupe_duration_s"} {
		if _, ok := opts[k]; ok {
			t.Errorf("unsupported OG option %s", k)
		}
	}
	templated := jsonObject(t, &m.TemplatedOGConfigRequest{OGConfigOptions: common, TemplateID: "t-example", TemplateValuesMapping: []m.OGTemplateValueMapping{{TemplateKey: "title", Fallback: hcti.Ptr(m.Titles)}, {TemplateKey: "description", MetaKey: hcti.Ptr("og:description")}}, Headers: map[string]string{}, AdditionalHeaderOrigins: []string{}})
	wantJSON(t, templated, "config_type", `"templated"`)
	wantJSON(t, templated, "template_id", `"t-example"`)
	wantJSON(t, templated, "template_version", "null")
	wantJSON(t, templated, "headers", `{}`)
	wantJSON(t, templated, "additional_header_origins", `[]`)
	wantJSON(t, templated, "template_values_mapping", `[{"template_key":"title","meta_key":null,"fallback":"titles"},{"template_key":"description","meta_key":"og:description","fallback":null}]`)
	if _, ok := templated["default_options"]; ok {
		t.Error("template config contains HTML options")
	}
}

func TestNumericResponseUnions(t *testing.T) {
	var config m.OGConfig
	if err := json.Unmarshal([]byte(`{"config_type":"html_css","template_version":"9223372036854775807","refresh_interval_s":86400,"default_options":{"device_scale":"1.5","viewport_width":"1200","ms_delay":0}}`), &config); err != nil {
		t.Fatal(err)
	}
	if config.TemplateVersion == nil || *config.TemplateVersion != 9223372036854775807 || config.DefaultOptions.DeviceScale == nil || *config.DefaultOptions.DeviceScale != 1.5 || *config.DefaultOptions.ViewportWidth != 1200 || *config.DefaultOptions.MSDelay != 0 {
		t.Fatalf("lost numeric values: %+v", config)
	}
	for _, body := range []string{`{"refresh_interval_s":-1}`, `{"template_version":"9223372036854775808"}`, `{"default_options":{"viewport_width":"1.5"}}`, `{"default_options":{"ms_delay":true}}`} {
		if err := json.Unmarshal([]byte(body), &config); err == nil {
			t.Errorf("accepted %s", body)
		}
	}
	if err := json.Unmarshal([]byte(`{"template_version":null,"default_options":null}`), &config); err != nil {
		t.Fatal(err)
	}
	if config.TemplateVersion != nil || config.DefaultOptions != nil {
		t.Error("reused response retained old fields")
	}
}

func TestUnicodeAndMalformedUTF8(t *testing.T) {
	for _, value := range []string{"emoji 👩🏽‍💻 💜", "日本語 e\u0301", "\x00\n\t\"\\<>&", "bad\xff\xfeutf8", "\xed\xa0\x80"} {
		want, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		fields := jsonObject(t, m.ProxyAuthentication{Username: value, Password: &value})
		wantJSON(t, fields, "username", string(want))
		wantJSON(t, fields, "password", string(want))
		opts := jsonObject(t, m.OGDefaultImageOptions{CSS: &value, Headers: map[string]string{"X-Test": value}})
		wantJSON(t, opts, "css", string(want))
		var headers map[string]string
		if err := json.Unmarshal(opts["headers"], &headers); err != nil {
			t.Fatal(err)
		}
		var normalized string
		if err := json.Unmarshal(want, &normalized); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(headers, map[string]string{"X-Test": normalized}) {
			t.Fatal("header serialization differs from encoding/json")
		}
	}
	// Nil and explicitly empty permissions must remain distinct.
	wantJSON(t, jsonObject(t, &m.APIKeyRequest{}), "permissions", "null")
	wantJSON(t, jsonObject(t, &m.APIKeyRequest{Permissions: []m.Permission{}}), "permissions", "[]")
	var proxy m.Proxy
	if err := json.Unmarshal([]byte(strings.TrimSuffix(proxyJSON, "}")+`,"password":"secret"}`), &proxy); err != nil {
		t.Fatal(err)
	}
	if _, ok := jsonObject(t, proxy)["password"]; ok {
		t.Error("read model retained password")
	}
}
