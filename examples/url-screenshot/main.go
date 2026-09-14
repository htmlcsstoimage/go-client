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

	image, err := client.CreateImage(ctx, hcti.URLImageRequest{
		URL:                 "https://example.com",
		FullScreen:          hcti.Ptr(true),
		BlockConsentBanners: hcti.Ptr(true),
		ImageOptions: hcti.ImageOptions{
			Format: hcti.WebP,
			RenderOptions: hcti.RenderOptions{
				ViewportWidth:  hcti.Ptr(1280),
				ViewportHeight: hcti.Ptr(720),
				DeviceScale:    hcti.Ptr(1.0),
				ColorScheme:    hcti.Ptr(hcti.Light),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create URL screenshot: %w", err)
	}
	fmt.Println(image.URL)
	return nil
}
