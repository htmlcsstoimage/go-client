package main

import (
	"fmt"
	"github.com/htmlcsstoimage/go-client"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: go run ./examples/image-url IMAGE_ID")
	}
	// Existing image retrieval URLs need no API credentials.
	client := hcti.NewClient("", "")
	imageURL, err := client.ImageURL(os.Args[1], hcti.RenderImageOptions{
		Format: hcti.WebP, Width: hcti.Ptr(600),
		Crop: &hcti.Crop{
			Horizontal:      &hcti.CropSpan{Size: &hcti.CropValue{Value: 100, Unit: hcti.CropPercent}},
			AspectRatio:     &hcti.AspectRatio{Width: 16, Height: 9},
			AspectRatioAxis: hcti.CropWidth, ComputedOrigin: hcti.CropCenter,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(imageURL)
}
