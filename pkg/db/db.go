package db

import (
	"os"

	"github.com/Matt-Gleich/logoru"
	"github.com/brianloveswords/airtable"
)

func projectsTable() airtable.Table {
	client := airtable.Client{
		APIKey: os.Getenv("AIRTABLE_API_KEY"),
		BaseID: os.Getenv("AIRTABLE_BASE_ID"),
	}

	return client.Table("Projects")
}

func CreateProjectIntent(fields ProjectFields) string {
	projects := projectsTable()

	fields.Status = ProjectStatusIntent
	project := Project{
		Fields: fields,
	}

	err := projects.Create(&project)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}

	return project.ID
}

func GetProject(id string) Project {
	projects := projectsTable()

	project := Project{}
	err := projects.Get(id, &project)

	if err != nil {
		return Project{}
	}

	return project
}

func UpdateProject(newProject Project) {
	projects := projectsTable()

	err := projects.Update(&newProject)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}
}

func DeleteProject(project Project) {
	projects := projectsTable()

	project.Fields.Status = ProjectStatusDeleted

	err := projects.Update(&project)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}
}

func GetAllProjects() []Project {
	table := projectsTable()

	projects := []Project{}

	err := table.List(&projects, &airtable.Options{Filter: "Status = 'project'", Sort: airtable.Sort{
		[2]string{"Category", "asc"},
		[2]string{"Language", "asc"},
		[2]string{"Name", "asc"},
	}})
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}

	return projects
}
