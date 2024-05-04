package mkbk

import (
	"html/template"
	"os"
	"path/filepath"
)

func RenderBookToHTMLSite(inputDir, outputDir string, book *Book) error {
	layoutDir := filepath.Join(inputDir, "layout")

	// make root output directory
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	// read layout configurations
	indexTemplate, err := template.ParseFiles(filepath.Join(layoutDir, "index.html"))
	if err != nil {
		return err
	}

	chapterTemplate, err := template.ParseFiles(filepath.Join(layoutDir, "_chapter.html"))
	if err != nil {
		return err
	}

	indexFile, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}

	_ = indexTemplate
	_ = chapterTemplate

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
