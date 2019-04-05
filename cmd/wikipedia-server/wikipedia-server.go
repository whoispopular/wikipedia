package main

import (
	"log"
	"os"
	"wikipedia"

	"github.com/gin-gonic/gin"
)

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	parser, err := wikipedia.NewParser(f)
	if err != nil {
		log.Fatal(err)
	}

	pipe1 := make(chan wikipedia.Page, 1000)
	go wikipedia.ReadAllPages(*parser, pipe1)
	indexPages(pipe1)
	f.Close()

	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.GET("/people/:name", personHandler)

	log.Println("Starting server...")
	r.Run()
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

var pagesCache = map[string][]wikipedia.MinifiedPage{}

func indexPages(pages <-chan wikipedia.Page) {
	for page := range pages {
		mpage := wikipedia.NewMinifiedPage(page)
		pageArr, ok := pagesCache[mpage.Title]

		if !ok {
			pagesCache[mpage.Title] = []wikipedia.MinifiedPage{mpage}
		} else {
			pagesCache[mpage.Title] = append(pageArr, mpage)
		}
	}
}

func personHandler(c *gin.Context) {
	page, ok := pagesCache[c.Param("name")]
	if !ok {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, page)
}
