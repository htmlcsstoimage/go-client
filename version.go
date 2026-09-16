package hcti

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionFile string

var sdkVersion = strings.TrimSpace(versionFile)

// Version returns the SDK version embedded at compile time from VERSION.
func Version() string { return sdkVersion }
