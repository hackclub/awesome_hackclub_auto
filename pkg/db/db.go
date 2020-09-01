package db

import (
	"context"

	_ "github.com/go-kivik/couchdb"

	"github.com/Matt-Gleich/logoru"
	"github.com/go-kivik/kivik"
)

func AddProjectToQueue(project Project) {
	db, err := kivik.New("couch", "http://admin:password@db:5984/")
	if err != nil {
		logoru.Error(err)
	}

	_, _, err = db.DB(context.TODO(), "queue").CreateDoc(context.TODO(), project)
	if err != nil {
		logoru.Error(err)
	}
}
