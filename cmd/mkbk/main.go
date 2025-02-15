package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JessebotX/mkbk"
)

const (
	CollectionConfigFileName = "mkbk.yml"
	BookConfigFileName = "mkbk-book.yml"
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
		var book mkbk.Book

		fullBookDirPath := filepath.Join(workingDir, bookDir)
		bookData, err := os.ReadFile(filepath.Join(fullBookDirPath, BookConfigFileName))
		if err != nil {
			log.Fatal(err)
		}

		err = mkbk.UnmarshalBook(bookData, &book)
		if err != nil {
			log.Fatal(err)
		}

		books = append(books, book)
	}
	collection.Books = books

	fmt.Printf("%#v\n", collection)
}
