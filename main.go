package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hackclub/awesome_hackclub_auto/pkg/config"
	"github.com/hackclub/awesome_hackclub_auto/pkg/handlers"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
	"github.com/honeybadger-io/honeybadger-go"
)

func main() {
	honeybadger.Configure(honeybadger.Configuration{APIKey: os.Getenv("HONEYBADGER_API_KEY")})

	notSet := logging.GetUnsetEnvVars([]string{"SLACK_TOKEN", "SLACK_SIGNING_SECRET", "REVIEW_CHANNEL", "AIRTABLE_API_KEY", "AIRTABLE_BASE_ID", "GH_APP_ID", "GH_INSTALLATION_ID", "GH_PRIVATE_KEY"})
	if len(notSet) > 0 {
		logging.Log(fmt.Sprintf("Startup error: the following env vars are unset: %s", strings.Join(notSet, ", ")), "critical", false)
	}

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
