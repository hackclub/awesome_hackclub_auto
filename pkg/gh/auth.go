package gh

import (
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
