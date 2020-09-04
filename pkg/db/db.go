package db

import (
	"context"

	_ "github.com/go-kivik/couchdb"
	"github.com/google/uuid"

	"github.com/Matt-Gleich/logoru"
	"github.com/go-kivik/kivik"
)

func CreateProjectIntent(project Project) string {
	db, err := kivik.New("couch", "http://admin:password@db:5984/")
	if err != nil {
		logoru.Error(err)
		return ""
	}

	project.ID = uuid.New().String()
	project.Status = ProjectStatusIntent

	id, _, err := db.DB(context.TODO(), "projects").CreateDoc(context.TODO(), project)
	if err != nil {
		logoru.Error(err)
		return ""
	}

	return id
}

func GetProject(id string) Project {
	db, err := kivik.New("couch", "http://admin:password@db:5984/")
	if err != nil {
		logoru.Error(err)
		return Project{}
	}

	row := db.DB(context.TODO(), "projects").Get(context.TODO(), id)
	project := Project{}
	err = row.ScanDoc(&project)
	if err != nil {
		logoru.Error(err)
	}
	return project
}

func UpdateProject(newProject Project) {
	db, err := kivik.New("couch", "http://admin:password@db:5984/")
	if err != nil {
		logoru.Error(err)
		return
	}

	_, err = db.DB(context.TODO(), "projects").Put(context.TODO(), newProject.ID, newProject)
	if err != nil {
		logoru.Error(err)
		return
	}
}

func DeleteProject(project Project) {
	db, err := kivik.New("couch", "http://admin:password@db:5984/")
	if err != nil {
		logoru.Error(err)
		return
	}

	_, err = db.DB(context.TODO(), "projects").Delete(context.TODO(), project.ID, project.Rev)
	if err != nil {
		logoru.Error(err)
		return
	}
}
