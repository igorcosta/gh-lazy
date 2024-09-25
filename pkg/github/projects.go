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

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse project items: %w", err)
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected JSON structure: 'items' is not an array")
	}

	var issues []models.IssueItem
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		content, ok := itemMap["content"].(map[string]interface{})
		if !ok {
			continue
		}

		if content["type"] != "Issue" {
			continue
		}

		number, _ := content["number"].(float64)
		title, _ := content["title"].(string)

		var repoFullName string
		if repo, ok := content["repository"].(map[string]interface{}); ok {
			repoName, _ := repo["name"].(string)
			if owner, ok := repo["owner"].(map[string]interface{}); ok {
				ownerLogin, _ := owner["login"].(string)
				repoFullName = fmt.Sprintf("%s/%s", ownerLogin, repoName)
			}
		}

		if repoFullName == "" || number == 0 {
			continue
		}

		issues = append(issues, models.IssueItem{
			Number:     int(number),
			Title:      title,
			Repository: repoFullName,
		})
	}

	return issues, nil
}

func (c *Client) GetIssueTitle(ctx context.Context, repo string, issueNumber int) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", "issue", "view", fmt.Sprintf("%d", issueNumber), "--repo", repo, "--json", "title", "--jq", ".title")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get title for issue #%d: %w", issueNumber, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (c *Client) DeleteProject(ctx context.Context, projectNumber string) error {
	owner, err := c.GetUsername()
	if err != nil {
		return fmt.Errorf("failed to get GitHub username: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gh", "project", "delete", projectNumber, "--owner", owner)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete project: %s - %w", string(output), err)
	}
	return nil
}

func (c *Client) LinkProjectToRepo(ctx context.Context, projectNumber, repoFullName string) error {
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s", repoFullName)
	}
	owner, repo := parts[0], parts[1]

	cmd := exec.CommandContext(ctx, "gh", "project", "link", projectNumber, "--owner", owner, "--repo", repo)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to link project to repository: %s - %w", string(output), err)
	}
	return nil
}
