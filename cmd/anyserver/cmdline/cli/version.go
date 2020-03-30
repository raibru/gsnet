package cli

import (
	"fmt"
	"io"
)

// PrintVersion prints the tool versions string
func PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "anyserver - v0.1.0\n")
	fmt.Fprintf(w, "  author : rbr <raibru@web.de>\n")
	fmt.Fprintf(w, "  license: MIT\n\n")
}
