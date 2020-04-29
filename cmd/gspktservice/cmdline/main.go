package main

import (
	gspktservice "github.com/raibru/gsnet/cmd/gspktservice/cmdline/cli"
	"github.com/raibru/gsnet/internal/sys"
)

func main() {
	defer sys.Exit(0)
	gspktservice.Execute()
}
