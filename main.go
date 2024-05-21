package main

import (
	"log"

)

func main() {
	book := Book{}

	err := UnmarshalBookConfigFile("mkbk.yml", &book)
	if err != nil {
		log.Fatal(err)
	}

	err = RenderBookToHTMLSite("./", book.OutputDir, &book)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Book Params:\n%#v\n",
	// 	book.Params["authors"].([]any)[0].(map[string]any)["bio"])
}
