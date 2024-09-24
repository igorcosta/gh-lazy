package main

import (
	"fmt"
	"os"

	"yourmodule/cmd"
	"yourmodule/pkg/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.LogError(err, "Command execution failed")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
