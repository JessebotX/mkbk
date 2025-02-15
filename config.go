package mkbk

import (
	"time"
)

type Collection struct {
	Params map[any]any
	Books []Book
	LayoutsDir string
	OutputDir string

	Title string
	BaseURL string
	LanguageCode string
}

type Book struct {
	Params map[any]any
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
	Params map[any]any

	Title string
	Description string
	Content string
	Weight uint8
	LastModified time.Time
	DatePublished time.Time
}
