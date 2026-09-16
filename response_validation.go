package hcti

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/htmlcsstoimage/go-client/internal/httpapi"
)

// ResponseError indicates a malformed successful API response.
// Its message excludes response contents.
type ResponseError = httpapi.ResponseError

func (r *Image) validateResponse() string {
	if strings.TrimSpace(r.ID) == "" {
		return "missing image ID"
	}
	u, err := url.Parse(r.URL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "missing or invalid image URL"
	}
	return ""
}

func (r *BatchResult) validateResponse() string {
	if r.Images == nil {
		return "missing images array"
	}
	for i := range r.Images {
		if problem := r.Images[i].validateResponse(); problem != "" {
			return fmt.Sprintf("images[%d]: %s", i, problem)
		}
	}
	return ""
}

func (r *TemplateVersion) validateResponse() string {
	if strings.TrimSpace(r.TemplateID) == "" {
		return "missing template ID"
	}
	if r.TemplateVersion <= 0 {
		return "missing or invalid template version"
	}
	return ""
}

func (r *TemplatePage) validateResponse() string {
	if r.Data == nil {
		return "missing template data array"
	}
	for i, t := range r.Data {
		if strings.TrimSpace(t.ID) == "" || t.Version <= 0 {
			return fmt.Sprintf("data[%d]: missing template ID or version", i)
		}
	}
	return ""
}
