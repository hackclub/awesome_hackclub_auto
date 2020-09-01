package block_kit

import (
	"encoding/json"

	"github.com/Matt-Gleich/logoru"
	"github.com/slack-go/slack"
)

func SubmitModal(privateMetadata interface{}, action SlackActionID) slack.ModalViewRequest {
	privateMetadataJson, err := json.Marshal(privateMetadata)
	if err != nil {
		logoru.Error(err)
	}
	return slack.ModalViewRequest{
		CallbackID:      "submit",
		Type:            "modal",
		Title:           slack.NewTextBlockObject("plain_text", "Submit", false, false),
		PrivateMetadata: string(privateMetadataJson),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Project GitHub URL", false, false),
					BlockID: "url",
					Element: slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "url",
						InitialValue: action.GitHubURL,
					},
				},
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Project Name", false, false),
					BlockID: "name",
					Element: slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "name",
						InitialValue: action.ProjectName,
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
