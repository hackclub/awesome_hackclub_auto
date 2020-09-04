package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hackclub/awesome_hackclub_auto/pkg/block_kit"
	"github.com/hackclub/awesome_hackclub_auto/pkg/util"

	"github.com/Matt-Gleich/logoru"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/slack-go/slack"
)

func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logoru.Error(err)
	}
	r.Form, err = url.ParseQuery(string(buf))
	if err != nil {
		logoru.Error(err)
	}

	parsed := slack.InteractionCallback{}

	err = json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)
	if err != nil {
		logoru.Error(err)
	}

	switch parsed.Type {
	case slack.InteractionTypeBlockActions:
		actionID := parsed.ActionCallback.BlockActions[0].ActionID
		if actionID == "submit" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))

			project := db.GetProject(parsed.ActionCallback.BlockActions[0].Value)

			switch project.Status {
			case db.ProjectStatusQueue:
				_, err := client.OpenView(parsed.TriggerID, block_kit.AlreadyInQueue())
				if err != nil {
					logoru.Error(err)
				}
				return
			case db.ProjectStatusIntent:
				_, err := client.OpenView(parsed.TriggerID, block_kit.SubmitModal(parsed.ActionCallback.BlockActions[0].Value, project))
				if err != nil {
					logoru.Error(err)
				}
			case db.ProjectStatusProject:
				// TODO
			}
		} else if actionID == "accept" {
			logoru.Debug("Approve project ", parsed.ActionCallback.BlockActions[0].Value)
		} else if actionID == "deny" {
			logoru.Debug("Deny project ", parsed.ActionCallback.BlockActions[0].Value)
		}
	case slack.InteractionTypeViewSubmission:
		values := parsed.View.State.Values

		project := db.GetProject(parsed.View.PrivateMetadata)

		project.Status = db.ProjectStatusQueue

		project.Description = values["description"]["description"].Value
		project.Name = values["name"]["name"].Value
		project.GitHubURL = values["url"]["url"].Value

		db.UpdateProject(project)
		util.SendReviewMessage(project)
	}
}
