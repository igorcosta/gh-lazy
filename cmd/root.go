package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh lazy",
	Short: "A GitHub CLI extension for managing projects, issues, and milestones",
	Long: `gh lazy is a GitHub CLI extension that helps you create project boards,
issues, milestones, and link them together efficiently and more cool lazy stuff.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("repo", "r", "", "The repository name (e.g., 'username/repo')")
	rootCmd.PersistentFlags().StringP("tasks", "t", "", "Path to the tasks JSON file")
	rootCmd.PersistentFlags().StringP("token-file", "f", ".token", "Path to the file containing the GitHub token")
}
