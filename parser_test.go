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

func TestUnmarshalCollectionWithExtraParams(t *testing.T) {
	collection := Collection{}
	err := Unmarshal(
		[]byte(`title: random title
extratag: 2
1: true`),
		&collection)

	if err != nil {
		t.Fatalf(`Unmarshal(...) = %v, want nil, error`, err)
	}

	{
		got := collection.Title
		want := "random title"
		if got != want {
			t.Fatalf(`collection.Title = %v, want %v, error`, got, want)
		}
	}
	{
		got := collection.Params["extratag"].(int)
		want := 2
		if got != want {
			t.Fatalf(`collection.Params["extratag"].(int) = %v, want %v, error`, got, want)
		}
	}

	{
		got := collection.Params["1"].(bool)
		want := true
		if got != true {
			t.Fatalf(`collection.Params["1"].(bool) = %v, want %v, error`, got, want)
		}
	}
}

