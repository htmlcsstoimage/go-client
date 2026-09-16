package management

import (
	"fmt"
	"net/url"
	"strconv"
)

// ListOptions selects one page; pagination is never fetched automatically.
type ListOptions struct {
	// Count is the page size, 1–100. Nil lets the API choose its default.
	Count *int
	// PageStart is the opaque cursor returned by the preceding page.
	PageStart *string
}

func (o ListOptions) query() (url.Values, error) {
	q := url.Values{}
	if o.Count != nil {
		if *o.Count < 1 || *o.Count > 100 {
			return nil, fmt.Errorf("hcti: count must be between 1 and 100")
		}
		q.Set("count", strconv.Itoa(*o.Count))
	}
	if o.PageStart != nil {
		q.Set("page_start", *o.PageStart)
	}
	return q, nil
}

// Pagination supplies the next page's cursor; nil means there are no more pages.
type Pagination struct {
	NextPageStart *string `json:"next_page_start"`
}

// Page contains one page of resources and its continuation cursor.
type Page[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func (p *Page[T]) validateResponse() string {
	if p.Data == nil {
		return "missing page data"
	}
	for i := range p.Data {
		if v, ok := any(&p.Data[i]).(interface{ validateResponse() string }); ok {
			if problem := v.validateResponse(); problem != "" {
				return problem
			}
		}
	}
	return ""
}
