package main

import (
	"os"

	anyclient "github.com/raibru/gsnet/cmd/anyclient/cmdline/cli"
)

func main() {
	defer os.Exit(0)
	anyclient.Execute()
}
