package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JessebotX/mkbk"
)

func main() {
	collection := mkbk.Collection{}

	data, err := os.ReadFile(filepath.Join("testdata", "1", "mkbk.yml"))
	if err != nil {
		log.Fatal(err)
	}

	err = mkbk.Unmarshal(data, &collection)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", collection)
}
