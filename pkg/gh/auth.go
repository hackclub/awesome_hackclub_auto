package gh

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/Matt-Gleich/logoru"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
)

// Authenticate with GitHub using the secret ssh key
// Return a github client instance
func Auth() *github.Client {
	appId, err := strconv.Atoi(os.Getenv("GH_APP_ID"))
	if err != nil {
		logoru.Critical("The GH_APP_ID environment variable should be set, and a number")
	}

	installationId, err := strconv.Atoi(os.Getenv("GH_INSTALLATION_ID"))
	if err != nil {
		logoru.Critical("The GH_INSTALLATION_ID environment variable should be set, and a number")
	}

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, int64(appId), int64(installationId), "private-key.pem")
	if err != nil {
		logoru.Critical("Failed to authenticate with GitHub using private key;", err)
		os.Exit(1)
	}
	return github.NewClient(&http.Client{Transport: itr})
}

// Get the current SHA for the file
func GetSHA() string {
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
