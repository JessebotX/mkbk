package main

import (
	"log"

	"github.com/JessebotX/mkbk"
)

func main() {
	book := mkbk.Book{}

	err := mkbk.UnmarshalBookConfigFile("mkbk-book.yml", &book)
	if err != nil {
		log.Fatal(err)
	}

	err = mkbk.RenderBookToHTMLSite("./", book.OutputDir, &book)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Book Params:\n%#v\n",
	// 	book.Params["authors"].([]any)[0].(map[string]any)["bio"])
}
