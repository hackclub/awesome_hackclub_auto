package gh

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Matt-Gleich/logoru"
	"github.com/google/go-github/v32/github"
)

func UpdateREADME(content string) {
	client := Auth()

	var (
		message = "✨ Add project"
		sha     = getSHA()
		branch  = "master"
	)

	_, _, err := client.Repositories.UpdateFile(
		context.Background(),
		"hackclub",
		"awesome-hackclub",
		"README.md",
		&github.RepositoryContentFileOptions{
			Message: &message,
			Content: []byte(content),
			SHA:     &sha,
			Branch:  &branch,
		},
	)
	if err != nil {
		logoru.Error("Failed to push change to repo;", err)
	}
	logoru.Success("Pushed changes to repo")
}

// Get the current SHA for the file
func getSHA() string {
	// Making request
	resp, err := http.Get("https://api.github.com/repos/hackclub/awesome-hackclub/contents/README.md")
	if err != nil || resp.StatusCode != http.StatusOK {
		logoru.Error("Failed to get SHA for README.md;", err)
	}
	defer resp.Body.Close()

	// Parsing response
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logoru.Error("Failed to get binary from response;", err)
	}
	var data struct {
		Sha string `json:"sha"`
	}
	err = json.Unmarshal(bin, &data)
	if err != nil {
		logoru.Error("Failed to parse json from response;", err)
	}
	return data.Sha
}