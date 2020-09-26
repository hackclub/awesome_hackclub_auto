package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Matt-Gleich/logoru"
	"github.com/hackclub/awesome_hackclub_auto/pkg/block_kit"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/hackclub/awesome_hackclub_auto/pkg/util"
	"github.com/slack-go/slack"
)

// HandleInteractivity is called when someone either clicks a button anywhere in the app or submits a modal
func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}
	r.Form, err = url.ParseQuery(string(buf))
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}

	parsed := slack.InteractionCallback{}

	// Attempt to parse the request body as JSON
	err = json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}

	switch parsed.Type {
	case slack.InteractionTypeBlockActions:
		// It was a button press

		actionID := parsed.ActionCallback.BlockActions[0].ActionID
		if actionID == "submit" {
			// Someone clicked the "Submit" button, so let's open the modal

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			project := db.GetProject(parsed.ActionCallback.BlockActions[0].Value)

			// Here, JSON-ify some state data
			metadata, _ := json.Marshal(struct {
				ProjectID string
				Channel   string
				TS        string
			}{
				ProjectID: parsed.ActionCallback.BlockActions[0].Value,
				Channel:   parsed.Channel.ID,
				TS:        parsed.Message.Timestamp,
			})

			// Show them the modal
			_, err := client.OpenView(parsed.TriggerID, block_kit.SubmitModal(string(metadata), project.Fields))
			if err != nil {
				logoru.Error(err)
				os.Exit(1)
			}
		} else if actionID == "accept" {
			// A project has been accepted into the awesome-hackclub repo!

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			// Get the project's metadata
			project := db.GetProject(parsed.ActionCallback.BlockActions[0].Value)

			// Set the status to Project (which means it got approved)
			project.Fields.Status = db.ProjectStatusProject

			// Push the new project to the DB
			db.UpdateProject(project)

			// Update the message to say "Cool Project was approved by @Cool Person"
			_, _, _, err := client.UpdateMessage(parsed.Channel.ID, parsed.Message.Timestamp, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>* was approved by <@%s>", project.Fields.GitHubURL, project.Fields.Name, parsed.User.ID), false, false), nil, nil),
			))
			if err != nil {
				logoru.Error(err)
				os.Exit(1)
			}

			// DM the project's creator to let them know it's been approved
			util.SendApprovedMessage(project)
		} else if actionID == "deny" {
			// A project got denied, so open a modal to ask for a reason (THE MODAL SUBMISSION IS HANDLED ELSEWHERE)

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			// Pass some state data through the modal
			privateMetadata, _ := json.Marshal(struct {
				ProjectID string
				ReviewTS  string
				Channel   string
			}{
				ProjectID: parsed.ActionCallback.BlockActions[0].Value,
				ReviewTS:  parsed.Message.Timestamp,
				Channel:   parsed.Channel.ID,
			})

			// Open a modal where the denier can provide a reason
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
				logoru.Error(err)
				os.Exit(1)
			}
		}
	case slack.InteractionTypeViewSubmission:
		// Somewhere, a modal was submitted

		values := parsed.View.State.Values

		if parsed.View.CallbackID == "submit" {
			// Someone just submitted the "Submit Project" modal

			// Check to make sure the project URL is a valid GitHub URL
			if !util.IsValidProjectURL(values["url"]["url"].Value) {
				// Reject the modal submission
				w.Header().Add("Content-Type", "application/json")
				resp, _ := json.Marshal(slack.ViewSubmissionResponse{
					ResponseAction: slack.RAErrors,
					Errors: map[string]string{
						"url": "This isn't a valid GitHub URL. It should look like the following: https://github.com/hackclub/hackclub",
					},
				})

				_, err = w.Write(resp)
				if err != nil {
					logoru.Error(err)
					os.Exit(1)
				}
				return
			}

			// Prepare to parse the JSON state data
			var metadata struct {
				ProjectID string
				Channel   string
				TS        string
			}

			// Actually parse it
			err := json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)
			if err != nil {
				logoru.Error(err)
			}

			project := db.GetProject(metadata.ProjectID)

			// Set the status to "queue"
			project.Fields.Status = db.ProjectStatusQueue

			// Set various fields on the project (based on the user's input)
			project.Fields.Description = values["description"]["description"].Value
			project.Fields.Name = values["name"]["name"].Value
			project.Fields.GitHubURL = values["url"]["url"].Value
			project.Fields.Category = values["category"]["category"].SelectedOption.Value
			project.Fields.Language = values["language"]["language"].SelectedOption.Value

			// Write the new project to the database
			db.UpdateProject(project)

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			_, _, _, err = client.UpdateMessage(metadata.Channel, metadata.TS, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Your project, *%s*, has successfully been submitted! :tada: You'll get another DM once it's been added.", project.Fields.Name), false, false), nil, nil),
			))
			if err != nil {
				logoru.Error(err)
				os.Exit(1)
			}

			// Post a message to the review channel
			util.SendReviewMessage(project)
		} else if parsed.View.CallbackID == "deny" {
			// The deny modal was submitted

			client := slack.New(os.Getenv("SLACK_TOKEN"))

			var metadata struct {
				ProjectID string
				ReviewTS  string
				Channel   string
			}
			err := json.Unmarshal([]byte(parsed.View.PrivateMetadata), &metadata)
			if err != nil {
				logoru.Error(err)
				os.Exit(1)
			}

			project := db.GetProject(metadata.ProjectID)

			// Delete the project
			db.DeleteProject(project)

			_, _, _, err = client.UpdateMessage(metadata.Channel, metadata.ReviewTS, slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>* was denied by <@%s>\n*Reason*: %s", project.Fields.GitHubURL, project.Fields.Name, parsed.User.ID, values["reason"]["reason"].Value), false, false), nil, nil),
			))
			if err != nil {
				logoru.Error(err)
				os.Exit(1)
			}

			// Let the user know that their project was denied
			util.SendDeniedMessage(project, values["reason"]["reason"].Value)
		}
	}
}
