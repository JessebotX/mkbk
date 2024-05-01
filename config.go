package mkbk

import (
	"time"
)

type Book struct {
	Params       map[string]any
	Title        string
	TitleSort    string
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

type Chapter struct {
	Params             map[string]any
	ParentBook         *Book
	ParentSectionTitle string
	Slug               string
	Title              string
	Description        string
	DatePublished      time.Time
	LastModified       time.Time
	Weight             int
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
