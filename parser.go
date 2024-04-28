package mkbk

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

func UnmarshalBookConfigFile(filepath string, book *Book) error {
	if book == nil {
		return nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return UnmarshalBookConfigData(data, book)
}

func UnmarshalBookConfigData(data []byte, book *Book) error {
	if data == nil {
		return nil
	}

	if book == nil {
		return nil
	}

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
