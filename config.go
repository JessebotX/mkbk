package mkbk

import (
	"time"
)

type Book struct {
	Params        map[string]any
	Title         string
	TitleSort     string
	Status        string
	LanguageCode  string
	Tags          []string
	Logline       string
	Content       string
	Authors       []Author
	Publisher     string
	CoverPath     string
	Chapters      []Chapter
	DatePublished time.Time
	LastModified  time.Time
	SeriesName    string
	SeriesNumber  float32
	IDs           []string
	Mirrors       []Address
}

type Chapter struct {
	Params        map[string]any
	ParentBook    *Book
	ParentChapter *Chapter
	Title         string
	Description   string
	DatePublished time.Time
	LastModified  time.Time
	Weight        int
	Subchapters   []Chapter
}

type Author struct {
	Params    map[string]any
	Name      string
	NameSort  string
	Bio       string
	Addresses []Address
	ImagePath string
}

type Address struct {
	Params  map[string]any
	Name    string
	Address string
	IsURL   bool
}
