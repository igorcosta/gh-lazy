package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, out *os.File) error {
	if err := checkGHCLIInstalled(); err != nil {
		return fmt.Errorf("GitHub CLI is required: %w", err)
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Usage = func() {
		printWelcome("GitHub User")
		printHelp()
	}
	repoPtr := flags.String("reponame", "", "The repository name (e.g., 'username/repo')")
	tasksFilePtr := flags.String("tasks", "", "Path to the tasks JSON file")
	tokenFilePtr := flags.String("tokenfile", ".token", "Path to the file containing the GitHub token")

	if len(args) == 1 || (len(args) == 2 && (args[1] == "-h" || args[1] == "help")) {
		flags.Usage()
		return nil
	}

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if *repoPtr == "" || *tasksFilePtr == "" {
		return fmt.Errorf("error: -reponame and -tasks flags are required")
	}

	token, err := getToken(*tokenFilePtr)
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}

	client, err := NewGitHubClient(token)
	if err != nil {
		return fmt.Errorf("error creating GitHub client: %w", err)
	}

	tasksFile, err := loadTasksFile(*tasksFilePtr)
	if err != nil {
		return fmt.Errorf("error loading tasks file: %w", err)
	}

	tool, err := NewLazyTool(client, *repoPtr, tasksFile)
	if err != nil {
		return fmt.Errorf("error creating LazyTool: %w", err)
	}

	if err := tool.Run(); err != nil {
		return fmt.Errorf("error running tool: %w", err)
	}

	fmt.Fprintln(out, "\nAll tasks completed successfully!")
	return nil
}
