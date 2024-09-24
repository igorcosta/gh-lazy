package models

type Issue struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Number int    `json:"number,omitempty"`
}

type Milestone struct {
	Title  string  `json:"title"`
	DueOn  string  `json:"due_on"`
	Issues []Issue `json:"issues"`
	Number int     `json:"number,omitempty"`
}

type TasksFile struct {
	ProjectTitle string      `json:"projectTitle"`
	Milestones   []Milestone `json:"milestones"`
}
