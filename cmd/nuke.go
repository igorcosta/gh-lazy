package cmd

import (
	"context"
	"fmt"
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

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Delete a GitHub project and optionally all linked issues",
	Long: `Delete a GitHub project and optionally all issues linked to it.

**Warning:** This operation is irreversible. Use with caution.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Parse flags
		projectIDOrURL, _ := cmd.Flags().GetString("projectid")
		deleteAll, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Create GitHub client using the token from config
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

		// If no project ID is provided, list projects and allow user to select
		if projectIDOrURL == "" {
			projects, err := client.ListUserProjects(ctx)
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}

			if len(projects) == 0 {
				fmt.Println("No projects found.")
				return nil
			}

			// Prompt user to select a project
			projectNames := []string{}
			for _, p := range projects {
				projectNames = append(projectNames, fmt.Sprintf("%s (ID: %d)", p.Title, p.Number))
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

			// Confirm deletion if not in dry-run mode
			if !cmd.Flags().Changed("dry-run") && !dryRun {
				confirmPrompt := promptui.Prompt{
					Label:     fmt.Sprintf("Are you sure you want to delete project '%s' and all linked issues?", selectedProject.Title),
					IsConfirm: true,
					Default:   "n",
				}
				result, err := confirmPrompt.Run()
				if err != nil || strings.ToLower(result) != "y" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			// Ask if the user wants to delete all associated issues
			if !cmd.Flags().Changed("all") {
				if !dryRun {
					deleteAllPrompt := promptui.Prompt{
						Label:     "Do you want to delete all issues associated with the project? (y/N)",
						IsConfirm: true,
						Default:   "n",
					}
					result, err := deleteAllPrompt.Run()
					if err != nil {
						deleteAll = false
					} else {
						deleteAll = (strings.ToLower(result) == "y")
					}
				}
			}
		}

		// Extract project number
		projectNumber, err := utils.ParseProjectID(projectIDOrURL)
		if err != nil {
			return fmt.Errorf("failed to parse project ID: %w", err)
		}

		// Fetch issues linked to the project
		issues, err := client.ListProjectIssues(ctx, projectNumber)
		if err != nil {
			return fmt.Errorf("failed to list issues linked to the project: %w", err)
		}

		totalTasks := 1 // For deleting the project
		if deleteAll {
			totalTasks += len(issues)
		}

		// Set up the progress bar
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

		// Delete issues if --all is set
		if deleteAll {
			for _, issue := range issues {
				if dryRun {
					color.Cyan("🗒️ Would delete issue #%d: %s", issue.Number, issue.Title)
				} else {
					err := client.DeleteIssue(ctx, issue.Repository, issue.Number)
					if err != nil {
						color.Red("❌ Failed to delete issue #%d: %v", issue.Number, err)
						failed++
					} else {
						color.Green("🗑️ Deleted issue #%d: %s", issue.Number, issue.Title)
						deleted++
					}
				}
				bar.Add(1)
			}
		} else if dryRun && len(issues) > 0 {
			color.Cyan("Issues linked to the project that would not be deleted:")
			for _, issue := range issues {
				color.Cyan("  - Issue #%d: %s", issue.Number, issue.Title)
			}
		}

		// Delete the project
		if dryRun {
			color.Cyan("🗒️ Would delete project %s", projectNumber)
			bar.Add(1)
		} else {
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

		// Summary
		color.Green("📊 Summary:")
		if deleteAll {
			if dryRun {
				color.Green("  🗒️ Issues that would be deleted: %d", len(issues))
			} else {
				color.Green("  🗑️ Deleted issues: %d", deleted)
			}
		}
		if dryRun {
			color.Green("  🗒️ Project that would be deleted: %s", projectNumber)
		} else {
			color.Green("  🗑️ Deleted project: %s", projectNumber)
		}
		if failed > 0 {
			color.Red("  ❌ Failed deletions: %d", failed)
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
