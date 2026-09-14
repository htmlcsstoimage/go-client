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

	image, err := client.CreateImage(ctx, hcti.HTMLImageRequest{
		HTML:        `<article><span>HTML/CSS to Image</span><h1>Hello from Go.</h1><p>A little HTML. A little CSS. An image.</p></article>`,
		CSS:         hcti.Ptr(`article { width: 600px; padding: 48px; background: #101827; color: white; font-family: 'Inter'; } span { color: #67e8f9; } h1 { font-size: 48px; margin: 20px 0; } p { color: #cbd5e1; }`),
		GoogleFonts: hcti.GoogleFonts{"Inter"},
		ImageOptions: hcti.ImageOptions{
			Format:        hcti.PNG,
			RenderOptions: hcti.RenderOptions{DeviceScale: hcti.Ptr(2.0)},
		},
	})
	if err != nil {
		return fmt.Errorf("create HTML image: %w", err)
	}
	fmt.Println("Image ID:", image.ID)
	fmt.Println("Image URL:", image.URL)
	return nil
}
