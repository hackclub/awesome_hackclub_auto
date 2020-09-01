package block_kit

type SlackActionID struct {
	Action      string `json:"act"`
	GitHubURL   string `json:"gh,omitempty"`
	ProjectName string `json:"pro,omitempty"`
	Timestamp   string `json:"ts"`
}

type SlackPrivateMetadata struct {
	Timestamp string `json:"ts"`
}
