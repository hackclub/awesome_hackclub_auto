package block_kit

import "github.com/slack-go/slack"

func AlreadyInQueue() slack.ModalViewRequest {
	return slack.ModalViewRequest{
		Type: "modal",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "This project is already in the review queue! You'll get a DM once it gets accepted.", false, false), nil, nil),
			},
		},
		Title: slack.NewTextBlockObject("plain_text", "Submit", false, false),
	}
}
