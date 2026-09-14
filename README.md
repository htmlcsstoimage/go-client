# HTML/CSS to Image — Go client

A Go client for the HTML/CSS to Image API, with support for Google Fonts, PDF options, templates, and signed rendering URLs.

## Installation

Requires Go 1.21 or newer.

```sh
go get github.com/htmlcsstoimage/go-client
```

Import `github.com/htmlcsstoimage/go-client` and use the `hcti` package name, as shown below.

## Create an image

```go
package main

import (
    "context"
    "log"

    "github.com/htmlcsstoimage/go-client"
)

func main() {
    client, err := hcti.NewClientFromEnv() // HCTI_API_ID and HCTI_API_KEY
    if err != nil { log.Fatal(err) }

    image, err := client.CreateImage(context.Background(), hcti.HTMLImageRequest{
        HTML: "<h1>Hello from Go</h1>",
        CSS: hcti.Ptr("h1 { font-family: 'Open Sans'; }"),
        GoogleFonts: hcti.GoogleFonts{"Open Sans", "Roboto"},
        ImageOptions: hcti.ImageOptions{
            Format: hcti.PNG,
            RenderOptions: hcti.RenderOptions{
                TransparentBackground: hcti.Ptr(false),
            },
        },
    })
    if err != nil { log.Fatal(err) }
    log.Println(image.URL)
}
```

Alternatively, use `hcti.NewClient(apiID, apiKey)`. Keep API credentials on the server.

Optional render fields use pointers: nil preserves API defaults; `hcti.Ptr(false)`, `hcti.Ptr(0)`, and `hcti.Ptr("")` send explicit values. For typed options, use `hcti.Ptr(hcti.Dark)` or `hcti.Ptr(hcti.Print)`. Font names are trimmed, deduplicated, and serialized to the API's pipe-delimited format.

## URL screenshots and PDF options

```go
image, err := client.CreateImage(ctx, hcti.URLImageRequest{
    URL: "https://example.com",
    FullScreen: hcti.Ptr(true),
    ImageOptions: hcti.ImageOptions{
        Format: hcti.PDF,
        PDFOptions: &hcti.PDFOptions{
            PrintBackground: hcti.Ptr(true),
            PageWidth: &hcti.PDFLength{Value: 8.5, Unit: hcti.Inches},
            PageHeight: &hcti.PDFLength{Value: 11, Unit: hcti.Inches},
            Margins: &hcti.PDFMargins{
                Top: hcti.PDFLength{Value: 10, Unit: hcti.Millimeters},
                Bottom: hcti.PDFLength{Value: 10, Unit: hcti.Millimeters},
            },
        },
    },
})
```

PDF units support pixels, inches, centimeters, and millimeters. Zero-value lengths are zero pixels. Format selects the returned URL's extension; it does not restrict the stored image to that format.

## Templates

```go
version, err := client.CreateTemplate(ctx, hcti.TemplateRequest{
    Name: "Social card",
    HTML: "<h1>{{title}}</h1>",
    GoogleFonts: hcti.GoogleFonts{"Roboto"},
})
if err != nil { return err }

image, err := client.CreateImage(ctx, hcti.TemplatedImageRequest{
    TemplateID: version.TemplateID,
    TemplateVersion: hcti.Ptr(version.TemplateVersion),
    TemplateValues: map[string]any{"title": "Hello"},
})
```

Use `CreateTemplateVersion` to add a version to a stable template ID. `ListTemplates` returns the latest versions; `ListTemplateVersions` lists versions of one template. Pass `page.Pagination.NextPageStart` as `TemplateListOptions.MaxVersion` to fetch subsequent pages; stop when nil. `DeleteTemplate` removes the entire template, not an individual version.

Listing returns common fields plus `Template.Raw`, which retains the complete JSON for template-editor block variants not yet modeled by the client.

## Batches and deletion

`CreateImageBatch` accepts a `BatchRequest` containing HTML/URL `Variations` and optional `DefaultOptions`. Omitted HTML/URL fields let variations inherit defaults. Templated image batches are not supported. Empty batches return locally without making an HTTP request. Deduplication settings are excluded from defaults and every variation without modifying the input requests.

A nil CSS pointer, map, or slice inherits the batch default. Use `CSS: hcti.Ptr("")`, `Headers: map[string]string{}`, `AdditionalHeaderOrigins: []string{}`, or `GoogleFonts: hcti.GoogleFonts{}` to explicitly clear that default:

```go
batch, err := client.CreateImageBatch(ctx, hcti.BatchRequest{
    DefaultOptions: hcti.HTMLImageRequest{
        HTML: "<h1>Hello</h1>",
        CSS: hcti.Ptr("h1 { color: red; }"),
        GoogleFonts: hcti.GoogleFonts{"Roboto"},
    },
    Variations: []hcti.ImageRequest{
        hcti.HTMLImageRequest{}, // Inherits CSS and fonts.
        hcti.HTMLImageRequest{CSS: hcti.Ptr(""), GoogleFonts: hcti.GoogleFonts{}},
    },
})
```

