package mkbk

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Unmarshal(data []byte, collection *Collection) error {
	err := yaml.Unmarshal(data, &collection.Params)
	if err != nil {
		return err
	}

	jsonbody, err := json.Marshal(collection.Params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonbody, &collection)
	if err != nil {
		return err
	}

	// Set default values
	if (strings.TrimSpace(collection.LayoutsDir) == "") {
		collection.LayoutsDir = filepath.Join("./layouts")
	}

	if (strings.TrimSpace(collection.OutputDir) == "") {
		collection.OutputDir = filepath.Join("./out")
	}

	if (strings.TrimSpace(collection.LanguageCode) == "") {
		collection.LanguageCode = "en"
	}

	return nil
}

func UnmarshalBook(data []byte, book *Book) error {
	err := yaml.Unmarshal(data, &book.Params)
	if err != nil {
		return err
	}

	jsonbody, err := json.Marshal(book.Params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonbody, &book)
	if err != nil {
		return err
	}

	return nil
}
