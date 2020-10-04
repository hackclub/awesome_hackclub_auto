package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/hackclub/awesome_hackclub_auto/pkg/config"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
)

type RepoMetaData struct {
	Language    string
	Description string
	Valid       bool
}

// Get information to match up with project fields
func RepoInfo(client *github.Client, owner string, name string) RepoMetaData {
	repo, _, err := client.Repositories.Get(context.Background(), owner, name)
	if err != nil {
		logging.Log(fmt.Sprintf("Failed to get info from repo: %v", err), "warning", true)
		return RepoMetaData{Valid: false}
	}

	metadata := RepoMetaData{}

	// Checking to see if the language is a valid language
	if *repo.Language != "" {
		// Repo has a language
		for _, lang := range config.Languages {
			if lang == *repo.Language {
				metadata.Language = *repo.Language
			}
		}
	}

	if metadata.Language == "" {
		logging.Log(*repo.Language+" isn't a supported language", "warning", false)
	}

	metadata.Valid = !*repo.Private

	if *repo.Description != "" {
		metadata.Description = *repo.Description
	}

	return metadata
}
