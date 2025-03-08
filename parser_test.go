package main

import (
	"testing"
)

func TestUnmarshalCollectionBasic(t *testing.T) {
	collection := Collection{}
	err := Unmarshal([]byte(`title:`), &collection)
	if err != nil {
		t.Fatalf(`Unmarshal(...) = %v, want nil, error`, err)
	}

	{
		got := collection.Title
		want := ""
		if got != want {
			t.Fatalf(`collection.Title = %v, want %v, error`, got, want)
		}
	}

	{
		got := collection.LanguageCode
		want := LanguageCodeDefault
		if got != want {
			t.Fatalf(`collection.LanguageCode = %v, want %v, error`, got, want)
		}
	}

	{
		got := collection.LayoutsDir
		want := LayoutsDirDefault
		if got != want {
			t.Fatalf(`collection.LayoutsDir = %v, want %v, error`, got, want)
		}
	}

	{
		got := collection.OutputDir
		want := OutputDirDefault
		if got != want {
			t.Fatalf(`collection.OutputDir = %v, want %v, error`, got, want)
		}
	}

	{
		got := collection.Params["title"]
		want := collection.Title
		if got == want {
			t.Fatalf(`collection.Params["title"] = %v, want %v, error`, got, want)
		}
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
