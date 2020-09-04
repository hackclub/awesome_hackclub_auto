package handlers

import (
	"encoding/json"
	"fmt"
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
				metadata, _ := json.Marshal(struct {
					ProjectID string
					Channel   string
					TS        string
				}{
					ProjectID: parsed.ActionCallback.BlockActions[0].Value,
					Channel:   parsed.Channel.ID,
					TS:        parsed.Message.Timestamp,
				})
				_, err := client.OpenView(parsed.TriggerID, block_kit.SubmitModal(string(metadata), project))
				if err != nil {
					logoru.Error(err)
				}
			case db.ProjectStatusProject:
				// TODO
			}
		} else if actionID == "accept" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))
			project := db.GetProject(parsed.ActionCallback.BlockActions[0].Value)

			project.Status = db.ProjectStatusProject
			db.UpdateProject(project)

			_, _, _, err := client.UpdateMessage(parsed.Channel.ID, parsed.Message.Timestamp, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>* was approved by <@%s>", project.GitHubURL, project.Name, parsed.User.ID), false, false), nil, nil),
			))
			if err != nil {
				logoru.Error(err)
			}
			util.SendApprovedMessage(project)
		} else if actionID == "deny" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))

			privateMetadata, _ := json.Marshal(struct {
				ProjectID string
				ReviewTS  string
				Channel   string
			}{
				ProjectID: parsed.ActionCallback.BlockActions[0].Value,
				ReviewTS:  parsed.Message.Timestamp,
				Channel:   parsed.Channel.ID,
			})

			client.OpenView(parsed.TriggerID, slack.ModalViewRequest{
				Type:            "modal",
				CallbackID:      "deny",
				PrivateMetadata: string(privateMetadata),
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						slack.NewInputBlock("reason", slack.NewTextBlockObject("plain_text", "Reason", false, false), &slack.PlainTextInputBlockElement{
							Type:      "plain_text_input",
							ActionID:  "reason",
							Multiline: true,
						}),
					},
				},
				Title:  slack.NewTextBlockObject("plain_text", "Deny Project", false, false),
				Submit: slack.NewTextBlockObject("plain_text", "Deny", false, false),
				Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
			})
		}
	case slack.InteractionTypeViewSubmission:
		values := parsed.View.State.Values

		if parsed.View.CallbackID == "submit" {
			var metadata struct {
				ProjectID string
				Channel   string
				TS        string
			}

			json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)

			project := db.GetProject(metadata.ProjectID)

			project.Status = db.ProjectStatusQueue

			project.Description = values["description"]["description"].Value
			project.Name = values["name"]["name"].Value
			project.GitHubURL = values["url"]["url"].Value

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			db.UpdateProject(project)
			_, _, _, err := client.UpdateMessage(metadata.Channel, metadata.TS, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Your project, *%s*, has successfully been submitted! :tada:", project.Name), false, false), nil, nil),
			))
			if err != nil {
				logoru.Error(err)
			}
			util.SendReviewMessage(project)
		} else if parsed.View.CallbackID == "deny" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))

			var metadata struct {
				ProjectID string
				ReviewTS  string
				Channel   string
			}
			json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)

			project := db.GetProject(metadata.ProjectID)

			db.DeleteProject(project)
			_, _, _, err := client.UpdateMessage(metadata.Channel, metadata.ReviewTS, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>* was denied by <@%s>\n*Reason*: %s", project.GitHubURL, project.Name, parsed.User.ID, values["reason"]["reason"].Value), false, false), nil, nil),
			))
			if err != nil {
				logoru.Error(err)
			}
			util.SendDeniedMessage(project, values["reason"]["reason"].Value)
		}
	}
}
