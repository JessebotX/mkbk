package mkbk

import (
	"testing"
)

func TestUnmarshalCollectionBasic(t *testing.T) {
	collection := Collection{}
	err := Unmarshal([]byte(`title:`), &collection)
	if err != nil {
		t.Fatalf(`Unmarshal(...) = %v, want nil, error`, err)
	}
	if collection.Title != "" {
		t.Fatalf(`collection.Title = %v, want %v, error`, collection.Title, "")
	}
	if collection.LanguageCode != LanguageCodeDefault {
		t.Fatalf(`collection.LanguageCode = %v, want %v, error`, collection.LanguageCode, LanguageCodeDefault)
	}
	if collection.LayoutsDir != LayoutsDirDefault {
		t.Fatalf(`collection.LayoutsDir = %v, want %v, error`, collection.LayoutsDir, LayoutsDirDefault)
	}
	if collection.OutputDir != OutputDirDefault {
		t.Fatalf(`collection.OutputDir = %v, want %v, error`, collection.OutputDir, OutputDirDefault)
	}
	if collection.Params["title"] == collection.Title {
		t.Fatalf(`collection.Params["title"] = %v, want %v, error`, collection.Params["title"], collection.Title)
	}
}

