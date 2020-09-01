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
		action := block_kit.SlackActionID{}
		json.Unmarshal([]byte(parsed.ActionCallback.BlockActions[0].ActionID), &action)

		if action.Action == "submit" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))
			if db.ProjectIsInQueue(action.Timestamp) {
				_, err := client.OpenView(parsed.TriggerID, block_kit.AlreadyInQueue())
				if err != nil {
					logoru.Error(err)
				}
				return
			}
			_, err := client.OpenView(parsed.TriggerID, block_kit.SubmitModal(block_kit.SlackPrivateMetadata{
				Timestamp: action.Timestamp,
			}, action))
			if err != nil {
				logoru.Error(err)
			}
		}
	case slack.InteractionTypeViewSubmission:
		values := parsed.View.State.Values

		metadata := block_kit.SlackPrivateMetadata{}
		json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)

		util.AddProjectToQueue(db.Project{
			Timestamp:   metadata.Timestamp,
			Name:        values["name"]["name"].Value,
			Description: values["description"]["description"].Value,
			GitHubURL:   values["url"]["url"].Value,
			UserID:      parsed.User.ID,
		})
	}
}
