package main

import (
	"html/template"
	"time"
)

type Book struct {
	Params       map[string]any
	LayoutDir    string
	OutputDir    string
	ChaptersDir  string
	Slug         string
	Title        string
	Status       string
	LanguageCode string
	Tags         []string
	Logline      string
	Content      string
	Authors      []Author
	Publisher    string
	CoverPath    string
	Chapters     []Chapter
	SeriesName   string
	SeriesNumber float32
	IDs          []string
	Mirrors      []Address
}

func (b Book) EpubBaseName() string {
	return b.Slug + ".epub"
}

type Chapter struct {
	Params             map[string]any
	Parent             *Book
	ParentSectionTitle string
	Slug               string
	Title              string
	Description        string
	Date               time.Time
	LastModified       time.Time
	Weight             int
	ContentHTML        template.HTML
	Previous           *Chapter
	Next               *Chapter
}

type Author struct {
	Name      string
	NameSort  string
	Bio       string
	Addresses []Address
	ImagePath string
}

type Address struct {
	Name    string
	Address string
	IsURL   bool
}
