package main

import (
	"fmt"
	"os"

	tfu "github.com/dirien/tfu/cmd"
)

var (
	version string
	commit  string
)

func main() {
	if err := tfu.Execute(version, commit); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
