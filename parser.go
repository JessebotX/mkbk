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
	if strings.TrimSpace(collection.LayoutsDir) == "" {
		collection.LayoutsDir = LayoutsDirDefault
	}

	if strings.TrimSpace(collection.OutputDir) == "" {
		collection.OutputDir = OutputDirDefault
	}

	if strings.TrimSpace(collection.LanguageCode) == "" {
		collection.LanguageCode = LanguageCodeDefault
	}

	return nil
}

func UnmarshalBook(data []byte, book *Book, collection *Collection) error {
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

	/// set default values
	// Override book's BaseURL if collection's BaseURL value exists
	if collection != nil && strings.TrimSpace(collection.BaseURL) != "" {
		book.BaseURL = collection.BaseURL
	}

	// Inherit collection's LanguageCode, else set to default value
	if strings.TrimSpace(book.LanguageCode) != "" {
		if strings.TrimSpace(collection.LanguageCode) != "" {
			book.LanguageCode = collection.LanguageCode
		} else {
			book.LanguageCode = LanguageCodeDefault
		}
	}

	if strings.TrimSpace(collection.LayoutsDir) == "" {
		collection.LayoutsDir = filepath.Join(LayoutsDirDefault)
	}

	if strings.TrimSpace(collection.OutputDir) == "" {
		collection.OutputDir = filepath.Join(OutputDirDefault)
	}

	if strings.TrimSpace(book.Status) == "" {
		book.Status = BookStatusDefault
	}

	return nil
}
