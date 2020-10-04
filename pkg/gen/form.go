package gen

import (
	"fmt"

	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
)

const header = `# 😎 Awesome Hack Club [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)
> A collection of super awesome projects made by hackclubbers
`

const footer = `
## ✨ Adding your project

Are you a hackclubber and want your project here? Its easy! Just react to your ship message with the awesome emoji (:awesome:) and you will get a DM from the awesome-hackclub bot. Click Submit, fill-in the information, and wait for your project to be reviewed. Once the review is complete you will get another DM from the awesome-hackclub bot. Thats it! As long as it was approved you can check back in this repo!
`

// Create the README
func FormREADME(groupedProjects map[string][]db.Project) string {
	var body string
	for category, projects := range groupedProjects {
		body = fmt.Sprintf("%v\n## %v\n", body, category)
		for _, project := range projects {
			body = fmt.Sprintf(
				"%v- [%v](%v) - (%v) %v\n",
				body,
				fmt.Sprintf("%s/%s", project.Fields.Username, project.Fields.Name),
				project.Fields.GitHubURL,
				project.Fields.Language,
				project.Fields.Description,
			)
		}
	}
	return header + body + footer
}
