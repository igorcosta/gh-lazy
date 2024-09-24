package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

func (c *Client) CreateProject(ctx context.Context, owner, title string) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", "project", "create", "--owner", owner, "--title", title, "--format", "json")
	output, err := cmd.Output()
	if err != nil {
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

func (c *Client) AddIssueToProject(ctx context.Context, owner, projectNumber string, issueURL string) error {
	cmd := exec.CommandContext(ctx, "gh", "project", "item-add", projectNumber,
		"--owner", owner, "--url", issueURL)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add issue to project: %w", err)
	}
	return nil
}

func (c *Client) LinkProjectToRepo(ctx context.Context, owner, repo, projectNumber string) error {
	cmd := exec.CommandContext(ctx, "gh", "project", "link", projectNumber,
		"--owner", owner, "--repo", fmt.Sprintf("%s/%s", owner, repo))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to link project to repository: %w", err)
	}
	return nil
}
