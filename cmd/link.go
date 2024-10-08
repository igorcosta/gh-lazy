package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/igorcosta/gh-lazy/pkg/config"
	"github.com/igorcosta/gh-lazy/pkg/github"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link project to repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return utils.WrapError(err, "failed to load config")
		}

		token, err := utils.ReadTokenFromFile(cfg.TokenFile)
		if err != nil {
			return utils.WrapError(err, "failed to read token")
		}

		client, err := github.NewClient(token)
		if err != nil {
			return utils.WrapError(err, "failed to create GitHub client")
		}

		projectNumber, err := cmd.Flags().GetString("project")
		if err != nil {
			return utils.WrapError(err, "failed to get project number")
		}

		repoName, err := cmd.Flags().GetString("repo")
		if err != nil {
			return utils.WrapError(err, "failed to get repository name")
		}

		if repoName == "" {
			return fmt.Errorf("repository name is required. Use -r or --repo flag to specify the name")
		}

		_, repo, err := splitRepoName(repoName)
		if err != nil {
			return utils.WrapError(err, "invalid repository name")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		err = client.LinkProjectToRepo(ctx, repo, projectNumber)
		if err != nil {
			return utils.WrapError(err, "failed to link project to repository")
		}

		utils.LogInfo(fmt.Sprintf("Successfully linked project %s to repository %s", projectNumber, repoName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().StringP("project", "p", "", "Project number to link")
	linkCmd.Flags().StringP("repo", "r", "", "The repository name (e.g., 'username/repo')")
	linkCmd.MarkFlagRequired("project")
	linkCmd.MarkFlagRequired("repo")
}
