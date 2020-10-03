package gh

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Matt-Gleich/logoru"
)

func LoadPrivateKey() []byte {
	key, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(os.Getenv("GH_PRIVATE_KEY"))))
	if err != nil {
		logoru.Critical("Error reading GitHub private key")
		return nil
	}
	return key
}
