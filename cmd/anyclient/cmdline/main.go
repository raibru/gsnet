package main

import (
	anyclient "github.com/raibru/gsnet/cmd/anyclient/cmdline/cli"
	"github.com/raibru/gsnet/internal/sys"
)

func main() {
	defer sys.Exit(0)
	anyclient.Execute()
}
