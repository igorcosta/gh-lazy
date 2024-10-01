package main

import (
	"fmt"
	"os"

	"github.com/igorcosta/gh-lazy/cmd"
	"github.com/igorcosta/gh-lazy/pkg/version"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("gh-lazy version %s\n", version.Version)
			fmt.Printf("commit: %s\n", version.Commit)
			fmt.Printf("built at: %s\n", version.BuildDate)
			return
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
