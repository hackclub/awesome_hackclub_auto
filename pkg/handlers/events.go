package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Matt-Gleich/logoru"
	"github.com/slack-go/slack/slackevents"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logoru.Error(err)
	}
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
			logoru.Debug(fmt.Sprintf("Reaction added by user %s", ev.User))
		}
	}
}
