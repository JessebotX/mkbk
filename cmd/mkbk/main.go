package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/JessebotX/mkbk"
)

const (
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
	collection := mkbk.Collection{}

	workingDir := filepath.Join("testdata", "1")
	data, err := os.ReadFile(filepath.Join(workingDir, CollectionConfigFileName))
	if err != nil {
		log.Fatal(err)
	}

	err = mkbk.Unmarshal(data, &collection)
	if err != nil {
		log.Fatal(err)
	}

	books := make([]mkbk.Book, len(collection.BookDirs))
	for _, bookDir := range collection.BookDirs {
		book := mkbk.Book{
			ID: filepath.Base(bookDir),
			BookDir: filepath.Join(workingDir, bookDir),
		}

		fullBookDirPath := filepath.Join(workingDir, bookDir)
		bookData, err := os.ReadFile(filepath.Join(fullBookDirPath, BookConfigFileName))
		if err != nil {
			log.Fatal(err)
		}

		err = mkbk.UnmarshalBook(bookData, &book, &collection)
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
