package main

import (
	"errors"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	epub "github.com/go-shiori/go-epub"
)

const (
	IndexTemplateName   = "index.html"
	ChapterTemplateName = "_chapter.html"
	RSSFeedTemplateName = "_rss.xml" // TODO implement RSS feeds
	ImagesFolderName    = "images"
	CSSFolderName       = "css"
)

func RenderBookToHTMLSite(inputDir, outputDir string, book *Book) error {
	webLayoutDir := filepath.Join(inputDir, book.WebLayoutDir)
	imagesDir := filepath.Join(inputDir, ImagesFolderName)
	outputImagesDir := filepath.Join(outputDir, ImagesFolderName)

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
	indexTemplate, err := template.ParseFiles(filepath.Join(webLayoutDir, IndexTemplateName))
	if err != nil {
		return err
	}

	chapterTemplate, err := template.ParseFiles(filepath.Join(webLayoutDir, ChapterTemplateName))
	if err != nil {
		return err
	}

	// copy other layout files (theme-specific stylesheets, images, etc.)
	err = copyDirectoryToOutput(webLayoutDir, outputDir, []string{IndexTemplateName, ChapterTemplateName, "README.md"}, true)
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

	// create folders
	err = os.MkdirAll(outputImagesDir, os.ModePerm)
	if err != nil {
		return err
	}

	// begin writing book epub file
	// TODO support images in epub
	bookEpub, err := epub.NewEpub(book.Title)
	if err != nil {
		return err
	}
	// add images into epub
	imagesDirItems, err := os.ReadDir(imagesDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		return err
	}
	for _, image := range imagesDirItems {
		imagePath, err := bookEpub.AddImage(filepath.Join(imagesDir, image.Name()), image.Name())
		if err != nil {
			return err
		}

		if image.Name() == book.CoverImageName {
			err = bookEpub.SetCover(imagePath, "")
			if err != nil {
				return err
			}
		}
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

	coverName := book.CoverImageName
	if strings.TrimSpace(coverName) != "" {
		coverPath := filepath.Join(imagesDir, coverName)

		err = os.Link(coverPath, filepath.Join(outputImagesDir, coverName))
		if err != nil {
			return err
		}

		/*
			dir := filepath.Dir(coverName)

			err = os.MkdirAll(filepath.Join(outputDir, dir), os.ModePerm)
			if err != nil {
				return err
			}

			err = os.Link(coverName, filepath.Join(outputDir, coverName))
			if err != nil {
				return err
			}*/
	}

	return nil
}

func copyDirectoryToOutput(inputDir, outputDir string, excludes []string, copySubdirs bool) error {
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

		if item.IsDir() && copySubdirs {
			err = os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return err
			}

			return copyDirectoryToOutput(fullPath, outputPath, []string{}, copySubdirs)
		} else if item.IsDir() && !copySubdirs {
			continue
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