Use `DeleteImage` or `DeleteImageBatch` to delete images. Successful deletion returns a nil error, including responses with no body.

## Signed URLs

`GenerateTemplatedImageURL` and `GenerateCreateAndRenderURL` produce URLs locally, without API calls. They return `(string, error)` and accept the same templated/URL request types used for image creation. `GenerateCreateAndRenderURL` excludes PDF options and dedupe duration. Signed URLs authorize rendering; do not include confidential template values or request headers in URLs you share.

## Image URLs, resizing, and cropping

`ImageURL` builds a URL for an existing image without calling the API. `RenderImageOptions` controls its format, dimensions, DPI, and crop. Cropping happens before resizing.

```go
options := hcti.RenderImageOptions{
    Format: hcti.WebP,
    Width: hcti.Ptr(600),
    Crop: &hcti.Crop{
        Horizontal: &hcti.CropSpan{
            Size: &hcti.CropValue{Value: 100, Unit: hcti.CropPercent},
        },
        AspectRatio: &hcti.AspectRatio{Width: 16, Height: 9},
        AspectRatioAxis: hcti.CropWidth,
        ComputedOrigin: hcti.CropCenter,
    },
}
imageURL, err := client.ImageURL("existing-image-id", options)
```

Pass the same options as the optional second argument to `GenerateTemplatedImageURL(request, options)` or `GenerateCreateAndRenderURL(request, options)`. A nonempty options format overrides the request format. Template variables that share transformation parameter names are preserved using the API's `__ro_` override parameters.

For rectangles, supply `Horizontal` and/or `Vertical` spans without an aspect ratio. Spans support a start position, start/end boundaries, start/size, or a size anchored with `CropStart`, `CropCenter`, or `CropEnd`. Positions and sizes accept pixels or percentages. URL helpers validate dimensions, units, crop combinations, and origins before returning a URL.

## HTTP behavior and errors

Every network operation takes a `context.Context`. The default timeout is 60 seconds. `WithHTTPClient` accepts custom transport/timeout settings; `WithBaseURL` selects another API origin. Redirects are disabled. The client never automatically retries requests, since retrying an ambiguous create can create another resource.

API requests send `User-Agent: HCTIGo/<version>`. The version is embedded from `VERSION` when the client is compiled.

Non-2xx responses return `*hcti.APIError`, inspectable with `errors.As`. It contains the HTTP status, API code/message, validation errors, and response headers (including rate-limit information). `Error()` intentionally omits response contents. Malformed successful responses return `*hcti.ResponseError`, including missing image IDs/URLs, invalid template versions, malformed arrays, and unexpected empty responses. Its message excludes response contents. Transport errors preserve their causes through Go error wrapping.

## Runnable examples

See [examples](examples/README.md) for complete programs covering HTML images, URL screenshots, PDFs, batches, template versioning, image transformations, and signed URLs.

## Development

### Releases

Update the root `VERSION` file (for example, from `0.1.0` to `0.1.1`) and push the change to `main`. After the Test workflow passes, the Release workflow creates the matching `v`-prefixed Git tag and a GitHub release with generated notes, then requests Go module indexing. It releases the exact commit that passed tests.

Existing tags and releases are left unchanged. If a release run fails after creating its tag, rerun the failed workflow to finish the release and indexing. No package registry secrets are required; the workflow uses GitHub's built-in token.

`VERSION` accepts stable semantic versions without a `v` prefix. Use a new version for every release; published tags must not be moved. Major versions starting with v2 require the matching `/v2` (or later) suffix in `go.mod` and imports.

### Checks

```sh
go test ./...
go vet ./...
```

Tests use fake HTTP transports and require no API credentials or live service.

Encoding tests cover emoji sequences, multilingual text, malformed UTF-8, control bytes, long values, buffer growth, and signature verification. Direct query encoding preserves arbitrary Go string bytes through percent-escaping, matching `net/url`. Template values preserve `encoding/json` semantics, including replacing each malformed UTF-8 byte in strings with U+FFFD. Ordinary scalar values are encoded directly; structured values and custom marshalers use `encoding/json`.

Run the encoding and signing fuzz targets separately:

```sh
go test -run '^$' -fuzz '^FuzzAppendQueryEscaped$' -fuzztime=30s
go test -run '^$' -fuzz '^FuzzScreenshotSigning$' -fuzztime=30s
go test -run '^$' -fuzz '^FuzzTemplateSigning$' -fuzztime=30s
go test -run '^$' -fuzz '^FuzzTemplateFloatEncoding$' -fuzztime=30s
```

Run the URL generation benchmarks:

```sh
go test -run '^$' -bench BenchmarkGenerate -benchmem
```
