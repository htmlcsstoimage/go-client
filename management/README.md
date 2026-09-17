# Management API

**Experimental:** this package is intended primarily for the official Terraform and Pulumi providers. It is publicly importable, but is not recommended for general application use yet. Its API may change as the IaC integrations develop. Use the root `hcti` package for application image creation, templates, and signing.

Manage API keys, proxies, storage destinations, and Open Graph configurations:

```go
import (
    hcti "github.com/htmlcsstoimage/go-client"
    "github.com/htmlcsstoimage/go-client/management"
)

client, err := management.NewClientFromEnv() // HCTI_API_ID and HCTI_API_KEY
if err != nil { return err }

page, err := client.ListProxies(ctx, management.ListOptions{Count: hcti.Ptr(100)})
if err != nil { return err }
for _, proxy := range page.Data {
    fmt.Println(proxy.ID, proxy.Name)
}
```

Both packages belong to the same Go module and use the same version and release. The shared `internal/httpapi` package is an implementation detail; applications and providers never import it directly.

Pass request payloads and storage connections as pointers, such as
`&management.ProxyRequest{...}` or `&management.AWSS3Connection{...}`.
Client methods reject nil requests and do not modify them. List options remain values.

See the [runnable pagination example](../examples/management/main.go) and [package reference](https://pkg.go.dev/github.com/htmlcsstoimage/go-client/management).

## Operations

| Resource | Methods |
| --- | --- |
| API keys | `CreateAPIKey`, `GetAPIKey`, `ListAPIKeys`, `UpdateAPIKey`, `DeleteAPIKey` |
| Proxies | `CreateProxy`, `GetProxy`, `ListProxies`, `UpdateProxy`, `DeleteProxy` |
| Storage destinations | `CreateStorageDestination`, `GetStorageDestination`, `ListStorageDestinations`, `UpdateStorageDestination`, `DeleteStorageDestination`, `GetAWSExternalID` |
| OG configs | `CreateOGConfig`, `GetOGConfig`, `ListOGConfigs`, `UpdateOGConfig`, `DeleteOGConfig` |

List methods fetch one page. Pass `Pagination.NextPageStart` as the next request's `PageStart`; nil ends pagination. API key lists also support `IncludeDisabled` and an all-of `WithPermission` filter.

## Updates and secrets

Updates replace editable settings. Send the complete desired configuration. Nullable request fields serialize nil as JSON null, while explicit false, zero, empty strings, maps, and slices remain distinct. Nonnullable optional OG settings (`OptimizationMode`, `RefreshIntervalSeconds`) are omitted when nil so the API chooses their defaults. The SDK does not fill rendering defaults or automatically retain secrets.

To retain a proxy password during an update, supply the unchanged username and explicitly set `RetainPassword`:

```go
auth := &management.ProxyAuthentication{
    Username: "proxy-user",
    RetainPassword: hcti.Ptr(true),
}
```

To replace it, supply `Password` and leave `RetainPassword` nil or false. An explicit empty proxy password is valid. Nil authentication removes proxy credentials.

Storage connections use `RetainSecretAccessKey` with the same explicit protocol: retaining requires unchanged provider and access key ID, existing credentials, and a nil `SecretAccessKey`. An empty storage secret is invalid. The API validates these rules, including create-versus-update restrictions; the client preserves the values you supply.

API key creation returns `APIKeyWithSecret.Secret` only once. Read/update responses cannot recover it. `DeleteAPIKey` disables the key, and its metadata remains readable. Use `APIKey.ID` for management operations and `APIKey.APIID` for authentication.

Read models contain no proxy password or storage secret fields. Header values and API error detail fields can still contain sensitive data; avoid logging entire objects. `APIError.Error()` and `ResponseError.Error()` exclude response contents.

## Storage and OG configuration types

`StorageDestinationRequest.ConnectionInfo` accepts `AWSS3Connection`, `CloudflareR2Connection`, `BackblazeB2Connection`, `DigitalOceanSpacesConnection`, `WasabiConnection`, `GoogleCloudStorageConnection`, or `OtherS3CompatibleConnection`. Each sets its own `provider` discriminator. AWS S3 uses a role ARN; other providers use `StorageCredentials`. All use `StorageBucket` to select an existing bucket and optional key prefix.

Storage creates and updates may test connectivity. Reads do not. Inspect the returned connection-test fields and actual settings, particularly when disabling a destination after a plan change.

OG writes accept `HTMLCSSOGConfigRequest` or `TemplatedOGConfigRequest`, with common settings in `OGConfigOptions`. Each sets its own `config_type` discriminator. HTML/CSS supports only the options listed in `OGDefaultImageOptions`. Template ID and optional version are JSON body fields for OG configs; a nil version follows the latest template. Updating an OG config can change its configuration type without changing its management ID.

## Transport and errors

Use `NewClient(apiID, apiKey)` or `NewClientFromEnv`. `WithHTTPClient` supplies transport/timeouts, `WithBaseURL` changes the API origin, and `WithUserAgentSuffix` appends an integration name/version while retaining the SDK identifier. To reuse a connection pool with the image client, supply the same `http.Client` to each package's `WithHTTPClient` option.

Both clients use HTTP Basic authentication and `HCTIGo/<version>`, with the version embedded from `VERSION`. The default timeout is 60 seconds. Redirects and automatic retries are disabled: a failed write response can follow a successful server-side change.

Errors support `errors.As` with either package's `APIError` and `ResponseError` aliases. Context cancellation is preserved through `errors.Is`. Successful responses are checked for required identity fields and recognized discriminators; malformed responses return an error rather than an empty resource.

## Rendering definitions for infrastructure providers

`CreateImageDefinition` supports `html_css`, `url`, and `templated` with typed `RenderDefinition` fields. It saves a definition without rendering, sends nullable inputs as explicit null, and disables deduplication so each create owns a distinct image. `GetImageMetadata` reads `/v1/images/{id}`; `DeleteImageDefinition` accepts asynchronous deletion. Template selectors are URL parameters, not request body fields.

`SaveTemplateDefinition` creates a template when its ID argument is empty, or creates a version under an existing ID. `GetTemplateDefinition` accepts nil for latest or searches the version listing for an exact int64 version. `DeleteTemplateDefinition` removes all versions. These methods share the existing transport and authentication. Use the root package for ordinary application image creation, templates, and signing.

`GetAWSExternalID` returns `ExternalID` for the trust policy's `sts:ExternalId` condition and `WriterRoleARN` for its AWS principal. The storage destination's role ARN is the role you create in your own AWS account.
