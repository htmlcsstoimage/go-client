package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	request := hcti.TemplateRequest{
		Name:        "Go example social card",
		Description: "A reusable card with a title and subtitle.",
		HTML:        `<article><h1>{{title}}</h1><p>{{subtitle}}</p></article>`,
		CSS:         hcti.Ptr(`article { width: 600px; padding: 48px; background: #172554; color: white; font-family: 'Inter'; } h1 { font-size: 42px; } p { color: #bfdbfe; }`),
		GoogleFonts: hcti.GoogleFonts{"Inter"},
	}
	created, err := client.CreateTemplate(ctx, request)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	// Print the ID immediately so the template can be found if a later request fails.
	fmt.Printf("Created template %s, version %d\n", created.TemplateID, created.TemplateVersion)

	// Version updates send the complete desired configuration.
	request.CSS = hcti.Ptr(`article { width: 600px; padding: 48px; background: #312e81; color: white; font-family: 'Inter'; } h1 { font-size: 42px; } p { color: #c7d2fe; }`)
	updated, err := client.CreateTemplateVersion(ctx, created.TemplateID, request)
	if err != nil {
		return fmt.Errorf("create template version: %w", err)
	}
	fmt.Println("New version:", updated.TemplateVersion)

	image, err := client.CreateImage(ctx, hcti.TemplatedImageRequest{
		TemplateID:      updated.TemplateID,
		TemplateVersion: hcti.Ptr(updated.TemplateVersion), // Omit to use the latest version.
		TemplateValues:  map[string]any{"title": "A new release", "subtitle": "Built with Go and HTML/CSS to Image."},
		Format:          hcti.PNG,
	})
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	fmt.Println("Image URL:", image.URL)

	options := hcti.TemplateListOptions{Count: 10}
	for {
		page, err := client.ListTemplateVersions(ctx, created.TemplateID, options)
		if err != nil {
			return fmt.Errorf("list template versions: %w", err)
		}
		for _, template := range page.Data {
			fmt.Printf("Version %d: %s\n", template.Version, template.Name)
		}
		if page.Pagination.NextPageStart == nil {
			break
		}
		options.MaxVersion = page.Pagination.NextPageStart
	}
	// Keep the template and image so the printed URL remains useful.
	return nil
}
