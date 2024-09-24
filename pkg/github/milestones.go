package github

import (
	"context"
	"fmt"

	"github.com/igorcosta/gh-lazy/pkg/models"
)

func (c *Client) CreateMilestone(ctx context.Context, owner, repo string, milestone models.Milestone) (int, error) {
	url := fmt.Sprintf("repos/%s/%s/milestones", owner, repo)
	var response struct {
		Number int `json:"number"`
	}
	if err := c.Post(ctx, url, milestone, &response); err != nil {
		return 0, fmt.Errorf("failed to create milestone: %w", err)
	}
	return response.Number, nil
}

func (c *Client) GetMilestoneByTitle(ctx context.Context, owner, repo, title string) (*models.Milestone, error) {
	url := fmt.Sprintf("repos/%s/%s/milestones", owner, repo)
	var milestones []models.Milestone
	if err := c.Get(ctx, url, &milestones); err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}
	for _, m := range milestones {
		if m.Title == title {
			return &m, nil
		}
	}
	return nil, nil
}
