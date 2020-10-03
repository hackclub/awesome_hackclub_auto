package db

import (
	"fmt"
	"os"

	"github.com/brianloveswords/airtable"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
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
		logging.Log(err, "error", false)
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

	logging.Log(fmt.Sprintf("Got project %s from Airtable", project.Fields.Name), "info", false)
	return project
}

func UpdateProject(newProject Project) {
	projects := projectsTable()

	err := projects.Update(&newProject)
	if err != nil {
		logging.Log(err, "error", false)
	}
	logging.Log(fmt.Sprintf("Updated project %s in Airtable", newProject.Fields.Name), "info", false)
}

func DeleteProject(project Project) {
	projects := projectsTable()

	project.Fields.Status = ProjectStatusDeleted

	err := projects.Update(&project)
	if err != nil {
		logging.Log(err, "error", false)
	}
	logging.Log(fmt.Sprintf("Removed project %s from Airtable", project.Fields.Name), "info", false)
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
		logging.Log(err, "error", false)
	}
	logging.Log("Got all projects from Airtable", "info", false)
	return projects
}
