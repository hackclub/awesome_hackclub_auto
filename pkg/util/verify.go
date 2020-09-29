package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
)

func VerifySlackRequest(r *http.Request, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(os.Getenv("SLACK_SIGNING_SECRET")))

	body = append([]byte(r.Header.Get("X-Slack-Request-Timestamp")+":"), body...)
	body = append([]byte("v0:"), body...)

	_, err := mac.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal([]byte("v0="+hex.EncodeToString(mac.Sum(nil))), []byte(r.Header.Get("X-Slack-Signature")))
}
