package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/igorcosta/gh-lazy/pkg/config"
	"github.com/igorcosta/gh-lazy/pkg/github"
	"github.com/igorcosta/gh-lazy/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func getCurrentRepo() (string, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current repository: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Delete a GitHub project and optionally all linked issues",
	Long: `Delete a GitHub project and optionally all issues linked to it.

**Warning:** This operation is irreversible. Use with caution.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		projectIDOrURL, _ := cmd.Flags().GetString("projectid")
		deleteAll, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		repoName, err := getCurrentRepo()
		if err != nil {
			return fmt.Errorf("failed to get current repository: %w", err)
		}
		fmt.Printf("Current repository: %s\n", repoName)

		token, err := utils.GetToken(cfg.TokenFile)
		if err != nil {
			utils.PrintUserGuide()
			return fmt.Errorf("authentication error: %w", err)
		}

		client, err := github.NewClient(token)
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if projectIDOrURL == "" {
			projects, err := client.ListUserProjects(ctx)
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}

			if len(projects) == 0 {
				fmt.Println("No projects found.")
				return nil
			}

			projectNames := make([]string, len(projects))
			for i, p := range projects {
				projectNames[i] = fmt.Sprintf("%s (ID: %d)", p.Title, p.Number)
			}

			prompt := promptui.Select{
				Label: "Select a project to nuke",
				Items: projectNames,
			}

			index, _, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			selectedProject := projects[index]
			projectIDOrURL = fmt.Sprintf("%d", selectedProject.Number)
			fmt.Printf("Selected project: %s\n", selectedProject.Title)

			if !dryRun {
				confirmPrompt := promptui.Prompt{
					Label:     fmt.Sprintf("Are you sure you want to delete project '%s' and all linked issues?", selectedProject.Title),
					IsConfirm: true,
				}
				result, err := confirmPrompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			if !cmd.Flags().Changed("all") {
				deleteAllPrompt := promptui.Prompt{
					Label:     "Do you want to delete all issues associated with the project?",
					IsConfirm: true,
				}
				result, err := deleteAllPrompt.Run()
				if err != nil {
					deleteAll = false
				} else {
					deleteAll = strings.ToLower(result) == "y"
				}
			}
		}

		projectNumber, err := utils.ParseProjectID(projectIDOrURL)
		if err != nil {
			return fmt.Errorf("failed to parse project ID: %w", err)
		}

		issues, err := client.ListProjectIssues(ctx, projectNumber)
		if err != nil {
			return fmt.Errorf("failed to list issues linked to the project: %w", err)
		}

		totalTasks := 1 // For deleting the project
		if deleteAll {
			totalTasks += len(issues)
		}

		bar := progressbar.NewOptions(totalTasks,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription("[cyan][1/2][reset] Processing..."),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}))

		if dryRun {
			color.Yellow("** Dry Run Mode Enabled **")
			color.Yellow("No actual deletions will occur.")
			fmt.Println()
		}

		failed := 0
		deleted := 0
		skipped := 0

		if deleteAll {
			color.Cyan("Deleting issues associated with the project:")
			for _, issue := range issues {
				if dryRun {
					color.Cyan("🗒️ Would delete issue #%d: %s (Repository: %s)", issue.Number, issue.Title, issue.Repository)
					bar.Add(1)
				} else {
					fmt.Printf("Deleting issue #%d: %s (Repository: %s)\n", issue.Number, issue.Title, issue.Repository)
					err := client.DeleteIssue(ctx, issue.Repository, issue.Number)
					if err != nil {
						color.Red("❌ Failed to delete issue #%d: %v", issue.Number, err)
						failed++
					} else {
						color.Green("🗑️ Deleted issue #%d: %s", issue.Number, issue.Title)
						deleted++
					}
					bar.Add(1)
				}
				time.Sleep(time.Second) // Small delay to avoid overwhelming the API
			}
		} else {
			skipped = len(issues)
			if skipped > 0 {
				color.Yellow("Skipping deletion of %d issues", skipped)
			}
		}

		if dryRun {
			color.Cyan("🗒️ Would delete project %s", projectNumber)
			bar.Add(1)
		} else {
			fmt.Printf("Deleting project %s\n", projectNumber)
			err = client.DeleteProject(ctx, projectNumber)
			if err != nil {
				color.Red("❌ Failed to delete project: %v", err)
				failed++
			} else {
				bar.Add(1)
				color.Green("✅ Project deleted successfully")
			}
		}

		bar.Finish()
		fmt.Println()

		color.Green("📊 Summary:")
		if deleteAll {
			if dryRun {
				color.Green("  🗒️ Issues that would be deleted: %d", len(issues))
			} else {
				color.Green("  🗑️ Deleted issues: %d", deleted)
				if failed > 0 {
					color.Red("  ❌ Failed deletions: %d", failed)
				}
			}
		} else {
			color.Yellow("  ⏭️ Skipped issues: %d", skipped)
		}
		if dryRun {
			color.Green("  🗒️ Project that would be deleted: %s", projectNumber)
		} else {
			color.Green("  🗑️ Deleted project: %s", projectNumber)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(nukeCmd)
	nukeCmd.Flags().StringP("projectid", "p", "", "Project ID or URL to nuke")
	nukeCmd.Flags().BoolP("all", "a", false, "Delete all issues linked to the project")
	nukeCmd.Flags().Bool("dry-run", false, "Show what would happen without making changes")
}
