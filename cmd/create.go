package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/igorcosta/gh-lazy/pkg/config"
	"github.com/igorcosta/gh-lazy/pkg/github"
	"github.com/igorcosta/gh-lazy/pkg/models"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create project, milestones, and issues",
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

		tasksFile, err := utils.LoadTasksFile(cfg.TasksFile)
		if err != nil {
			return utils.WrapError(err, "failed to load tasks file")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		owner, repo, err := splitRepoName(cfg.Repo)
		if err != nil {
			return utils.WrapError(err, "invalid repository name")
		}

		projectURL, err := client.CreateProject(ctx, owner, tasksFile.ProjectTitle)
		if err != nil {
			return utils.WrapError(err, "failed to create project")
		}
		utils.LogInfo(fmt.Sprintf("Created new project: %s", projectURL))

		for _, milestone := range tasksFile.Milestones {
			milestoneNumber, err := createOrGetMilestone(ctx, client, owner, repo, milestone)
			if err != nil {
				return utils.WrapError(err, fmt.Sprintf("failed to create or get milestone %s", milestone.Title))
			}

			for _, issue := range milestone.Issues {
				issueNumber, err := createOrGetIssue(ctx, client, owner, repo, issue)
				if err != nil {
					return utils.WrapError(err, fmt.Sprintf("failed to create or get issue %s", issue.Title))
				}

				err = client.UpdateIssueMilestone(ctx, owner, repo, issueNumber, milestoneNumber)
				if err != nil {
					return utils.WrapError(err, fmt.Sprintf("failed to associate issue #%d with milestone #%d", issueNumber, milestoneNumber))
				}

				err = client.AddIssueToProject(ctx, owner, projectURL, fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, issueNumber))
				if err != nil {
					return utils.WrapError(err, fmt.Sprintf("failed to add issue #%d to project", issueNumber))
				}
			}
		}

		utils.LogInfo("All tasks completed successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createOrGetMilestone(ctx context.Context, client *github.Client, owner, repo string, milestone models.Milestone) (int, error) {
	existingMilestone, err := client.GetMilestoneByTitle(ctx, owner, repo, milestone.Title)
	if err != nil {
		return 0, utils.WrapError(err, "checking existing milestone")
	}
	if existingMilestone != nil {
		return existingMilestone.Number, nil
	}

	number, err := client.CreateMilestone(ctx, owner, repo, milestone)
	if err != nil {
		return 0, utils.WrapError(err, "creating milestone")
	}
	return number, nil
}

func createOrGetIssue(ctx context.Context, client *github.Client, owner, repo string, issue models.Issue) (int, error) {
	existingIssue, err := client.GetIssueByTitle(ctx, owner, repo, issue.Title)
	if err != nil {
		return 0, utils.WrapError(err, "checking existing issue")
	}
	if existingIssue != nil {
		return existingIssue.Number, nil
	}

	number, err := client.CreateIssue(ctx, owner, repo, issue)
	if err != nil {
		return 0, utils.WrapError(err, "creating issue")
	}
	return number, nil
}

func splitRepoName(repoName string) (string, string, error) {
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository name format. Expected 'owner/repo', got '%s'", repoName)
	}
	return parts[0], parts[1], nil
}
