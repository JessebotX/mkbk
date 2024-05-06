package mkbk

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const (
	IndexTemplateName   = "index.html"
	ChapterTemplateName = "_chapter.html"
)

func RenderBookToHTMLSite(inputDir, outputDir string, book *Book) error {
	layoutDir := filepath.Join(inputDir, "layout")

	// make root output directory
	err := os.MkdirAll(outputDir, os.ModePerm)
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

	indexFile, err := os.Create(filepath.Join(outputDir, IndexTemplateName))
	if err != nil {
		return err
	}

	err = copyDirectoryToOutput(layoutDir, outputDir, []string{IndexTemplateName, ChapterTemplateName, "README.md"})
	if err != nil {
		return err
	}

	// create index
	err = indexTemplate.Execute(indexFile, book)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// chapters
	for _, chapter := range book.Chapters {
		chapterFile, err := os.Create(filepath.Join(outputDir, chapter.Slug+".html"))
		if err != nil {
			return err
		}
		defer chapterFile.Close()

		err = chapterTemplate.Execute(chapterFile, &chapter)
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
