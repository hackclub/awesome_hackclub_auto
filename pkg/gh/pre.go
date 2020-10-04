package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/hackclub/awesome_hackclub_auto/pkg/config"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
)

type RepoMetaData struct {
	Language    *string
	Description *string
}

// Get information to match up with project fields
func RepoInfo(client *github.Client, owner string, name string) RepoMetaData {
	repo, _, err := client.Repositories.Get(context.Background(), owner, name)
	if err != nil {
		logging.Log(fmt.Sprintf("Failed to get info from repo: %v", err), "warning", false)
		return RepoMetaData{}
	}
	// Checking to see if the language is a valid language
	for _, lang := range config.Languages {
		if lang == *repo.Language {
			return RepoMetaData{
				Language:    repo.Language,
				Description: repo.Description,
			}
		}
	}
	logging.Log(*repo.Language+" isn't a supported language", "warning", false)
	return RepoMetaData{
		Description: repo.Description,
	}
}
