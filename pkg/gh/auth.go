package gh

import (
	"net/http"
	"os"

	"github.com/Matt-Gleich/logoru"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
)

// Authenticate with GitHub using the secret ssh key
// Return a github client instance
func Auth() *github.Client {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 81465, 99, "private-key.pem")
	if err != nil {
		logoru.Critical("Failed to authenticate with GitHub using private key;", err)
		os.Exit(1)
	}
	return github.NewClient(&http.Client{Transport: itr})
}
