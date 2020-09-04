package block_kit

import (
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/slack-go/slack"
)

func SubmitModal(intentID string, project db.Project) slack.ModalViewRequest {
	return slack.ModalViewRequest{
		CallbackID:      "submit",
		Type:            "modal",
		Title:           slack.NewTextBlockObject("plain_text", "Submit", false, false),
		PrivateMetadata: intentID,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Project GitHub URL", false, false),
					BlockID: "url",
					Element: slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "url",
						InitialValue: project.GitHubURL,
					},
				},
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Project Name", false, false),
					BlockID: "name",
					Element: slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "name",
						InitialValue: project.Name,
					},
				},
				slack.InputBlock{
					Optional: true,
					Type:     "input",
					Label:    slack.NewTextBlockObject("plain_text", "Project Description", false, false),
					BlockID:  "description",
					Element: slack.PlainTextInputBlockElement{
						Type:      "plain_text_input",
						ActionID:  "description",
						Multiline: true,
					},
				},
			},
		},
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
	}
}
