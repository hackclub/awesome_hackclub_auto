package util

import "testing"

func TestParseGitHubURL(t *testing.T) {
	url, valid := ParseGitHubURL("https://github.com/cjdenio/replier")
	if !valid {
		t.Error()
	} else if url.Owner != "cjdenio" {
		t.Error()
	} else if url.Name != "replier" {
		t.Error()
	}

	url, valid = ParseGitHubURL("https://github.com/cjdenio/replier is my project")
	if valid {
		t.Log(url)
		t.Error()
	}
}
