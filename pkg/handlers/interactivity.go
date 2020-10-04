package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hackclub/awesome_hackclub_auto/pkg/block_kit"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/hackclub/awesome_hackclub_auto/pkg/gen"
	"github.com/hackclub/awesome_hackclub_auto/pkg/gh"
	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
	"github.com/hackclub/awesome_hackclub_auto/pkg/util"
	"github.com/slack-go/slack"
)

func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Log(err, "error", false)
	}

	if !util.VerifySlackRequest(r, buf) {
		logging.Log("invalid Slack request", "warning", false)
		_, err = w.Write(nil)
		if err != nil {
			logging.Log(err, "error", false)
		}
		return
	}

	r.Form, err = url.ParseQuery(string(buf))
	if err != nil {
		logging.Log(err, "error", false)
	}

	parsed := slack.InteractionCallback{}

	err = json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)
	if err != nil {
		logging.Log(err, "error", false)
	}

	switch parsed.Type {
	case slack.InteractionTypeBlockActions:
		actionID := parsed.ActionCallback.BlockActions[0].ActionID
		if actionID == "submit" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))

			project := db.GetProject(parsed.ActionCallback.BlockActions[0].Value)

			switch project.Fields.Status {
			case db.ProjectStatusQueue:
				_, err := client.OpenView(parsed.TriggerID, block_kit.AlreadyInQueue())
				if err != nil {
					logging.Log(err, "error", false)
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
				_, err := client.OpenView(parsed.TriggerID, block_kit.SubmitModal(string(metadata), project.Fields))
				if err != nil {
					logging.Log(err, "error", false)
				}
			case db.ProjectStatusProject:
				// TODO
			}
		} else if actionID == "accept" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))
			project := db.GetProject(parsed.ActionCallback.BlockActions[0].Value)

			if project.Fields.Status == db.ProjectStatusProject {
				// Exit if the project's already approved
				return
			}

			project.Fields.Status = db.ProjectStatusProject
			db.UpdateProject(project)

			_, _, _, err := client.UpdateMessage(parsed.Channel.ID, parsed.Message.Timestamp, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>* was approved by <@%s>", project.Fields.GitHubURL, project.Fields.Name, parsed.User.ID), false, false), nil, nil),
			))
			if err != nil {
				logging.Log(err, "error", false)
			}
			projects := gen.GroupProjects(db.GetAllProjects())
			readme := gen.FormREADME(projects)
			gh.UpdateREADME(readme, project)
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

			_, err = client.OpenView(parsed.TriggerID, slack.ModalViewRequest{
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
			if err != nil {
				logging.Log(err, "error", false)
			}
		}
	case slack.InteractionTypeViewSubmission:
		values := parsed.View.State.Values

		if parsed.View.CallbackID == "submit" {
			if !util.IsValidProjectURL(values["url"]["url"].Value) {
				w.Header().Add("Content-Type", "application/json")
				resp, _ := json.Marshal(slack.ViewSubmissionResponse{
					ResponseAction: slack.RAErrors,
					Errors: map[string]string{
						"url": "This isn't a valid GitHub URL. It should look like the following: https://github.com/hackclub/hackclub",
					},
				})

				_, err = w.Write(resp)
				if err != nil {
					logging.Log(err, "error", false)
				}
				return
			}
			var metadata struct {
				ProjectID string
				Channel   string
				TS        string
			}

			err := json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)
			if err != nil {
				logging.Log(err, "error", false)
			}

			project := db.GetProject(metadata.ProjectID)

			project.Fields.Status = db.ProjectStatusQueue

			project.Fields.Description = values["description"]["description"].Value
			project.Fields.Name = values["name"]["name"].Value
			project.Fields.GitHubURL = values["url"]["url"].Value
			project.Fields.Category = values["category"]["category"].SelectedOption.Value
			project.Fields.Language = values["language"]["language"].SelectedOption.Value

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			db.UpdateProject(project)
			_, _, _, err = client.UpdateMessage(metadata.Channel, metadata.TS, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Your project, *%s*, has successfully been submitted! :tada: You'll get another DM once it's been added.", project.Fields.Name), false, false), nil, nil),
			))
			if err != nil {
				logging.Log(err, "error", false)
			}
			util.SendReviewMessage(project)
		} else if parsed.View.CallbackID == "deny" {
			client := slack.New(os.Getenv("SLACK_TOKEN"))

			var metadata struct {
				ProjectID string
				ReviewTS  string
				Channel   string
			}
			err := json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)
			if err != nil {
				logging.Log(err, "error", false)
			}

			project := db.GetProject(metadata.ProjectID)

			db.DeleteProject(project)
			_, _, _, err = client.UpdateMessage(metadata.Channel, metadata.ReviewTS, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>* was denied by <@%s>\n*Reason*: %s", project.Fields.GitHubURL, project.Fields.Name, parsed.User.ID, values["reason"]["reason"].Value), false, false), nil, nil),
			))
			if err != nil {
				logging.Log(err, "error", false)
			}
			util.SendDeniedMessage(project, values["reason"]["reason"].Value)
		}
	}
}
