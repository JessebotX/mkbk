package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	Version = "1.0"
	CollectionConfigFileName = "mkbk.yml"
	BookConfigFileName = "mkbk-book.yml"

	IndexTemplateName   = "index.html"
	BookTemplateName = "index.html"
	BookFolderName = "_book"
	ChapterTemplateName = "_chapter.html"
	RSSFeedTemplateName = "_rss.xml" // TODO implement RSS feeds
	ImagesFolderName    = "images"
	CSSFolderName       = "css"
	BookOutputIndexFile = "index.html"
	CollectionOutputIndexFile = "index.html"
)

func main() {
	workingDir := "./"
	if len(os.Args) > 1 {
		switch arg := os.Args[1]; arg {
		case "-h":
			fmt.Printf("USAGE: %v <directory|-h|-v>\n", os.Args[0])
			os.Exit(0)
		case "-v":
			fmt.Printf("%v v%v\n", os.Args[0], Version)
			os.Exit(0)
		default:
			workingDir = filepath.Join(arg)
		}
	}

	collection := Collection{}

	data, err := os.ReadFile(filepath.Join(workingDir, CollectionConfigFileName))
	if err != nil {
		log.Fatal(err)
	}

	err = Unmarshal(data, &collection)
	if err != nil {
		log.Fatal(err)
	}

	books := make([]Book, 0)
	for _, bookDir := range collection.BookDirs {
		book := Book{
			ID: filepath.Base(bookDir),
			BookDir: filepath.Join(workingDir, bookDir),
		}

		fullBookDirPath := filepath.Join(workingDir, bookDir)
		bookData, err := os.ReadFile(filepath.Join(fullBookDirPath, BookConfigFileName))
		if err != nil {
			log.Fatal(err)
		}

		err = UnmarshalBook(bookData, &book, &collection)
		if err != nil {
			log.Fatal(err)
		}

		books = append(books, book)
	}
	collection.Books = books

	//fmt.Printf("%#v\n", collection)
	err = RenderCollectionToHTML(workingDir, collection)
	if err != nil {
		log.Fatal(err)
	}
}
