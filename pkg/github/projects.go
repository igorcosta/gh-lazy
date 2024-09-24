package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/igorcosta/gh-lazy/pkg/models"
)

func (c *Client) CreateProject(ctx context.Context, owner, title string) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", "project", "create", "--owner", owner, "--title", title, "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("failed to create project: %s", string(exitError.Stderr))
		}
		return "", fmt.Errorf("failed to create project: %w", err)
	}
	var response struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return "", fmt.Errorf("failed to parse project creation output: %w", err)
	}
	return response.URL, nil
}

func (c *Client) AddIssueToProject(ctx context.Context, owner, projectURL, issueURL string) error {
	// Extract project number from URL
	parts := strings.Split(projectURL, "/")
	projectNumber := parts[len(parts)-1]

	cmd := exec.CommandContext(ctx, "gh", "project", "item-add", projectNumber,
		"--owner", owner, "--url", issueURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add issue to project: %s - %w", string(output), err)
	}

	return nil
}

// List issues linked to the project
func (c *Client) ListProjectIssues(ctx context.Context, owner, projectNumber string) ([]models.IssueItem, error) {
	cmd := exec.CommandContext(ctx, "gh", "project", "item-list", projectNumber, "--owner", owner, "--json", "content")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list project items: %w", err)
	}

	var items []struct {
		Content struct {
			TypeName   string `json:"__typename"`
			Number     int    `json:"number"`
			Repository struct {
				NameWithOwner string `json:"nameWithOwner"`
			} `json:"repository"`
		} `json:"content"`
	}

	if err := json.Unmarshal(output, &items); err != nil {
		return nil, fmt.Errorf("failed to parse project items: %w", err)
	}

	var issues []models.IssueItem
	for _, item := range items {
		if item.Content.TypeName == "Issue" {
			issues = append(issues, models.IssueItem{
				Number:     item.Content.Number,
				Repository: item.Content.Repository,
			})
		}
	}

	return issues, nil
}

// Delete the project
func (c *Client) DeleteProject(ctx context.Context, owner, projectNumber string) error {
	cmd := exec.CommandContext(ctx, "gh", "project", "delete", projectNumber, "--confirm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete project: %s", string(output))
	}
	return nil
}

func (c *Client) LinkProjectToRepo(ctx context.Context, owner, repo, projectNumber string) error {
	cmd := exec.CommandContext(ctx, "gh", "project", "link", projectNumber,
		"--owner", owner, "--repo", fmt.Sprintf("%s/%s", owner, repo))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to link project to repository: %s - %w", string(output), err)
	}

	return nil
}
