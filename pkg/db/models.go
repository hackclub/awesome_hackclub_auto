package db

type ProjectStatus string

const (
	ProjectStatusIntent  ProjectStatus = "intent"
	ProjectStatusQueue   ProjectStatus = "queue"
	ProjectStatusProject ProjectStatus = "project"
	ProjectStatusDeleted ProjectStatus = "deleted"
)

type Project struct {
	ID     string        `json:"id,omitempty"`
	Fields ProjectFields `json:"fields"`
}

type ProjectFields struct {
	Status      ProjectStatus `json:"Status"`
	Timestamp   string        `json:"Timestamp"`
	UserID      string        `json:"User ID"`
	Name        string        `json:"Project Name"`
	GitHubURL   string        `json:"GitHub URL"`
	Description string        `json:"Description"`
	Language    string        `json:"Language"`
	Category    string        `json:"Category"`
	Channel     string        `json:"Channel"`
	Username    string        `json:"Username"`
}
