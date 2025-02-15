package mkbk

import (
	"time"
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

	Params map[string]any
	ChaptersDir string
	Chapters []Chapter
	LastModified time.Time
	DatePublished time.Time
}

type Chapter struct {
	Title string
	Description string
	Content string
	Weight uint8
	LastModified time.Time
	DatePublished time.Time

	Params map[string]any
}
