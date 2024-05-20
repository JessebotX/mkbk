package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	epub "github.com/go-shiori/go-epub"
)

const (
	IndexTemplateName   = "index.html"
	ChapterTemplateName = "_chapter.html"
	RSSFeedTemplateName = "_rss.xml" // TODO implement RSS feeds
)

func RenderBookToHTMLSite(inputDir, outputDir string, book *Book) error {
	layoutDir := filepath.Join(inputDir, book.LayoutDir)

	// clean up any existing output dir
	err := os.RemoveAll(outputDir)
	if err != nil {
		return err
	}

	// make root output directory
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	// read layout configurations
	indexTemplate, err := template.ParseFiles(filepath.Join(layoutDir, IndexTemplateName))
	if err != nil {
		return err
	}

	chapterTemplate, err := template.ParseFiles(filepath.Join(layoutDir, ChapterTemplateName))
	if err != nil {
		return err
	}

	// copy other layout files (theme-specific stylesheets, images, etc.)
	err = copyDirectoryToOutput(layoutDir, outputDir, []string{IndexTemplateName, ChapterTemplateName, "README.md"})
	if err != nil {
		return err
	}

	// create index file using index template
	indexFile, err := os.Create(filepath.Join(outputDir, IndexTemplateName))
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// write index index
	err = indexTemplate.Execute(indexFile, book)
	if err != nil {
		return err
	}

	// begin writing book epub file
	// TODO support images in epub
	bookEpub, err := epub.NewEpub(book.Title)
	if err != nil {
		return err
	}

	// parse chapters
	for _, chapter := range book.Chapters {
		// create chapter html file
		chapterFile, err := os.Create(filepath.Join(outputDir, chapter.Slug+".html"))
		if err != nil {
			return err
		}
		defer chapterFile.Close()

		// write chapter html
		err = chapterTemplate.Execute(chapterFile, &chapter)
		if err != nil {
			return err
		}

		// write to epub
		sectionContent := "<h1>" + chapter.Title + "</h1>" + string(chapter.ContentHTML)
		_, err = bookEpub.AddSection(sectionContent, chapter.Title, "", "")
		if err != nil {
			return err
		}
	}

	// create epub
	bookEpubPath := filepath.Join(outputDir, book.EpubBaseName())
	err = bookEpub.Write(bookEpubPath)
	if err != nil {
		return err
	}

	if strings.TrimSpace(book.CoverPath) != "" {
		dir := filepath.Dir(book.CoverPath)

		err = os.MkdirAll(filepath.Join(outputDir, dir), os.ModePerm)
		if err != nil {
			return err
		}

		err = os.Link(book.CoverPath, filepath.Join(outputDir, book.CoverPath))
		if err != nil {
			return err
		}
	}

	for _, author := range book.Authors {
		if strings.TrimSpace(author.ImagePath) == "" {
			continue
		}

		dir := filepath.Dir(author.ImagePath)

		err = os.MkdirAll(filepath.Join(outputDir, dir), os.ModePerm)
		if err != nil {
			return err
		}

		err = os.Link(author.ImagePath, filepath.Join(outputDir, author.ImagePath))
		if err != nil {
			return err
		}
	}

	return nil
}

func copyDirectoryToOutput(inputDir, outputDir string, excludes []string) error {
	items, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		excluded := false
		for _, exclude := range excludes {
			if strings.ToLower(exclude) == item.Name() {
				excluded = true
				break
			}
		}

		if excluded {
			continue
		}

		fullPath := filepath.Join(inputDir, item.Name())
		outputPath := filepath.Join(outputDir, item.Name())

		if item.IsDir() {
			err = os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return err
			}

			return copyDirectoryToOutput(fullPath, outputPath, []string{})
		}

		err = os.RemoveAll(outputPath)
		if err != nil {
			return err
		}

		err = os.Link(fullPath, outputPath)
		if err != nil {
			return err
		}
	}

	return nil
}
