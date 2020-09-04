package util

import (
	"fmt"
	"os"

	"github.com/Matt-Gleich/logoru"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/slack-go/slack"
)

func SendReviewMessage(project db.Project) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	description := project.Description
	if description == "" {
		description = "_<no description>_"
	}

	_, _, err := client.PostMessage(os.Getenv("REVIEW_CHANNEL"), slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s> just submitted a project for review: <%s|%s>", project.UserID, project.GitHubURL, project.Name), false, false), []*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Description*\n%s", description), false, false),
		}, nil),
		slack.NewActionBlock(
			"",
			slack.ButtonBlockElement{
				Type:     "button",
				Style:    slack.StylePrimary,
				Text:     slack.NewTextBlockObject("plain_text", "Approve", false, false),
				ActionID: "accept",
				Value:    project.ID,
				Confirm: slack.NewConfirmationBlockObject(
					slack.NewTextBlockObject("plain_text", "Approve?", false, false),
					slack.NewTextBlockObject("plain_text", "Are you sure you want to add this project to awesome-hackclub?", false, false),
					slack.NewTextBlockObject("plain_text", "Yes", false, false),
					slack.NewTextBlockObject("plain_text", "No", false, false),
				),
			},
			slack.NewButtonBlockElement("deny", project.ID, slack.NewTextBlockObject("plain_text", "Deny", false, false)).WithStyle(slack.StyleDanger),
		),
	))
	if err != nil {
		logoru.Error(err)
	}
}

func SendApprovedMessage(project db.Project) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	_, _, err := client.PostMessage(project.UserID, slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("WOO-HOO!! :tada: Your project, *%s*, has been added to <https://github.com/hackclub/awesome-hackclub|awesome-hackclub>!!! :fastparrot:", project.Name), false, false), nil, nil),
	))

	if err != nil {
		logoru.Error(err)
	}
}

func SendDeniedMessage(project db.Project, reason string) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	_, _, err := client.PostMessage(project.UserID, slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Your project, *%s*, was denied :cry:\n*Reason*: %s", project.Name, reason), false, false), nil, nil),
	))

	if err != nil {
		logoru.Error(err)
	}
}
