package main

import (
	"os"

	gspktservice "github.com/raibru/gsnet/cmd/gspktservice/cmdline/cli"
)

func main() {
	defer os.Exit(0)
	gspktservice.Execute()
}
