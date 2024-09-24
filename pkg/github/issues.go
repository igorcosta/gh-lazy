package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/igorcosta/gh-lazy/pkg/models"
)

func (c *Client) CreateIssue(ctx context.Context, owner, repo string, issue models.Issue) (int, error) {
	url := fmt.Sprintf("repos/%s/%s/issues", owner, repo)
	var response struct {
		Number int `json:"number"`
	}

	payload, err := json.Marshal(issue)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal issue: %w", err)
	}

	if err := c.Post(ctx, url, bytes.NewReader(payload), &response); err != nil {
		return 0, fmt.Errorf("failed to create issue: %w", err)
	}
	return response.Number, nil
}

func (c *Client) GetIssueByTitle(ctx context.Context, owner, repo, title string) (*models.Issue, error) {
	url := fmt.Sprintf("repos/%s/%s/issues?state=all", owner, repo)
	var issues []models.Issue
	if err := c.Get(ctx, url, &issues); err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}
	for _, i := range issues {
		if i.Title == title {
			return &i, nil
		}
	}
	return nil, nil
}

func (c *Client) UpdateIssueMilestone(ctx context.Context, owner, repo string, issueNumber, milestoneNumber int) error {
	url := fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, issueNumber)
	payload := map[string]interface{}{
		"milestone": milestoneNumber,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var response interface{}
	if err := c.Patch(ctx, url, bytes.NewReader(jsonPayload), &response); err != nil {
		return fmt.Errorf("failed to update issue milestone: %w", err)
	}
	return nil
}
