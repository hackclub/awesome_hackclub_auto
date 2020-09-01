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

func ProjectIsInQueue(ts string) bool {
	db, err := kivik.New("couch", "http://admin:password@db:5984/")
	if err != nil {
		logoru.Error(err)
	}
	rows, err := db.DB(context.TODO(), "queue").Find(context.TODO(), map[string]map[string]string{
		"selector": {
			"ts": ts,
		},
	})
	if err != nil {
		logoru.Error(err)
	}

	logoru.Debug(rows.TotalRows())
	logoru.Debug(ts)

	isInQueue := rows.Next()
	rows.Close()

	return isInQueue
}
