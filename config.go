package main

import (
	"html/template"
	"time"
)

type Book struct {
	Params map[string]any

	WebLayoutDir string
	EpubLayoutDir string
	TextDir string
	OutputDir string

	BaseURL string
	Slug string
	Title string
	Status string
	LanguageCode string
	Tags []string
	Logline string
	Content string
	ContentHTML template.HTML
	Authors []Author
	PublisherName string
	CoverImageName string
	Mirrors []Address
	Series []BookSeries
	Chapters []Chapter
}

func (b Book) EpubBaseName() string {
	return b.Slug + ".epub"
}

type BookSeries struct {
	Name string
	Number float64
	URL string
}

type Chapter struct {
	Params map[string]any

	Book *Book
	Slug string
	Title string
	Description string
	ParsedDate time.Time
	ParsedLastModified time.Time
	Weight int
	Content string
	ContentHTML template.HTML

	Previous *Chapter
	Next *Chapter

	Date string
	LastMod string
}

type Author struct {
	Name      string
	NameSort  string
	Addresses []Address
}

type Address struct {
	Name    string
	Address string
	IsURL   bool
	Alternate []Address
}
