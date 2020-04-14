package cli

import (
	"fmt"
	"io"
)

// var section of global version variables
var (
	Version     = "0.1.0-snapshot"
	Build       = "-"
	serviceName = "anyserver"
	author      = "rbr <github.com/raibru>"
	license     = "MIT License (c) 2020 raibru"
)

// PrintVersion prints the tool versions string
func PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "%s - v%s (build:%s)\n", serviceName, Version, Build)
	fmt.Fprintf(w, "  author : %s\n", author)
	fmt.Fprintf(w, "  license: %s\n\n", license)
}
