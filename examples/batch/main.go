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

	result, err := client.CreateImageBatch(ctx, hcti.BatchRequest{
		DefaultOptions: hcti.HTMLImageRequest{
			HTML:         `<div class="card">One design. Three colors.</div>`,
			GoogleFonts:  hcti.GoogleFonts{"Inter"},
			ImageOptions: hcti.ImageOptions{Format: hcti.PNG},
		},
		Variations: []hcti.ImageRequest{
			hcti.HTMLImageRequest{CSS: hcti.Ptr(`.card { width: 600px; padding: 48px; font: 32px 'Inter'; background: #dbeafe; color: #1e3a8a; }`)},
			hcti.HTMLImageRequest{CSS: hcti.Ptr(`.card { width: 600px; padding: 48px; font: 32px 'Inter'; background: #dcfce7; color: #14532d; }`)},
			hcti.HTMLImageRequest{CSS: hcti.Ptr(`.card { width: 600px; padding: 48px; font: 32px 'Inter'; background: #fce7f3; color: #831843; }`)},
		},
	})
	if err != nil {
		return fmt.Errorf("create batch: %w", err)
	}
	for i, image := range result.Images {
		fmt.Printf("Variation %d: %s (ID: %s)\n", i+1, image.URL, image.ID)
	}
	return nil
}
