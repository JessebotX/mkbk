package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"github.com/JessebotX/mkbk"
)

func main() {
	yamlBody, err := os.ReadFile(filepath.Join("testdata", "1", "mkbk.yml"))
	if err != nil {
		log.Fatal(err)
	}

	var collection mkbk.Collection

	err = yaml.Unmarshal(yamlBody, &collection)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", collection)
}
