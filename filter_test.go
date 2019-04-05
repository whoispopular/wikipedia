package wikipedia

import (
	"testing"
)

func pagesToChan(pages []Page) chan Page {
	c := make(chan Page, len(pages))
	for _, page := range pages {
		c <- page
	}
	close(c)
	return c
}

func chanToPages(c <-chan Page) []Page {
	pages := []Page{}
	for page := range c {
		pages = append(pages, page)
	}
	return pages
}

func isPagesEqual(c1 []Page, c2 []Page, t *testing.T) {
	if len(c1) != len(c2) {
		t.Fatalf(
			"Length for output is incorrect got %d expected %d \n\tTest case: `%#v`\n\tOutput: `%#v`",
			len(c1),
			len(c2),
			c1,
			c2,
		)
		return
	}

	for k := range c1 {
		if c1[k] != c2[k] {
			t.Fatalf("Output doesn't match input: \nExpected %#v\nGot: %#v", c1[k], c2[k])
		}
	}
}

func TestFilterByCategory(t *testing.T) {
	tests := []struct {
		filter CategoryFilter
		input  []Page
		output []Page
	}{
		{
			filter: CategoryFilter{
				Pattern: "\\[\\[Category:(.+)]\\]",
				Allowed: []string{},
			},
			input:  []Page{Page{Text: "[[Category:hello]]"}},
			output: []Page{},
		},
		{
			filter: CategoryFilter{
				Pattern: "\\[\\[Category:(.+)]\\]",
				Allowed: []string{"hello"},
			},
			input:  []Page{Page{Text: "[[Category:hello]]"}},
			output: []Page{Page{Text: "[[Category:hello]]"}},
		},
		{
			filter: CategoryFilter{
				Pattern: "\\[\\[Category:(.+)]\\]",
				Allowed: []string{"hello"},
				Denied:  []string{"hello world"},
			},
			input: []Page{
				Page{Text: "[[Category:hello world]]"},
				Page{Text: "[[Category:hello there]]"},
			},
			output: []Page{Page{Text: "[[Category:hello there]]"}},
		},
	}

	for _, test := range tests {
		inputChan := pagesToChan(test.input)
		outputChan := make(chan Page, len(test.input))
		FilterByCategory(inputChan, outputChan, test.filter)

		output := chanToPages(outputChan)
		isPagesEqual(test.output, output, t)
	}
}

func TestFilterByTitleLength(t *testing.T) {
	tests := []struct {
		filter TitleLengthFilter
		input  []Page
		output []Page
	}{
		{
			filter: TitleLengthFilter{Min: 0, Max: 5},
			input:  []Page{Page{Title: "1234567890"}},
			output: []Page{},
		},
		{
			filter: TitleLengthFilter{Min: 5, Max: 10},
			input:  []Page{Page{Title: "1234567890"}},
			output: []Page{Page{Title: "1234567890"}},
		},
		{
			filter: TitleLengthFilter{Min: 5, Max: 10},
			input:  []Page{Page{Title: "1234"}},
			output: []Page{},
		},
	}

	for _, test := range tests {
		inputChan := pagesToChan(test.input)
		outputChan := make(chan Page, len(test.input))
		FilterByTitleLength(inputChan, outputChan, test.filter)

		output := chanToPages(outputChan)
		isPagesEqual(test.output, output, t)
	}
}

func TestFilterByRedirect(t *testing.T) {
	tests := []struct {
		input  []Page
		output []Page
	}{
		{
			input:  []Page{Page{Redir: Redirect{Title: "another page"}}},
			output: []Page{},
		},
		{
			input:  []Page{Page{Redir: Redirect{Title: ""}}},
			output: []Page{Page{Redir: Redirect{Title: ""}}},
		},
		{
			input:  []Page{Page{}},
			output: []Page{Page{}},
		},
	}

	for _, test := range tests {
		inputChan := pagesToChan(test.input)
		outputChan := make(chan Page, len(test.input))
		FilterByRedirect(inputChan, outputChan)

		output := chanToPages(outputChan)
		isPagesEqual(test.output, output, t)
	}
}
