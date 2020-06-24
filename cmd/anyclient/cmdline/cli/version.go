package cli

import (
	"fmt"
	"io"
)

// var section
var (
	major     = "0"
	minor     = "1"
	patch     = "0"
	buildTag  = "-"
	buildDate = "-"
	appName   = "anyclient"
	author    = "raibru <github.com/raibru>"
	license   = "MIT License (c) 2020 raibru"
)

// PrintVersion prints the tool versions string
func PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "%s - v%s.%s.%s\n", appName, major, minor, patch)
	fmt.Fprintf(w, "  Build-%s (%s)\n", buildTag, buildDate)
	fmt.Fprintf(w, "  author : %s\n", author)
	fmt.Fprintf(w, "  license: %s\n\n", license)
}

// VersionShort returns short version info
func VersionShort() string {
	s := ""
	s += fmt.Sprintf("%s - v%s.%s.%s\n", appName, major, minor, patch)
	s += fmt.Sprintf("  Build-%s (%s)\n", buildTag, buildDate)
	return s
}
