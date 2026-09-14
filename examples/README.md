# Examples

Run these commands from the repository root. Each directory contains a complete program using the local client module.

Set `HCTI_API_ID` and `HCTI_API_KEY` in your environment before running, except for the image-URL example, which needs no credentials. The signed-URL example optionally reads `HCTI_TEMPLATE_ID` for an existing template.

| Example | Run | Demonstrates |
| --- | --- | --- |
| [HTML image](html-image/main.go) | `go run ./examples/html-image` | HTML/CSS card, Google Fonts, device scale |
| [URL screenshot](url-screenshot/main.go) | `go run ./examples/url-screenshot` | Full-page screenshot, viewport dimensions, WebP |
| [PDF](pdf/main.go) | `go run ./examples/pdf` | A4 report, millimeter dimensions, margins, print backgrounds |
| [Batch](batch/main.go) | `go run ./examples/batch` | Three color variations sharing HTML and font defaults |
| [Templates](templates/main.go) | `go run ./examples/templates` | Create a template, add a version, render that version, paginate versions |
| [Image URL](image-url/main.go) | `go run ./examples/image-url IMAGE_ID` | Crop an existing image to 16:9, resize, and select WebP locally |
| [Signed URLs](signed-urls/main.go) | `go run ./examples/signed-urls` | Generate resized WebP screenshot and template URLs locally |

The image examples call the API and can consume image credits. The template example creates a new template and two versions on each run. Created resources are kept so you can inspect the printed URLs; use `DeleteImage` or `DeleteTemplate` when you no longer need them.

Use credentials with image-creation permission for the image examples. The template example additionally needs template create/update and read permissions. Signed URLs are generated without API calls; opening them requests rendering. Treat the printed signed URLs as shareable rendering credentials.

Programs print render URLs rather than downloading files. They stop on errors and use a two-minute context deadline for API workflows, alongside the client's default per-request timeout.

To compile all examples and run the offline SDK tests without making API calls:

```sh
go test ./...
```
