package wikipedia

import (
	"regexp"
	"strings"
)

type RegexReplacer struct {
	r *regexp.Regexp
	s string
}

func (r *RegexReplacer) Replace(str string) string {
	return r.r.ReplaceAllString(str, r.s)
}

func NewRegexReplacer(reg string, substitute string) *RegexReplacer {
	return &RegexReplacer{
		r: regexp.MustCompile(reg),
		s: substitute,
	}
}

func RenderMarkupToHTML(str string) string {
	replacers := []*RegexReplacer{
		NewRegexReplacer(`(?s){{[^{}]*}}`, ""),
		NewRegexReplacer(`<ref[^>]*>([^<]*)</ref>`, ""),
		NewRegexReplacer(`<ref[^>]*/>`, ""),
		NewRegexReplacer(`(?s)\{\|.+\|\}`, ""),
		NewRegexReplacer(`(?m)^======(.*)======`, "<br /><h6>$1</h6>"),
		NewRegexReplacer(`(?m)^=====(.*)=====`, "<br /><h5>$1</h5>"),
		NewRegexReplacer(`(?m)^====(.*)====`, "<br /><h4>$1</h4>"),
		NewRegexReplacer(`(?m)^===(.*)===`, "<br /><h3>$1</h3>"),
		NewRegexReplacer(`(?m)^==(.*)==`, "<br /><h2>$1</h2>"),
		NewRegexReplacer(`(?m)^=(.*)=`, "<br /><h1>$1</h1>"),
		NewRegexReplacer(`----`, "<hr />"),
		NewRegexReplacer(`__FORCETOC__`, ""),
		NewRegexReplacer(`__TOC__`, ""),
		NewRegexReplacer(`__NOTOC__`, ""),
		NewRegexReplacer(`(?m)^\*+(.*)$`, "<li>$1</li>"),
		NewRegexReplacer(`(?m)^#+(.*)$`, "<li>$1</li>"),
		NewRegexReplacer(`'''''(.*)'''''`, "<strong><i>$1</i></strong>"),
		NewRegexReplacer(`'''(.*)'''`, "<strong>$1</strong>"),
		NewRegexReplacer(`''(.*)''`, "<i>$1</i>"),
		NewRegexReplacer(`\[(https?://[^\[\]]+)\s([^\[\]]+)\]`, "$2"),
		NewRegexReplacer(`\[\[([^\[\]]+)\|([^\[\]]+)\]\]`, "$2"),
		NewRegexReplacer(`\[\[([^\[\]]+)\|\]\]`, "$1"),
		NewRegexReplacer(`\[\[([^\[\]]*)\]\]`, "$1"),
	}
	for _, replacer := range replacers {
		for {
			newstr := replacer.Replace(str)
			if newstr == str {
				break
			}
			str = newstr
		}
	}

	str = strings.TrimSpace(str)
	str = NewRegexReplacer(`\n\n`, "<br />").Replace(str)

	return str
}

func RenderMarkupToText(str string) string {
	replacers := []*RegexReplacer{
		NewRegexReplacer(`(?s){{[^{}]*}}`, ""),
		NewRegexReplacer(`<ref[^>]*>([^<]*)</ref>`, ""),
		NewRegexReplacer(`<ref[^>]*/>`, ""),
		NewRegexReplacer(`(?s)\{\|.+\|\}`, ""),
		NewRegexReplacer(`(?m)^======(.*)======`, "\n$1"),
		NewRegexReplacer(`(?m)^=====(.*)=====`, "\n$1"),
		NewRegexReplacer(`(?m)^====(.*)====`, "\n$1"),
		NewRegexReplacer(`(?m)^===(.*)===`, "\n$1"),
		NewRegexReplacer(`(?m)^==(.*)==`, "\n$1"),
		NewRegexReplacer(`(?m)^=(.*)=`, "\n$1"),
		NewRegexReplacer(`----`, ""),
		NewRegexReplacer(`__FORCETOC__`, ""),
		NewRegexReplacer(`__TOC__`, ""),
		NewRegexReplacer(`__NOTOC__`, ""),
		NewRegexReplacer(`(?m)^\*+(.*)$`, "$1"),
		NewRegexReplacer(`(?m)^#+(.*)$`, "$1"),
		NewRegexReplacer(`'''''(.*)'''''`, "$1"),
		NewRegexReplacer(`'''(.*)'''`, "$1"),
		NewRegexReplacer(`''(.*)''`, "$1"),
		NewRegexReplacer(`\[(https?://[^\[\]]+)\s([^\[\]]+)\]`, "$2"),
		NewRegexReplacer(`\[\[([^\[\]]+)\|([^\[\]]+)\]\]`, "$2"),
		NewRegexReplacer(`\[\[([^\[\]]+)\|\]\]`, "$1"),
		NewRegexReplacer(`\[\[([^\[\]]*)\]\]`, "$1"),
	}
	for _, replacer := range replacers {
		for {
			newstr := replacer.Replace(str)
			if newstr == str {
				break
			}
			str = newstr
		}
	}

	str = strings.TrimSpace(str)

	return str
}

func FirstParagraph(str string) string {
	str = RenderMarkupToText(str)
	return strings.Split(str, "\n\n")[0]
}
