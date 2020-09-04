package main

import (
	"net/http"
	"os"

	"github.com/Matt-Gleich/logoru"
	"github.com/gorilla/mux"
	"github.com/hackclub/awesome_hackclub_auto/pkg/config"
	"github.com/hackclub/awesome_hackclub_auto/pkg/handlers"
)

func main() {
	config.PopulateConfig()

	r := mux.NewRouter()

	r.HandleFunc("/slack/events", handlers.HandleEvents).Methods("POST")
	r.HandleFunc("/slack/interactivity", handlers.HandleInteractivity).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		logoru.Critical(err)
	}
}
