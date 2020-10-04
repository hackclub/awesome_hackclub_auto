package util

import "regexp"

type ParsedGitHubURL struct {
	Owner string
	Name  string
	URL   string
}

func ParseGitHubURL(url string) (ParsedGitHubURL, bool) {
	re := regexp.MustCompile(`^https?:\/\/github\.com\/([^\/>\|\s]+)\/([^\/>\|\s]+)$`).FindStringSubmatch(url)
	if re == nil {
		return ParsedGitHubURL{}, false
	}

	return ParsedGitHubURL{
		Owner: re[1],
		Name:  re[2],
		URL:   re[0],
	}, true
}

func ParseGitHubURLInString(url string) (ParsedGitHubURL, bool) {
	re := regexp.MustCompile(`https?:\/\/github\.com\/([^\/>\|\s]+)\/([^\/>\|\s]+)`).FindStringSubmatch(url)
	if re == nil {
		return ParsedGitHubURL{}, false
	}

	return ParsedGitHubURL{
		Owner: re[1],
		Name:  re[2],
		URL:   re[0],
	}, true
}
