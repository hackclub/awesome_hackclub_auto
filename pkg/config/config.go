package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Matt-Gleich/logoru"
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
		logoru.Error(err)
		os.Exit(1)
	}

	categories, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(categories, &Categories)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}
}

func populateLanguages() {
	languagesPath, _ := filepath.Abs("./pkg/config/languages.json")
	file, err := os.Open(languagesPath)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}

	languages, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(languages, &Languages)
	if err != nil {
		logoru.Error(err)
		os.Exit(1)
	}
}
