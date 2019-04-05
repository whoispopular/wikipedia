package main

import (
	"compress/bzip2"
	"flag"
	"log"
	"os"
	"wikipedia"
)

func main() {

	config_path := flag.String("config", "", "pages filters configuration json file")
	format := flag.String("format", "xml", "output format, available formats: xml, json, title, title-cat, text")
	channelCap := flag.Int("channel-cap", 1000, "processing channels capacity between filters")
	flag.Parse()

	config, err := wikipedia.ReadConfig(*config_path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config: %#v", config)

	decompressed := bzip2.NewReader(os.Stdin)
	parser, err := wikipedia.NewParser(decompressed)
	if err != nil {
		log.Fatal(err)
	}

	pipe1 := make(chan wikipedia.Page, *channelCap)
	go wikipedia.ReadAllPages(*parser, pipe1)

	pipe2 := make(chan wikipedia.Page, *channelCap)
	go wikipedia.FilterByRedirect(pipe1, pipe2)

	pipe3 := make(chan wikipedia.Page, *channelCap)
	go wikipedia.FilterByCategory(pipe2, pipe3, config.Category)

	pipe4 := make(chan wikipedia.Page, *channelCap)
	go wikipedia.FilterByTitleLength(pipe3, pipe4, config.TitleLength)

	pipe5 := make(chan wikipedia.Page, *channelCap)
	go wikipedia.FilterByTitle(pipe4, pipe5, config.Title)

	switch *format {
	case "json":
		wikipedia.WritePagesToJSON(pipe5)
	case "xml":
		wikipedia.WritePagesToXML(pipe5)
	case "title":
		wikipedia.PrintTitles(pipe5)
	case "title-cat":
		wikipedia.PrintTitlesAndCategories(pipe5, config.Category)
	case "text":
		wikipedia.PrintAsText(pipe5, config.Category)
	default:
		log.Fatalf("Can't use format `%s`", *format)
	}
}
