package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/yuin/goldmark-meta"

	"gopkg.in/yaml.v3"
)

func UnmarshalBookConfigFile(path string, book *Book) error {
	if book == nil {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return UnmarshalBookConfigData(filepath.Dir(path), data, book)
}

func UnmarshalBookConfigData(dir string, data []byte, book *Book) error {
	if data == nil {
		return nil
	}

	if book == nil {
		return nil
	}

	err := yaml.Unmarshal(data, &book.Params)
	if err != nil {
		return err
	}

	// apparently, json unmarshalling/marshalling fixes problems with
	// yaml marshalling/unmarshalling into things like camel case
	// (e.g. coverPath == CoverPath == coverpath, etc.)
	jsonbody, err := json.Marshal(book.Params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonbody, &book)
	if err != nil {
		return err
	}

	// default directory configurations
	if strings.TrimSpace(book.OutputDir) == "" {
		book.OutputDir = "out"
	}

	if strings.TrimSpace(book.LayoutDir) == "" {
		book.LayoutDir = "layout"
	}

	if strings.TrimSpace(book.ChaptersDir) == "" {
		book.ChaptersDir = "src"
	}

	if strings.TrimSpace(book.Slug) == "" {
		sanitizer := regexp.MustCompile("([^a-zA-Z0-9]+)")
		book.Slug = strings.ToLower(sanitizer.ReplaceAllString(book.Title, "-"))
	}

	book.Chapters, err = readChaptersDir(filepath.Join(dir, book.ChaptersDir), book)
	if err != nil {
		return err
	}

	return nil
}

func readChaptersDir(dir string, book *Book) ([]Chapter, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(),
	)

	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	chapters := make([]Chapter, 0)
	for _, item := range items {
		chapterPath := filepath.Join(dir, item.Name())

		chapter := Chapter{
			Parent: book,
			Slug:   strings.TrimSuffix(filepath.Base(chapterPath), ".md"),
		}

		err = unmarshalChapter(chapterPath, md, &chapter)
		if err != nil {
			return chapters, err
		}

		chapters = append(chapters, chapter)
	}

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

	return chapters, nil
}

func unmarshalChapter(path string, md goldmark.Markdown, chapter *Chapter) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	context := parser.NewContext()
	err = md.Convert(content, &buffer, parser.WithContext(context))
	if err != nil {
		return err
	}
	metadata := meta.Get(context)
	chapter.Params = metadata

	// set title (default: the base file name)
	if metadata["title"] != nil {
		switch v := metadata["title"].(type) {
		case string:
			chapter.Title = v
		default:
			return fmt.Errorf("%s chapter title is of the wrong type (expected string)", path)
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
			return fmt.Errorf("%s chapter description is of the wrong type (expected string)", path)
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
			return fmt.Errorf("%s chapter weight is of the wrong type (expected int)", path)
		}
	} else {
		chapter.Weight = 1
	}

	// set date if provided
	if metadata["date"] != nil {
		switch v := metadata["date"].(type) {
		case string:
			dateString := v
			date, err := time.Parse("2006-01-02T03:04:05", dateString)
			if err != nil {
				return err
			}

			chapter.Date = date
		default:
			return fmt.Errorf("%s chapter date is of the wrong type (expected string)", path)
		}
	}

	// set date if provided
	if metadata["last_modified"] != nil {
		switch v := metadata["last_modified"].(type) {
		case string:
			dateString := v
			date, err := time.Parse("2006-01-02T03:04:05", dateString)
			if err != nil {
				return err
			}

			chapter.LastModified = date
		default:
			return fmt.Errorf("%s chapter last_modified is of the wrong type (expected string)", path)
		}
	}
	chapter.ContentHTML = template.HTML(buffer.String())

	return nil
}
