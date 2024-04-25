package main

import (
	"fmt"
	"log"

	"github.com/JessebotX/mkbk"
)

func main() {
	fmt.Println("Hello, world")

	book := mkbk.Book{
		Title: "Hello",
	}

	err := mkbk.UnmarshalBookConfigFile("mkbk-book.yml", &book)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("----\nBook\n----\n%#v\n", book)
	// fmt.Printf("Book Params:\n%#v\n",
	// 	book.Params["authors"].([]any)[0].(map[string]any)["bio"])
}
