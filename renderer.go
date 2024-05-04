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

	chapterTemplate, err := template.ParseFiles(filepath.Join(layoutDir, "_chapter", "index.html"))
	if err != nil {
		return err
	}

	_ = indexTemplate
	_ = chapterTemplate

	return nil
}
