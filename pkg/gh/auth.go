package gh

import (
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
)

// Authenticate with GitHub using the secret ssh key
// Return a github client instance
func Auth() *github.Client {
	appId, err := strconv.Atoi(os.Getenv("GH_APP_ID"))
	if err != nil {
		logging.Log("The GH_APP_ID environment variable should be set, and a number", "critical", false)
	}

	installationId, err := strconv.Atoi(os.Getenv("GH_INSTALLATION_ID"))
	if err != nil {
		logging.Log("The GH_INSTALLATION_ID environment variable should be set, and a number", "critical", false)
	}

	itr, err := ghinstallation.New(http.DefaultTransport, int64(appId), int64(installationId), LoadPrivateKey())
	if err != nil {
		logging.Log("Failed to authenticate with GitHub using private key; "+err.Error(), "critical", false)
		os.Exit(1)
	}
	return github.NewClient(&http.Client{Transport: itr})
}
