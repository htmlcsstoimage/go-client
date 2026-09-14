package main

import (
	"fmt"
	"log"
	"os"

	"github.com/htmlcsstoimage/go-client"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client, err := hcti.NewClientFromEnv()
	if err != nil {
		return err
	}
	renderOptions := hcti.RenderImageOptions{Width: hcti.Ptr(600), Format: hcti.WebP}

	screenshotURL, err := client.GenerateCreateAndRenderURL(hcti.URLImageRequest{
		URL:          "https://example.com",
		FullScreen:   hcti.Ptr(true),
		ImageOptions: hcti.ImageOptions{Format: hcti.PNG},
	}, renderOptions)
	if err != nil {
		return fmt.Errorf("sign screenshot URL: %w", err)
	}
	fmt.Println("Screenshot URL:", screenshotURL)

	// Optional: use a template belonging to the same organization as the credentials.
	if templateID := os.Getenv("HCTI_TEMPLATE_ID"); templateID != "" {
		templateURL, err := client.GenerateTemplatedImageURL(hcti.TemplatedImageRequest{
			TemplateID:     templateID,
			TemplateValues: map[string]any{"title": "Hello from Go", "subtitle": "A signed image URL."},
			Format:         hcti.PNG,
		}, renderOptions)
		if err != nil {
			return fmt.Errorf("sign template URL: %w", err)
		}
		fmt.Println("Template URL:", templateURL)
	}
	return nil
}
