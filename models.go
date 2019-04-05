package wikipedia

import (
	"encoding/xml"
	"io"
	"regexp"
	"strings"
)

type Parser struct {
	x *xml.Decoder
}

func NewParser(r io.Reader) (*Parser, error) {
	d := xml.NewDecoder(r)
	_, err := d.Token()
	if err != nil {
		return nil, err
	}

	return &Parser{
		x: d,
	}, nil
}

func (p *Parser) Next() (*Page, error) {
	rv := &Page{}
	return rv, p.x.Decode(rv)
}

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	ID    uint64   `xml:"id"`
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

func (p *Page) Categories(categoryRegexp *regexp.Regexp) string {
	matches := categoryRegexp.FindAllStringSubmatch(p.Text, -1)
	categories := []string{}
	for _, c := range matches {
		categories = append(categories, strings.ToLower(c[1]))
	}
	return strings.Join(categories, ",")
}

type MinifiedPage struct {
	ID    uint64 `xml:"id"`
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

func NewMinifiedPage(p Page) MinifiedPage {
	return MinifiedPage{
		ID:    p.ID,
		Title: p.Title,
		Text:  FirstParagraph(p.Text),
	}
}
