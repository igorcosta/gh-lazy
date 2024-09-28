package main

import (
	"fmt"
	"os"

	"github.com/igorcosta/gh-lazy/cmd"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/igorcosta/gh-lazy/pkg/version"
)

func main() {
	// Handle version flag
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version.GetVersionInfo())
		return
	}

	// Print version information at startup
	fmt.Printf("gh-lazy %s\n", version.Version)

	if err := cmd.Execute(); err != nil {
		utils.LogError(err, "Command execution failed")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
