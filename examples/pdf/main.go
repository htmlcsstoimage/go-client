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

	margin := hcti.PDFLength{Value: 12, Unit: hcti.Millimeters}
	image, err := client.CreateImage(ctx, hcti.HTMLImageRequest{
		HTML:        `<main><h1>Monthly report</h1><p>Revenue increased by 12% this month.</p><table><tr><th>Month</th><th>Revenue</th></tr><tr><td>September</td><td>$12,400</td></tr></table></main>`,
		CSS:         hcti.Ptr(`body { font-family: 'Roboto'; color: #172554; } table { border-collapse: collapse; width: 100%; } th, td { padding: 12px; border: 1px solid #cbd5e1; text-align: left; } th { background: #eff6ff; }`),
		GoogleFonts: hcti.GoogleFonts{"Roboto"},
		ImageOptions: hcti.ImageOptions{
			Format:        hcti.PDF,
			RenderOptions: hcti.RenderOptions{MediaType: hcti.Ptr(hcti.Print)},
			PDFOptions: &hcti.PDFOptions{
				PrintBackground: hcti.Ptr(true),
				PageWidth:       &hcti.PDFLength{Value: 210, Unit: hcti.Millimeters},
				PageHeight:      &hcti.PDFLength{Value: 297, Unit: hcti.Millimeters},
				Margins:         &hcti.PDFMargins{Top: margin, Right: margin, Bottom: margin, Left: margin},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create PDF: %w", err)
	}
	fmt.Println("PDF URL:", image.URL)
	return nil
}
