package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Matt-Gleich/logoru"
	"github.com/hackclub/awesome_hackclub_auto/pkg/db"
	"github.com/hackclub/awesome_hackclub_auto/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logoru.Error(err)
	}

	if !util.VerifySlackRequest(r, buf.Bytes()) {
		logoru.Warning("invalid Slack request")
		_, err = w.Write(nil)
		if err != nil {
			logoru.Error(err)
		}
		return
	}

	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logoru.Error(e)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logoru.Error(err)
			os.Exit(1)
		}
		w.Header().Set("Content-Type", "text")
		_, err = w.Write([]byte(r.Challenge))
		if err != nil {
			logoru.Error(err)
		}
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.ReactionAddedEvent:
			if ev.Reaction == "awesome" && ev.User == ev.ItemUser {
				client := slack.New(os.Getenv("SLACK_TOKEN"))
				//client.PostMessage(ev.ItemUser, slack.MsgOptionText("You reacted with the sacred emoji", false))
				resp, err := client.GetConversationHistory(&slack.GetConversationHistoryParameters{
					ChannelID: ev.Item.Channel,
					Limit:     1,
					Inclusive: true,
					Latest:    ev.Item.Timestamp,
					Oldest:    ev.Item.Timestamp,
				})
				if err != nil {
					logoru.Error(err)
				} else if len(resp.Messages) >= 1 {
					logoru.Debug(resp.Messages[0].Text)

					projectIntent := util.GenerateProjectIntent(resp.Messages[0].Text)

					projectIntent.UserID = ev.User
					projectIntent.Timestamp = ev.Item.Timestamp
					projectIntent.Channel = ev.Item.Channel

					permalink, err := client.GetPermalink(&slack.PermalinkParameters{
						Channel: ev.Item.Channel,
						Ts:      ev.Item.Timestamp,
					})
					if err != nil {
						logoru.Error(err)
					}

					intentID := db.CreateProjectIntent(projectIntent)

					_, _, err = client.PostMessage(ev.ItemUser, slack.MsgOptionBlocks(
						slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Howdy! :wave: I see that you've reacted to <%s|one of your messages> with :awesome:. You're halfway to getting your project on <https://github.com/hackclub/awesome-hackclub|awesome-hackclub>! :sunglasses: Just click that button down there :arrow_down: to fill in some info and finish your submission! :tada:", permalink), false, false), nil, nil),
						slack.NewActionBlock(
							"",
							slack.NewButtonBlockElement(
								"submit",
								intentID,
								slack.NewTextBlockObject("plain_text", "Submit", true, false),
							),
						),
						slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "Was this a mistake? No worries! just ignore this message and you'll be fine.", false, false)),
					), slack.MsgOptionText("You're halfway to getting your project on <https://github.com/hackclub/awesome-hackclub|awesome-hackclub>! :sunglasses:", false), slack.MsgOptionDisableLinkUnfurl())
					if err != nil {
						logoru.Error(err)
					}
				}
			}
		}
	}
}
