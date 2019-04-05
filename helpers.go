package wikipedia

import (
	"regexp"
	"strings"
)

func ReadAllPages(parser Parser, pages chan<- Page) {
	for {
		page, err := parser.Next()

		if err != nil {
			close(pages)
			return
		}

		page.Title = CleanTitle(page.Title)
		pages <- *page
	}
}

var ignoreTitlePattern = regexp.MustCompile("\\(.+\\)|-|_|\n")
var multipleSpacesPattern = regexp.MustCompile("\\s+")

func CleanTitle(str string) string {
	cleantitle := ignoreTitlePattern.ReplaceAllLiteralString(str, " ")
	compactSpaces := multipleSpacesPattern.ReplaceAllLiteralString(cleantitle, " ")
	return strings.TrimSpace(compactSpaces)
}
