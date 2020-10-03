package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hackclub/awesome_hackclub_auto/pkg/config"
	"github.com/hackclub/awesome_hackclub_auto/pkg/handlers"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
	"github.com/honeybadger-io/honeybadger-go"
)

func main() {
	honeybadger.Configure(honeybadger.Configuration{APIKey: os.Getenv("HONEYBADGER_API_KEY")})

	config.PopulateConfig()

	r := mux.NewRouter()

	r.HandleFunc("/slack/events", handlers.HandleEvents).Methods("POST")
	r.HandleFunc("/slack/interactivity", handlers.HandleInteractivity).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err := http.ListenAndServe(":"+port, honeybadger.Handler(r))
	if err != nil {
		logging.Log(err, "critical", false)
		os.Exit(1)
	}
}
