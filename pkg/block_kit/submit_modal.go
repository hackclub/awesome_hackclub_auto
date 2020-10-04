package block_kit

import (
	"github.com/hackclub/awesome_hackclub_auto/pkg/config"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/slack-go/slack"
)

func SubmitModal(metadata string, project db.ProjectFields) slack.ModalViewRequest {
	// Category stuff

	categoryOptions := []*slack.OptionBlockObject{}

	for _, category := range config.Categories {
		categoryOptions = append(categoryOptions, &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", category, false, false),
			Value: category,
		})
	}

	var initialCategory *slack.OptionBlockObject = nil
	if project.Category != "" {
		initialCategory = &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", project.Category, false, false),
			Value: project.Category,
		}
	}

	// Language stuff

	languageOptions := []*slack.OptionBlockObject{}

	for _, language := range config.Languages {
		languageOptions = append(languageOptions, &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", language, false, false),
			Value: language,
		})
	}

	var initialLanguage *slack.OptionBlockObject = nil
	if project.Language != "" {
		initialLanguage = &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", project.Language, false, false),
			Value: project.Language,
		}
	}

	return slack.ModalViewRequest{
		CallbackID:      "submit",
		Type:            "modal",
		Title:           slack.NewTextBlockObject("plain_text", "Submit Project", false, false),
		PrivateMetadata: metadata,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "GitHub URL", false, false),
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
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Your GitHub username", false, false),
					BlockID: "username",
					Element: slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "username",
						InitialValue: project.Username,
					},
				},
				slack.InputBlock{
					Optional: true,
					Type:     "input",
					Label:    slack.NewTextBlockObject("plain_text", "Project Description", false, false),
					BlockID:  "description",
					Element: slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "description",
						InitialValue: project.Description,
						Multiline:    true,
					},
				},
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Category", false, false),
					BlockID: "category",
					Element: slack.SelectBlockElement{
						Type:          "static_select",
						Placeholder:   slack.NewTextBlockObject("plain_text", "Select one...", false, false),
						ActionID:      "category",
						Options:       categoryOptions,
						InitialOption: initialCategory,
					},
				},
				slack.InputBlock{
					Type:    "input",
					Label:   slack.NewTextBlockObject("plain_text", "Language", false, false),
					BlockID: "language",
					Element: slack.SelectBlockElement{
						Type:          "static_select",
						Placeholder:   slack.NewTextBlockObject("plain_text", "Select one...", false, false),
						ActionID:      "language",
						Options:       languageOptions,
						InitialOption: initialLanguage,
					},
				},
			},
		},
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
	}
}
