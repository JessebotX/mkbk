package mkbk

import (
	"time"
)

const (
	LanguageCodeDefault = "en"
	BookStatusDefault = "Completed"
	LayoutsDirDefault = "./layouts"
	OutputDirDefault = "./out"
)

type Collection struct {
	Title string
	BaseURL string
	BookDirs []string
	LanguageCode string

	Params map[string]any
	Books []Book
	LayoutsDir string
	OutputDir string
}

type Book struct {
	Title string
	TitleSort string
	Authors string
	AuthorsSort string
	BaseURL string
	Description string
	LanguageCode string
	Content string
	Status string
	CoverImageName string
	DatePublished string

	LayoutsDir string
	OutputDir string
	Params map[string]any
	ChaptersDir string
	Chapters []Chapter
	LastModifiedParsed time.Time
	DatePublishedParsed time.Time
}

type Chapter struct {
	Title string
	Description string
	Content string
	Weight int
	DatePublished string
	LastModified string

	DatePublishedParsed time.Time
	LastModifiedParsed time.Time
	Params map[string]any
}

