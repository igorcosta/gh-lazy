package cmd

import (
	"fmt"

	"github.com/igorcosta/gh-lazy/pkg/github"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh lazy",
	Short: "A GitHub CLI extension for managing projects, issues, and milestones",
	Long: `gh lazy is a GitHub CLI extension that helps you create project boards,
issues, milestones, and link them together efficiently.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tokenFile, _ := cmd.Flags().GetString("token-file")
		token, err := utils.GetToken(tokenFile)
		if err != nil {
			utils.PrintUserGuide()
			return fmt.Errorf("authentication error: %w", err)
		}

		client, err := github.NewClient(token)
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		username, err := client.GetUsername()
		if err != nil {
			return fmt.Errorf("failed to get GitHub username: %w", err)
		}

		utils.PrintWelcome(username)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("repo", "r", "", "The repository name (e.g., 'username/repo')")
	rootCmd.PersistentFlags().StringP("tasks", "t", "", "Path to the tasks JSON file")
	rootCmd.PersistentFlags().StringP("token-file", "f", "", "Path to the file containing the GitHub token")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}
