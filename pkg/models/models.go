package models

import "time"

type Issue struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Number int    `json:"number,omitempty"`
}

type Milestone struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	DueOn       time.Time `json:"due_on,omitempty"`
	State       string    `json:"state,omitempty"`
	Number      int       `json:"number,omitempty"`
}

type MilestoneWithIssues struct {
	Milestone
	Issues []Issue `json:"issues"`
}

type TasksFile struct {
	ProjectTitle string                `json:"projectTitle"`
	Milestones   []MilestoneWithIssues `json:"milestones"`
}

type IssueItem struct {
	Number     int    `json:"number"`
	Repository string `json:"repository"`
}

type Project struct {
	Closed           bool   `json:"closed"`
	ID               string `json:"id"`
	Number           int    `json:"number"`
	Title            string `json:"title"`
	URL              string `json:"url"`
	ShortDescription string `json:"shortDescription"`
}

type ProjectListResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int       `json:"totalCount"`
}
