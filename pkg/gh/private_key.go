package gh

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
)

func LoadPrivateKey() []byte {
	key, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(os.Getenv("GH_PRIVATE_KEY"))))
	if err != nil {
		logging.Log("Error reading GitHub private key", "critical", false)
		return nil
	}
	return key
}
