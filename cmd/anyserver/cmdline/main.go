package main

import (
	"os"

	anyserver "github.com/raibru/gsnet/cmd/anyserver/cmdline/cli"
)

func main() {
	defer os.Exit(0)
	anyserver.Execute()
}
