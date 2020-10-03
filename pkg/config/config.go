package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hackclub/awesome_hackclub_auto/pkg/logging"
)

var (
	Categories []string
	Languages  []string
)

func PopulateConfig() {
	populateCategories()
	populateLanguages()
}

func populateCategories() {
	categoriesPath, _ := filepath.Abs("./pkg/config/categories.json")
	file, err := os.Open(categoriesPath)
	if err != nil {
		logging.Log(err, "error", false)
	}

	categories, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(categories, &Categories)
	if err != nil {
		logging.Log(err, "error", false)
	}
}

func populateLanguages() {
	languagesPath, _ := filepath.Abs("./pkg/config/languages.json")
	file, err := os.Open(languagesPath)
	if err != nil {
		logging.Log(err, "error", false)
	}

	languages, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(languages, &Languages)
	if err != nil {
		logging.Log(err, "error", false)
	}
}
