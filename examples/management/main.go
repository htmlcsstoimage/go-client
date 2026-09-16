// List proxy metadata without creating resources or rendering images.
package main

import (
	"context"
	"log"
	"time"

	hcti "github.com/htmlcsstoimage/go-client"
	"github.com/htmlcsstoimage/go-client/management"
)

func main() {
	client, err := management.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	options := management.ListOptions{Count: hcti.Ptr(100)}
	for {
		page, err := client.ListProxies(ctx, options)
		if err != nil {
			log.Fatal(err)
		}
		for _, proxy := range page.Data {
			log.Printf("proxy %s: enabled=%t", proxy.ID, proxy.Enabled)
		}
		if page.Pagination.NextPageStart == nil {
			break
		}
		options.PageStart = page.Pagination.NextPageStart
	}
}
