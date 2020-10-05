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

Are you a hackclubber and want your project here? Its easy! Just react to your ship message with the awesome emoji (:awesome:) and you will get a DM from the awesome-hackclub bot. Click Submit, fill-in the information, and wait for your project to be reviewed. Once the review is complete you will get another DM from the awesome-hackclub bot. Thats it! As long as it was approved you can check back in this repo! **Please only submit public, finished repositories that align with the [Hack Club code of conduct](https://hackclub.com/conduct/)**.
`

// Create the README
func FormREADME(groupedProjects map[string][]db.Project) string {
	var body string
	for category, projects := range groupedProjects {
		body = fmt.Sprintf("%v\n## %v\n", body, category)
		for _, project := range projects {
			if project.Fields.Description != "" {
				project.Fields.Description = "_" + project.Fields.Description + "_"
			}
			body = fmt.Sprintf(
				"%v- **[%v](%v)** - [@%[4]v](https://github.com/%[4]v) - **(%v)** %v\n",
				body,
				project.Fields.Name,
				project.Fields.GitHubURL,
				project.Fields.Username,
				project.Fields.Language,
				project.Fields.Description,
			)
		}
	}
	return header + body + footer
}
