package main

import (
	pktservice "github.com/raibru/gsnet/cmd/pktservice/cmdline/cli"
	"github.com/raibru/gsnet/internal/sys"
)

func main() {
	defer sys.Exit(0)
	pktservice.Execute()
}
