package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/igorcosta/gh-lazy/pkg/config"
	"github.com/igorcosta/gh-lazy/pkg/github"
	"github.com/igorcosta/gh-lazy/pkg/models"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create project, milestones, and issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoName, err := cmd.Flags().GetString("repo")
		if err != nil {
			return fmt.Errorf("failed to get repository name: %w", err)
		}

		if repoName == "" {
			return fmt.Errorf("repository name is required. Use -r or --repo flag to specify the name")
		}

		tasksFile, err := cmd.Flags().GetString("tasks")
		if err != nil {
			return fmt.Errorf("failed to get tasks file path: %w", err)
		}

		if tasksFile == "" {
			return fmt.Errorf("tasks file path is required. Use -t or --tasks flag to specify the path")
		}

		if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
			return fmt.Errorf("tasks file does not exist: %s", tasksFile)
		}

		absTasksFile, err := filepath.Abs(tasksFile)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of tasks file: %w", err)
		}

		token, err := utils.GetToken(cfg.TokenFile)
		if err != nil {
			utils.PrintUserGuide()
			return fmt.Errorf("authentication error: %w", err)
		}

		client, err := github.NewClient(token)
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		tasks, err := utils.LoadTasksFile(absTasksFile)
		if err != nil {
			return fmt.Errorf("failed to load tasks file: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		owner, repo, err := splitRepoName(repoName)
		if err != nil {
			return fmt.Errorf("invalid repository name: %w", err)
		}

		totalTasks := len(tasks.Milestones) + 2 // +2 for project creation and linking
		for _, m := range tasks.Milestones {
			totalTasks += len(m.Issues)
		}

		bar := progressbar.NewOptions(totalTasks,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription("[cyan][1/3][reset] Creating project, milestones, and issues..."),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}))

		completed := 0
		skipped := 0
		failed := 0
		createdIssues := []string{}

		projectURL, err := client.CreateProject(ctx, tasks.ProjectTitle)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		bar.Add(1)
		completed++

		// Extract project number from URL
		parts := strings.Split(projectURL, "/")
		projectNumber := parts[len(parts)-1]

		// Link the project to the repository
		err = client.LinkProjectToRepo(ctx, projectNumber, repoName)
		if err != nil {
			color.Yellow("⚠️ Failed to link project to repository: %v", err)
			skipped++
		} else {
			color.Green("✅ Project linked to repository %s", repoName)
			completed++
		}
		bar.Add(1)

		for _, milestone := range tasks.Milestones {
			milestoneNumber, err := createOrGetMilestone(ctx, client, owner, repo, milestone)
			if err != nil {
				color.Red("❌ Failed to create/get milestone %s: %v", milestone.Title, err)
				failed++
				bar.Add(1)
				continue
			}
			bar.Add(1)
			completed++

			for _, issue := range milestone.Issues {
				issueNumber, err := createOrGetIssue(ctx, client, owner, repo, issue)
				if err != nil {
					color.Red("❌ Failed to create/get issue %s: %v", issue.Title, err)
					failed++
					bar.Add(1)
					continue
				}

				err = client.UpdateIssueMilestone(ctx, owner, repo, issueNumber, milestoneNumber)
				if err != nil {
					color.Yellow("⚠️ Failed to associate issue #%d with milestone #%d: %v", issueNumber, milestoneNumber, err)
					skipped++
				}

				issueURL := fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, issueNumber)
				err = client.AddIssueToProject(ctx, projectURL, issueURL)
				if err != nil {
					color.Yellow("⚠️ Failed to add issue #%d to project: %v", issueNumber, err)
					skipped++
				}

				createdIssues = append(createdIssues, issueURL)
				bar.Add(1)
				completed++
			}
		}

		bar.Finish()
		fmt.Println()

		color.Green("✅ Project created successfully: %s", projectURL)
		fmt.Println("Created issues:")
		for _, issueURL := range createdIssues {
			color.Cyan("  • %s", issueURL)
		}

		fmt.Println()
		color.Green("📊 Summary:")
		color.Green("  ✅ Completed tasks: %d", completed)
		color.Yellow("  ⚠️ Skipped tasks: %d", skipped)
		color.Red("  ❌ Failed tasks: %d", failed)
		color.Cyan("  🔗 Project URL: %s", projectURL)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("repo", "r", "", "The repository name (e.g., 'username/repo')")
	createCmd.Flags().StringP("tasks", "t", "", "Path to the tasks JSON file")
	createCmd.MarkFlagRequired("repo")
	createCmd.MarkFlagRequired("tasks")
}

func createOrGetMilestone(ctx context.Context, client *github.Client, owner, repo string, milestoneWithIssues models.MilestoneWithIssues) (int, error) {
	existingMilestone, err := client.GetMilestoneByTitle(ctx, owner, repo, milestoneWithIssues.Title)
	if err != nil {
		return 0, fmt.Errorf("checking existing milestone: %w", err)
	}
	if existingMilestone != nil {
		return existingMilestone.Number, nil
	}

	number, err := client.CreateMilestone(ctx, owner, repo, milestoneWithIssues.Milestone)
	if err != nil {
		return 0, fmt.Errorf("creating milestone: %w", err)
	}
	return number, nil
}

func createOrGetIssue(ctx context.Context, client *github.Client, owner, repo string, issue models.Issue) (int, error) {
	existingIssue, err := client.GetIssueByTitle(ctx, owner, repo, issue.Title)
	if err != nil {
		return 0, fmt.Errorf("checking existing issue: %w", err)
	}
	if existingIssue != nil {
		return existingIssue.Number, nil
	}

	number, err := client.CreateIssue(ctx, owner, repo, issue)
	if err != nil {
		return 0, fmt.Errorf("creating issue: %w", err)
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
