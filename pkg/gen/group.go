package gen

import (
	"sort"

	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
)

// Group projects based off category
func GroupProjects(projects []db.Project) map[string][]db.Project {
	groupedProjects := map[string][]db.Project{}
	for _, project := range projects {
		category := project.Fields.Category
		groupedProjects[category] = append(groupedProjects[category], project)
	}

	// Sorting alphabetically based off the categories
	sortedCategories := []string{}
	for c := range groupedProjects {
		sortedCategories = append(sortedCategories, c)
	}
	sort.Strings(sortedCategories)

	sortedProjects := map[string][]db.Project{}
	for _, c := range sortedCategories {
		sortedProjects[c] = groupedProjects[c]
	}

	return sortedProjects
}
