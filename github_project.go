package lazytool

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync/atomic"
)

func (lt *LazyTool) addIssueToProjectUsingGHCLI(ctx context.Context, issueNumber int, projectNumber string) error {
	cmd := exec.CommandContext(ctx, "gh", "project", "item-add", projectNumber,
		"--owner", lt.owner,
		"--url", fmt.Sprintf("https://github.com/%s/%s/issues/%d", lt.owner, lt.repo, issueNumber))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add issue to project using gh CLI: %w\nStderr: %s", err, stderr.String())
	}
	return nil
}

func (lt *LazyTool) linkProjectToRepo(ctx context.Context, projectNumber string) error {
	lt.progressChan <- "Linking project to repository"
	cmd := exec.CommandContext(ctx, "gh", "project", "link", projectNumber,
		"--owner", lt.owner,
		"--repo", fmt.Sprintf("%s/%s", lt.owner, lt.repo))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to link project to repository using gh CLI: %w\nStderr: %s", err, stderr.String())
	}
	atomic.AddInt32(&lt.tasksCompleted, 1)
	return nil
}
