package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/yuin/goldmark-meta"

	"gopkg.in/yaml.v3"
)

const (
	DateLayoutString = "2006-01-02T03:04:05"
)

func Unmarshal(data []byte, collection *Collection) error {
	err := yaml.Unmarshal(data, &collection.Params)
	if err != nil {
		return err
	}

	jsonbody, err := json.Marshal(collection.Params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonbody, &collection)
	if err != nil {
		return err
	}

	// Set default values
	if strings.TrimSpace(collection.LayoutsDir) == "" {
		collection.LayoutsDir = LayoutsDirDefault
	}

	if strings.TrimSpace(collection.OutputDir) == "" {
		collection.OutputDir = OutputDirDefault
	}

	if strings.TrimSpace(collection.LanguageCode) == "" {
		collection.LanguageCode = LanguageCodeDefault
	}

	return nil
}

func UnmarshalBook(data []byte, book *Book, collection *Collection) error {
	err := yaml.Unmarshal(data, &book.Params)
	if err != nil {
		return err
	}

	jsonbody, err := json.Marshal(book.Params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonbody, &book)
	if err != nil {
		return err
	}

	/// set default values
	// Override book's BaseURL if collection's BaseURL value exists
	if collection != nil && strings.TrimSpace(collection.BaseURL) != "" {
		book.BaseURL = collection.BaseURL
	}

	// Inherit collection's LanguageCode, else set to default value
	if strings.TrimSpace(book.LanguageCode) != "" {
		if strings.TrimSpace(collection.LanguageCode) != "" {
			book.LanguageCode = collection.LanguageCode
		} else {
			book.LanguageCode = LanguageCodeDefault
		}
	}

	if strings.TrimSpace(book.ChaptersDir) == "" {
		book.ChaptersDir = filepath.Join(ChaptersDirDefault)
	}

	if strings.TrimSpace(collection.LayoutsDir) == "" {
		collection.LayoutsDir = filepath.Join(LayoutsDirDefault)
	}

	if strings.TrimSpace(collection.OutputDir) == "" {
		collection.OutputDir = filepath.Join(OutputDirDefault)
	}

	if strings.TrimSpace(book.Status) == "" {
		book.Status = BookStatusDefault
	}

	// parse blurb
	contentHTML, _, err := convertMarkdownToHTML([]byte(book.Content))
	if err != nil {
		return err
	}
	book.ContentHTML = contentHTML

	// read chapters
	fullChaptersDir := filepath.Join(book.BookDir, book.ChaptersDir)
	chapterFiles, err := os.ReadDir(fullChaptersDir)
	if err != nil {
		return err
	}

	chapters := make([]Chapter, 0)
	for _, file := range chapterFiles {
		// ignore hidden files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		chapterFilePath := filepath.Join(fullChaptersDir, file.Name())
		chapter, err := parseChapter(chapterFilePath)
		if err != nil {
			return err
		}
		chapter.Parent = book
		chapters = append(chapters, chapter)
	}

	// sort and set next and previous chapters
	slices.SortFunc(chapters, func(a, b Chapter) int {
		// TODO also compare date, then title
		return a.Weight - b.Weight
	})

	for i, _ := range chapters {
		if i >= 1 {
			chapters[i].Previous = &chapters[i-1]
		}

		if i < (len(chapters) - 1) {
			chapters[i].Next = &chapters[i+1]
		}
	}

	book.Chapters = chapters
	book.Parent = collection

	return nil
}

func parseChapter(path string) (Chapter, error) {
	chapter := Chapter{
		ID: strings.TrimSuffix(filepath.Base(path), ".md"),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return chapter, err
	}

	html, metadata, err := convertMarkdownToHTML(data)
	if err != nil {
		return chapter, err
	}

	chapter.Content = string(data)
	chapter.ContentHTML = html
	chapter.Params = metadata

	// set title (default: the base file name)
	if metadata["title"] != nil {
		switch v := metadata["title"].(type) {
		case string:
			chapter.Title = v
		default:
			return chapter, fmt.Errorf("%s chapter title is of the wrong type (expected string)", path)
		}
	} else {
		chapter.Title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	// set description (default: the chapter title)
	if metadata["description"] != nil {
		switch v := metadata["description"].(type) {
		case string:
			chapter.Description = v
		default:
			return chapter, fmt.Errorf("%s chapter description is of the wrong type (expected string)", path)
		}
	} else {
		chapter.Description = "Read " + chapter.Title
	}

	// set chapter weight if provided (default: 1)
	if metadata["weight"] != nil {
		switch v := metadata["weight"].(type) {
		case int:
			chapter.Weight = v
		default:
			return chapter, fmt.Errorf("%s chapter weight is of the wrong type (expected int)", path)
		}
	} else {
		chapter.Weight = 1
	}

	// set date if provided
	if metadata["date"] != nil {
		switch v := metadata["date"].(type) {
		case string:
			dateString := v
			date, err := time.Parse(DateLayoutString, dateString)
			if err != nil {
				return chapter, err
			}

			chapter.DatePublishedParsed = date
		default:
			return chapter, fmt.Errorf("%s chapter date is of the wrong type (expected string)", path)
		}
	}

	// set date if provided
	if metadata["last_modified"] != nil {
		switch v := metadata["last_modified"].(type) {
		case string:
			dateString := v
			date, err := time.Parse(DateLayoutString, dateString)
			if err != nil {
				return chapter, err
			}

			chapter.LastModifiedParsed = date
		default:
			return chapter, fmt.Errorf("%s chapter last_modified is of the wrong type (expected string)", path)
		}
	}

	return chapter, nil
}

func convertMarkdownToHTML(content []byte) (template.HTML, map[string]any, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			extension.GFM,
			extension.Footnote,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(
					extension.TypographicSubstitutions{
						// replace -- with an em-dash, ignore en-dashes
						extension.EnDash: []byte("&mdash;"),
						extension.EmDash: nil,
					},
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(),
	)

	var buffer bytes.Buffer
	context := parser.NewContext()
	err := md.Convert(content, &buffer, parser.WithContext(context))
	if err != nil {
		return template.HTML(""), nil, err
	}

	metadata := meta.Get(context)

	return template.HTML(buffer.String()), metadata, nil
}
