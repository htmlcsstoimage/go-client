package hcti

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionFile string

var userAgent = "HCTIGo/" + strings.TrimSpace(versionFile)
