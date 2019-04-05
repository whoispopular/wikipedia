package wikipedia

import (
	"regexp"
	"strings"
)

type CategoryFilter struct {
	Pattern string
	Allowed []string
	Denied  []string
}

func (c *CategoryFilter) CategoryRegexp() *regexp.Regexp {
	return regexp.MustCompile(c.Pattern)
}

type TitleLengthFilter struct {
	Min int
	Max int
}

type TitleFilter struct {
	AllowContain  []string
	RemoveContain []string
	RemoveExact   []string
}

func FilterByTitleLength(pages <-chan Page, filtered chan<- Page, filter TitleLengthFilter) {
	for page := range pages {
		if len(page.Title) >= filter.Min && len(page.Title) <= filter.Max {
			filtered <- page
		}
	}

	close(filtered)
}

func FilterByCategory(pages <-chan Page, filtered chan<- Page, filter CategoryFilter) {

	inAnyCategory := func(page_categories string, categories []string) bool {
		for _, cat := range categories {
			if strings.Index(page_categories, cat) != -1 {
				return true
			}
		}
		return false
	}

	for page := range pages {
		categories := page.Categories(filter.CategoryRegexp())
		if inAnyCategory(categories, filter.Allowed) && !inAnyCategory(categories, filter.Denied) {
			filtered <- page
		}
	}

	close(filtered)
}

func FilterByRedirect(pages <-chan Page, filtered chan<- Page) {
	for page := range pages {
		if len(page.Redir.Title) == 0 {
			filtered <- page
		}
	}

	close(filtered)
}

func FilterByTitle(pages <-chan Page, filtered chan<- Page, filter TitleFilter) {
	inArray := func(s string, arr []string) bool {
		for _, match := range arr {
			if s == match {
				return true
			}
		}
		return false
	}

	containsAny := func(s string, arr []string) bool {
		for _, match := range arr {
			if strings.Index(s, match) != -1 {
				return true
			}
		}
		return false
	}

	for page := range pages {
		if containsAny(page.Title, filter.AllowContain) &&
			!containsAny(page.Title, filter.RemoveContain) &&
			!inArray(page.Title, filter.RemoveExact) {
			filtered <- page
		}
	}

	close(filtered)
}
