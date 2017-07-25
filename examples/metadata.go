package main

import (
	"io"
	"log"
	"os"

	"github.com/gohugoio/hugo/parser"
)

func main() {
	page, err := parser.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	metadata := page.FrontMatter()
	if len(metadata) < 1 {
		log.Fatalln("No FrontMatter")
	}
	data, err := parser.DetectFrontMatter(rune(metadata[0])).Parse(metadata)
	if err != nil {
		log.Fatalln(err)
	}
	r, err := MakeTable(data)
	if err != nil {
		log.Fatalln(err)
	}
	io.Copy(os.Stdout, r)
}
