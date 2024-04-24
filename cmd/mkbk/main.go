package main

import (
	"fmt"

	"github.com/JessebotX/mkbk"
)

func main() {
	fmt.Println("Hello, world")

	book := mkbk.Book{
		Title: "Hello",
	}

	fmt.Println(book)
}
