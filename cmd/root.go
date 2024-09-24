package cmd

import (
	"os"

	"github.com/igorcosta/gh-lazy/pkg/github"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh lazy",
	Short: "A GitHub CLI extension for managing projects, issues, and milestones",
	Long: `gh lazy is a GitHub CLI extension that helps you create project boards,
issues, milestones, and link them together efficiently.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Get the GitHub username
			token, err := utils.ReadTokenFromFile(".token") // You might want to make this configurable
			if err != nil {
				utils.LogError(err, "Failed to read GitHub token")
				os.Exit(1)
			}

			client, err := github.NewClient(token)
			if err != nil {
				utils.LogError(err, "Failed to create GitHub client")
				os.Exit(1)
			}

			username, err := client.GetUsername()
			if err != nil {
				utils.LogError(err, "Failed to get GitHub username")
				os.Exit(1)
			}

			utils.PrintWelcome(username)
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("repo", "r", "", "The repository name (e.g., 'username/repo')")
	rootCmd.PersistentFlags().StringP("tasks", "t", "", "Path to the tasks JSON file")
	rootCmd.PersistentFlags().StringP("token-file", "f", ".token", "Path to the file containing the GitHub token")

	// Disable the completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Remove the 'help' subcommand
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}
