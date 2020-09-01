package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Matt-Gleich/logoru"
	"github.com/slack-go/slack"
)

func HandleInteractivity(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logoru.Error(err)
	}
	r.Form, err = url.ParseQuery(string(buf))
	if err != nil {
		logoru.Error(err)
	}

	parsed := slack.InteractionCallback{}

	err = json.Unmarshal([]byte(r.Form.Get("payload")), &parsed)
	if err != nil {
		logoru.Error(err)
	}

	switch parsed.Type {
	case slack.InteractionTypeBlockActions:
		logoru.Debug(parsed.ActionCallback.BlockActions[0].ActionID)
	}
}
