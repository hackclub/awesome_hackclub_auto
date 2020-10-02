package gen

import "github.com/hackclub/awesome_hackclub_auto/pkg/db"

// Group projects based off category
func GroupProjects(projects []db.Project) map[string][]db.Project {
	groupedProjects := map[string][]db.Project{}
	for _, project := range projects {
		category := project.Fields.Category
		groupedProjects[category] = append(groupedProjects[category], project)
	}
	return groupedProjects
}
