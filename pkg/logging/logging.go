package logging

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"crypto/rand"

	"github.com/gleich/logoru"
	"github.com/honeybadger-io/honeybadger-go"
)

func Log(v interface{}, level string, localOnly bool) {
	// Only log to Honeybadger if an API key is set and it's an error
	if os.Getenv("HONEYBADGER_API_KEY") != "" && !localOnly && (level == "error" || level == "critical" || level == "warning") {
		logoru.Info("Logging an error to Honeybadger...")
		_, err := honeybadger.Notify(fmt.Sprintf("%s: %v", strings.ToUpper(level), v), honeybadger.Fingerprint{Content: genFingerprint()})
		if err != nil {
			logoru.Error("Ironically enough, something went wrong while trying to log an error :(")
		}
	}

	switch level {
	case "debug":
		logoru.Debug(v)
	case "info":
		logoru.Info(v)
	case "warning":
		logoru.Warning(v)
	case "error":
		logoru.Error(v)
	case "critical":
		logoru.Critical(v)
	case "success":
		logoru.Success(v)
	}
}

func genFingerprint() string {
	fingerprint := make([]byte, 16)
	_, err := rand.Read(fingerprint)
	if err != nil {
		_, err := honeybadger.Notify("Error generating Honeybadger fingerprint")
		if err != nil {
			logoru.Error("Ironically enough, something went wrong while trying to log an error :(")
		}
		return ""
	}
	return hex.EncodeToString(fingerprint)
}
