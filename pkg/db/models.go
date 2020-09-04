package db

type ProjectStatus string

const (
	ProjectStatusIntent  ProjectStatus = "intent"
	ProjectStatusQueue   ProjectStatus = "queue"
	ProjectStatusProject ProjectStatus = "project"
)

type Project struct {
	ID          string        `json:"_id"`
	Rev         string        `json:"_rev"`
	Status      ProjectStatus `json:"status"`
	Timestamp   string        `json:"ts"`
	UserID      string        `json:"user_id"`
	Name        string        `json:"name"`
	GitHubURL   string        `json:"github_url"`
	Description string        `json:"description"`
}
