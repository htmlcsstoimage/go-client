// Package management provides clients for API keys, proxies, storage destinations,
// Open Graph configurations, and provider-oriented rendering definitions.
//
// This package is experimental and intended primarily for the official Terraform
// and Pulumi providers. It is publicly importable, but is not recommended for
// general application use yet. Its API may change as the IaC integrations develop.
// Use the root hcti package for application image creation, templates, and signing.
// Both packages share authentication and HTTP behavior.
// Pass request payloads and storage connections as pointers. Nil requests are
// rejected, and client methods do not modify requests. List options remain values.
//
// Updates replace resource settings. Nullable fields are sent as JSON null when
// unset; the API applies its defaults. Retaining credentials requires an explicit
// retention flag. Requests are never automatically retried.
package management
