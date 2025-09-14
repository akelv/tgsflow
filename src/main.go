package main

import (
	"os"

	cmdpkg "github.com/kelvin/tgsflow/src/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	code := cmdpkg.Execute(version, commit, date)
	os.Exit(code)
}
