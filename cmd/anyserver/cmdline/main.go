package main

import (
	anyserver "github.com/raibru/gsnet/cmd/anyserver/cmdline/cli"
	"github.com/raibru/gsnet/internal/sys"
)

func main() {
	defer sys.Exit(0)
	anyserver.Execute()
}
