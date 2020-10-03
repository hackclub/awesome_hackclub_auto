package util

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Matt-Gleich/logoru"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/slack-go/slack"
)

func SendReviewMessage(project db.Project) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	description := project.Fields.Description
	if description == "" {
		description = "_<no description>_"
	}

	_, _, err := client.PostMessage(os.Getenv("REVIEW_CHANNEL"), slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s> just submitted a project for review: <%s|%s>", project.Fields.UserID, project.Fields.GitHubURL, project.Fields.Name), false, false), []*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Description*\n%s", description), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Category*\n%s", project.Fields.Category), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Language*\n%s", project.Fields.Language), false, false),
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

	_, _, err := client.PostMessage(project.Fields.UserID, slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("WOO-HOO!! :tada: Your project, *%s*, has been added to <https://github.com/hackclub/awesome-hackclub|awesome-hackclub>!!! :fastparrot:", project.Fields.Name), false, false), nil, nil),
	))

	if err != nil {
		logoru.Error(err)
	}
	logoru.Info("Sent approved message for", project.Fields.Name)
}

func SendDeniedMessage(project db.Project, reason string) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	_, _, err := client.PostMessage(project.Fields.UserID, slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Your project, *%s*, was denied :cry:\n*Reason*: %s", project.Fields.Name, reason), false, false), nil, nil),
	))

	if err != nil {
		logoru.Error(err)
	}
	logoru.Info("Sent denied message for", project.Fields.Name)
}

// GenerateProjectIntent generates some pre-filled project data, given the text of the message
func GenerateProjectIntent(messageText string) db.ProjectFields {
	re := regexp.MustCompile(`https?:\/\/github\.com\/([^\/>\|]+)\/([^\/>\|]+)`).FindStringSubmatch(messageText)

	if re != nil {
		// TODO: automatically pre-fill repo language and description
		return db.ProjectFields{
			GitHubURL: re[0],
			Name:      re[2],
		}
	} else {
		// There aren't any GitHub URLs in the message
		return db.ProjectFields{}
	}
}

func IsValidProjectURL(url string) bool {
	return regexp.MustCompile(`^https?:\/\/github\.com\/([^\/>\|]+)\/([^\/>\|]+)\/?$`).MatchString(url)
}
