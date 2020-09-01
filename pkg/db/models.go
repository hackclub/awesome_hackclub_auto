package db

type Project struct {
	Timestamp   string `json:"ts"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	GitHubURL   string `json:"github_url"`
	Description string `json:"description"`
}
