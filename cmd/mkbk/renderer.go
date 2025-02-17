package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/JessebotX/mkbk"
	"golang.org/x/sync/errgroup"
)

func RenderCollectionToHTML(workingDir string, collection mkbk.Collection) error {
	layoutsDir := filepath.Join(workingDir, collection.LayoutsDir)
	outputDir := filepath.Join(workingDir, collection.OutputDir)

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
	indexTemplatePath := filepath.Join(layoutsDir, IndexTemplateName)
	indexTemplate, err := template.ParseFiles(indexTemplatePath)
	if err != nil {
		return err
	}

	bookTemplatePath := filepath.Join(layoutsDir, BookFolderName, BookTemplateName)
	bookTemplate, err := template.ParseFiles(bookTemplatePath)
	if err != nil {
		return err
	}

	chapterTemplatePath := filepath.Join(layoutsDir, BookFolderName, ChapterTemplateName)
	chapterTemplate, err := template.ParseFiles(chapterTemplatePath)
	if err != nil {
		return err
	}

	// copy contents in layouts directory
	err = copyDirectoryToOutput(layoutsDir, outputDir, []string{IndexTemplateName, BookTemplateName, ChapterTemplateName, "README.md"}, true)
	if err != nil {
		return err
	}

	// create index file using index template
	indexFile, err := os.Create(filepath.Join(outputDir, CollectionOutputIndexFile))
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// write index index
	err = indexTemplate.Execute(indexFile, collection)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	// create book indexes and cover image
	for _, book := range collection.Books {
		if strings.TrimSpace(book.Title) == "" {
			continue
		}

		g.Go(func() error {
			bookInputDir := filepath.Join(workingDir, book.BookDir)
			bookOutputDir := filepath.Join(outputDir, book.ID)
			err = os.MkdirAll(bookOutputDir, os.ModePerm)
			if err != nil {
				return err
			}

			// create book index file
			indexFile, err := os.Create(filepath.Join(bookOutputDir, BookOutputIndexFile))
			if err != nil {
				return err
			}
			defer indexFile.Close()

			err = bookTemplate.Execute(indexFile, book)
			if err != nil {
				return err
			}

			// cover image
			coverName := book.CoverImageName
			if strings.TrimSpace(coverName) != "" {
				oldCoverPath := filepath.Join(bookInputDir, book.CoverImageName)
				newCoverPath := filepath.Join(bookOutputDir, book.CoverImageName)
				err = os.Link(oldCoverPath, newCoverPath)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return err
	}

	// create chapters
	g = new(errgroup.Group)
	for _, book := range collection.Books {
		if strings.TrimSpace(book.Title) == "" {
			continue
		}

		bookDir := filepath.Join(outputDir, book.ID)
		for _, chapter := range book.Chapters {
			g.Go(func() error {
				// create chapter index file
				chapterFile, err := os.Create(filepath.Join(bookDir, chapter.ID+".html"))
				if err != nil {
					return err
				}
				defer chapterFile.Close()

				err = chapterTemplate.Execute(chapterFile, chapter)
				if err != nil {
					return err
				}

				return nil
			})
		}
	}
	err = g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func copyDirectoryToOutput(inputDir, outputDir string, excludes []string, copySubdirs bool) error {
	dirItems, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, item := range dirItems {
		// check if item is excluded
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

		_ = excludes

		fullPath := filepath.Join(inputDir, item.Name())
		outputPath := filepath.Join(outputDir, item.Name())

		// copy subdirectories if needed
		if item.IsDir() && copySubdirs {
			err = copyDirectoryToOutput(fullPath, outputPath, excludes, copySubdirs)
			if err != nil {
				return err
			}
			continue
		} else if item.IsDir() && !copySubdirs {
			continue
		}

		// override outputPath contents
		err = os.RemoveAll(outputPath)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		if err != nil {
			return err
		}

		// copy to outputPath
		err = os.Link(fullPath, outputPath)
		if err != nil {
			return err
		}
	}

	return nil
}
