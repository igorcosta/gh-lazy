package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/igorcosta/gh-lazy/pkg/models"
)

func (c *Client) CreateProject(ctx context.Context, title string) (string, error) {
	owner, err := c.GetUsername()
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub username: %w", err)
	}

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

func (c *Client) AddIssueToProject(ctx context.Context, projectURL, issueURL string) error {
	owner, err := c.GetUsername()
	if err != nil {
		return fmt.Errorf("failed to get GitHub username: %w", err)
	}

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

func (c *Client) ListUserProjects(ctx context.Context) ([]models.Project, error) {
	owner, err := c.GetUsername()
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub username: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gh", "project", "list", "--owner", owner, "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %s - %w", string(output), err)
	}

	var result models.ProjectListResponse
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return result.Projects, nil
}

func (c *Client) ListProjectIssues(ctx context.Context, projectNumber string) ([]models.IssueItem, error) {
	owner, err := c.GetUsername()
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub username: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gh", "project", "item-list", projectNumber, "--owner", owner, "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list project items: %s - %w", string(output), err)
	}

	var result struct {
		Items []struct {
			Content struct {
				TypeName   string `json:"__typename"`
				Number     int    `json:"number"`
				Repository string `json:"repository"`
			} `json:"content"`
		} `json:"items"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse project items: %w", err)
	}

	var issues []models.IssueItem
	for _, item := range result.Items {
		if item.Content.TypeName == "Issue" {
			issues = append(issues, models.IssueItem{
				Number:     item.Content.Number,
				Repository: item.Content.Repository,
			})
		}
	}

	return issues, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectNumber string) error {
	owner, err := c.GetUsername()
	if err != nil {
		return fmt.Errorf("failed to get GitHub username: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gh", "project", "delete", projectNumber, "--owner", owner, "--yes")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete project: %s - %w", string(output), err)
	}
	return nil
}

func (c *Client) LinkProjectToRepo(ctx context.Context, repo, projectNumber string) error {
	owner, err := c.GetUsername()
	if err != nil {
		return fmt.Errorf("failed to get GitHub username: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gh", "project", "link", projectNumber,
		"--owner", owner, "--repo", fmt.Sprintf("%s/%s", owner, repo))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to link project to repository: %s - %w", string(output), err)
	}

	return nil
}
