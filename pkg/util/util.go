package util

import (
	"fmt"
	"os"

	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/slack-go/slack"
)

func AddProjectToQueue(project db.Project) {
	db.AddProjectToQueue(project)
	client := slack.New(os.Getenv("SLACK_TOKEN"))
	client.PostMessage(project.UserID, slack.MsgOptionText(fmt.Sprintf("Your project, *%s*, has successfully been submitted! :tada:", project.Name), false))

	client.PostMessage(os.Getenv("REVIEW_CHANNEL"), slack.MsgOptionBlocks(
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s> just submitted a project for review: <%s|%s>", project.UserID, project.GitHubURL, project.Name), false, false), nil, nil),
	))
}
