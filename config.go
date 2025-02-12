package mkbk

import (
	"time"
)

type Collection struct {
	Params []any
	Books []Book
	LayoutsDir string
	OutputDir string

	Title string
	BaseURL string
	LanguageCode string
}

type Book struct {
	Params []any
	Chapters []Chapter
	LastModified time.Time
	DatePublished time.Time

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
}

type Chapter struct {
	Params []any

	Title string
	Description string
	Content string
	Weight uint8
	LastModified time.Time
	DatePublished time.Time
}
