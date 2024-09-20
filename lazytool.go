package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LazyTool struct {
	client         GitHubClient
	owner          string
	repo           string
	tasksFile      *TasksFile
	progressChan   chan string
	logger         *log.Logger
	mu             sync.Mutex
	tasksTotal     int32
	tasksCompleted int32
	tasksSkipped   int32
}

func NewLazyTool(client GitHubClient, repoFullName string, tasksFile *TasksFile) (*LazyTool, error) {
	parts := strings.Split(repoFullName, "/")
	var owner, repo string
	if len(parts) == 2 {
		owner, repo = parts[0], parts[1]
	} else if len(parts) == 1 {
		repo = parts[0]
		fmt.Print("Enter the owner for the repository: ")
		fmt.Scanln(&owner)
	} else {
		return nil, fmt.Errorf("invalid repository name format. Expected 'owner/repo' or 'repo', got '%s'", repoFullName)
	}

	return &LazyTool{
		client:       client,
		owner:        owner,
		repo:         repo,
		tasksFile:    tasksFile,
		progressChan: make(chan string),
		logger:       log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (lt *LazyTool) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	go lt.showProgress()
	defer close(lt.progressChan)

	username, err := lt.getUsername(ctx)
	if err != nil {
		return fmt.Errorf("error getting username: %w", err)
	}

	printWelcome(username)

	projectURL, err := lt.createProjectUsingGHCLI(ctx, lt.owner)
	if err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}
	fmt.Printf("Created new project: %s\n", projectURL)

	projectNumber, err := lt.getProjectNumberFromURL(projectURL)
	if err != nil {
		return fmt.Errorf("error getting project number: %w", err)
	}

	lt.calculateTotalTasks()

	if err := lt.createAllMilestones(ctx); err != nil {
		return fmt.Errorf("error creating milestones: %w", err)
	}

	if err := lt.createAllIssues(ctx); err != nil {
		return fmt.Errorf("error creating issues: %w", err)
	}

	if err := lt.associateIssuesWithMilestones(ctx); err != nil {
		return fmt.Errorf("error associating issues with milestones: %w", err)
	}

	if err := lt.addAllIssuesToProject(ctx, projectNumber); err != nil {
		return fmt.Errorf("error adding issues to project: %w", err)
	}

	if err := lt.linkProjectToRepo(ctx, projectNumber); err != nil {
		return fmt.Errorf("error linking project to repository: %w", err)
	}

	return nil
}

func (lt *LazyTool) calculateTotalTasks() {
	var total int32
	for _, milestone := range lt.tasksFile.Milestones {
		total++
		total += int32(len(milestone.Issues))
	}
	atomic.StoreInt32(&lt.tasksTotal, total)
}

func (lt *LazyTool) showProgress() {
	for message := range lt.progressChan {
		completed := atomic.LoadInt32(&lt.tasksCompleted)
		skipped := atomic.LoadInt32(&lt.tasksSkipped)
		total := atomic.LoadInt32(&lt.tasksTotal)
		fmt.Printf("\r[%d/%d] Completed: %d, Skipped: %d - %s", completed+skipped, total, completed, skipped, message)
	}
	fmt.Println()
}

func (lt *LazyTool) createProjectUsingGHCLI(ctx context.Context, owner string) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", "project", "create", "--owner", owner, "--title", lt.tasksFile.ProjectTitle, "--format", "json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to create project using gh CLI: %w\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	var response struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return "", fmt.Errorf("failed to parse project creation output: %w\nOutput: %s", err, stdout.String())
	}

	if response.URL == "" {
		return "", fmt.Errorf("project URL is empty in the response\nOutput: %s", stdout.String())
	}

	return response.URL, nil
}

func (lt *LazyTool) getProjectNumberFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid project URL format")
	}
	return parts[len(parts)-1], nil
}

func (lt *LazyTool) createAllMilestones(ctx context.Context) error {
	for i, milestone := range lt.tasksFile.Milestones {
		lt.progressChan <- fmt.Sprintf("Creating milestone: %s", milestone.Title)
		milestoneNumber, err := lt.createOrGetMilestone(ctx, milestone)
		if err != nil {
			return fmt.Errorf("failed to create or get milestone %s: %w", milestone.Title, err)
		}
		lt.tasksFile.Milestones[i].Number = milestoneNumber
		atomic.AddInt32(&lt.tasksCompleted, 1)
	}
	return nil
}

func (lt *LazyTool) createAllIssues(ctx context.Context) error {
	for i, milestone := range lt.tasksFile.Milestones {
		for j, issue := range milestone.Issues {
			lt.progressChan <- fmt.Sprintf("Creating issue: %s", issue.Title)
			issueNumber, created, err := lt.createOrGetIssue(ctx, issue)
			if err != nil {
				return fmt.Errorf("failed to create or get issue %s: %w", issue.Title, err)
			}
			lt.tasksFile.Milestones[i].Issues[j].Number = issueNumber
			if created {
				atomic.AddInt32(&lt.tasksCompleted, 1)
			} else {
				atomic.AddInt32(&lt.tasksSkipped, 1)
			}
		}
	}
	return nil
}

