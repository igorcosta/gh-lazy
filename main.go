package main

import (
	"fmt"
	"os"

	"github.com/igorcosta/gh-lazy/cmd"
	"github.com/igorcosta/gh-lazy/pkg/utils"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.LogError(err, "Command execution failed")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Version: %s, Commit: %s\n", version, commit)
		return
	}
}
