package db

import (
	"os"

	"github.com/brianloveswords/airtable"

	"github.com/Matt-Gleich/logoru"
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
	}
}

func DeleteProject(project Project) {
	projects := projectsTable()

	project.Fields.Status = ProjectStatusDeleted

	err := projects.Update(&project)
	if err != nil {
		logoru.Error(err)
	}
}

func GetAllProjects() []Project {
	table := projectsTable()

	projects := []Project{}

	err := table.List(&projects, &airtable.Options{Filter: "Status = 'project'"})
	if err != nil {
		logoru.Error(err)
	}

	return projects
}