func (lt *LazyTool) associateIssuesWithMilestones(ctx context.Context) error {
	for _, milestone := range lt.tasksFile.Milestones {
		for _, issue := range milestone.Issues {
			lt.progressChan <- fmt.Sprintf("Associating issue #%d with milestone #%d", issue.Number, milestone.Number)
			if err := lt.updateIssueMilestone(ctx, issue.Number, milestone.Number); err != nil {
				return fmt.Errorf("failed to associate issue #%d with milestone #%d: %w", issue.Number, milestone.Number, err)
			}
			atomic.AddInt32(&lt.tasksCompleted, 1)
		}
	}
	return nil
}

func (lt *LazyTool) addAllIssuesToProject(ctx context.Context, projectNumber string) error {
	for _, milestone := range lt.tasksFile.Milestones {
		for _, issue := range milestone.Issues {
			lt.progressChan <- fmt.Sprintf("Adding issue #%d to project", issue.Number)
			if err := lt.addIssueToProjectUsingGHCLI(ctx, issue.Number, projectNumber); err != nil {
				return fmt.Errorf("failed to add issue #%d to project: %w", issue.Number, err)
			}
			atomic.AddInt32(&lt.tasksCompleted, 1)
		}
	}
	return nil
}

func (lt *LazyTool) getUsername(ctx context.Context) (string, error) {
	var response struct {
		Login string `json:"login"`
	}
	if err := lt.client.Get("user", &response); err != nil {
		return "", err
	}
	return response.Login, nil
}

func (lt *LazyTool) createOrGetMilestone(ctx context.Context, milestone Milestone) (int, error) {
	existingMilestone, err := lt.getMilestoneByTitle(ctx, milestone.Title)
	if err != nil {
		return 0, fmt.Errorf("checking existing milestone: %w", err)
	}
	if existingMilestone != nil {
		atomic.AddInt32(&lt.tasksSkipped, 1)
		return existingMilestone.Number, nil
	}

	var response struct {
		Number int `json:"number"`
	}

	dueOn, err := lt.formatDate(milestone.DueOn)
	if err != nil {
		return 0, fmt.Errorf("formatting due date: %w", err)
	}

	payload := map[string]interface{}{
		"title":  milestone.Title,
		"state":  "open",
		"due_on": dueOn,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("marshaling milestone: %w", err)
	}
	url := fmt.Sprintf("repos/%s/%s/milestones", lt.owner, lt.repo)
	if err := lt.client.Post(url, bytes.NewReader(payloadBytes), &response); err != nil {
		return 0, fmt.Errorf("creating milestone: %w", err)
	}
	return response.Number, nil
}

func (lt *LazyTool) formatDate(dateStr string) (string, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("parsing date: %w", err)
	}
	return date.Format(time.RFC3339), nil
}

func (lt *LazyTool) getMilestoneByTitle(ctx context.Context, title string) (*Milestone, error) {
	var milestones []Milestone
	url := fmt.Sprintf("repos/%s/%s/milestones", lt.owner, lt.repo)
	if err := lt.client.Get(url, &milestones); err != nil {
		return nil, fmt.Errorf("fetching milestones: %w", err)
	}
	for _, m := range milestones {
		if m.Title == title {
			return &m, nil
		}
	}
	return nil, nil
}

func (lt *LazyTool) createOrGetIssue(ctx context.Context, issue Issue) (int, bool, error) {
	existingIssue, err := lt.getIssueByTitle(ctx, issue.Title)
	if err != nil {
		return 0, false, fmt.Errorf("checking existing issue: %w", err)
	}
	if existingIssue != nil {
		return existingIssue.Number, false, nil
	}

	var response struct {
		Number int `json:"number"`
	}
	payload := map[string]interface{}{
		"title": issue.Title,
		"body":  issue.Body,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, false, fmt.Errorf("marshaling issue: %w", err)
	}
	url := fmt.Sprintf("repos/%s/%s/issues", lt.owner, lt.repo)
	if err := lt.client.Post(url, bytes.NewReader(payloadBytes), &response); err != nil {
		return 0, false, fmt.Errorf("creating issue: %w", err)
	}
	return response.Number, true, nil
}

func (lt *LazyTool) getIssueByTitle(ctx context.Context, title string) (*Issue, error) {
	var issues []Issue
	url := fmt.Sprintf("repos/%s/%s/issues?state=all", lt.owner, lt.repo)
	if err := lt.client.Get(url, &issues); err != nil {
		return nil, fmt.Errorf("fetching issues: %w", err)
	}
	for _, i := range issues {
		if i.Title == title {
			return &i, nil
		}
	}
	return nil, nil
}

func (lt *LazyTool) updateIssueMilestone(ctx context.Context, issueNumber, milestoneNumber int) error {
	payload := map[string]interface{}{
		"milestone": milestoneNumber,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling issue update: %w", err)
	}
	url := fmt.Sprintf("repos/%s/%s/issues/%d", lt.owner, lt.repo, issueNumber)
	var response interface{}
	if err := lt.client.Patch(url, bytes.NewReader(payloadBytes), &response); err != nil {
		return fmt.Errorf("updating issue milestone: %w", err)
	}
	return nil
}

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

// Additional utility methods can be added here if needed

func (lt *LazyTool) logError(format string, args ...interface{}) {
	lt.logger.Printf(format, args...)
}

func (lt *LazyTool) logInfo(format string, args ...interface{}) {
	lt.logger.Printf(format, args...)
}
