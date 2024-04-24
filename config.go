package mkbk

import (
	"time"
)

type BookStatus string

const (
	STATUS_COMPLETED BookStatus = "completed"
	STATUS_ONGOING   BookStatus = "ongoing"
	STATUS_HIATUS    BookStatus = "hiatus"
)

type Book struct {
	Params        map[string]any
	Title         string
	TitleSort     string
	Status        BookStatus
	LanguageCode  string
	Tags          []string
	Logline       string
	Content       string
	Authors       []Author
	Publisher     Publisher
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
	ParentBook    *Book
	ParentChapter *Chapter
	Title         string
	Description   string
	DatePublished time.Time
	LastModified  time.Time
	Weight        int
	Subchapters   []Chapter
	Params        map[string]any
}

type Publisher struct {
	Name      string
	Addresses []Address
	Bio       string
	ImagePath string
	Params    map[string]any
}

type Author struct {
	Name      string
	NameSort  string
	Bio       string
	Addresses []Address
	ImagePath string
	Params    map[string]any
}

type Address struct {
	Name    string
	Address string
	IsURL   bool
	Params  map[string]any
}
