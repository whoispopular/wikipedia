package wikipedia

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

func PrintTitlesAndCategories(pages <-chan Page, categoryFilter CategoryFilter) {
	for p := range pages {
		fmt.Printf("%s \n\tCategories: %s\n", p.Title, p.Categories(categoryFilter.CategoryRegexp()))
	}
}

func PrintTitles(pages <-chan Page) {
	for p := range pages {
		fmt.Println(p.Title)
	}
}

func WritePagesToJSON(pages <-chan Page) {
	for page := range pages {
		output, err := json.Marshal(page)
		if err == nil {
			fmt.Println(string(output))
		}
	}
}

func WritePagesToXML(pages <-chan Page) {
	for page := range pages {
		output, err := xml.Marshal(page)
		if err == nil {
			fmt.Println(string(output))
		} else {
			log.Printf("Error while marshaling page `%s`", err)
		}
	}
}

func PrintAsText(pages <-chan Page, categoryFilter CategoryFilter) {
	for p := range pages {
		fmt.Printf(
			"%s \n\tCategories: %s\n\tExcerpt: %s\n",
			p.Title,
			p.Categories(categoryFilter.CategoryRegexp()),
			FirstParagraph(p.Text),
		)
	}
}
