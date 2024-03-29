package wiki

import (
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	linkExpression = `\[\[([^\]\[:]+)\]\]`
)

type Page struct {
	XMLName  xml.Name   `xml:"page"`
	Title    string     `xml:"title"`
	Revision []Revision `xml:"revision"`
}

type Revision struct {
	XMLName xml.Name `xml:"revision"`
	Comment string   `xml:"comment"`
	Text    string   `xml:"text"`
}

var blacklist = []string{
	"User",
	"Talk",
	"Category",
	"Wikipedia",
	"MediaWiki",
	"Template",
	"User talk",
	"Wikipedia talk",
	"File talk",
	"MediaWiki talk",
	"Template talk",
	"Help talk",
	"Category talk",
	"Portal talk",
	"Book talk",
	"Draft talk",
	"TimedText talk",
	"Module talk",
	"Special",
	"Help",
	"Module",
}

type Result map[string][]string

func Process(r io.Reader, count *int) (Result, error) {
	result := Result{}

	decoder := xml.NewDecoder(r)

PageIteration:
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to parse token")
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var page Page
				if err := decoder.DecodeElement(&page, &se); err != nil {
					return nil, errors.Wrap(err, "failed to decode element")
				}

				if len(page.Revision) == 0 {
					continue PageIteration
				}

				title := page.Title
				content := page.Revision[0].Text

				for _, prefix := range blacklist {
					if strings.HasPrefix(title, prefix+":") {
						continue PageIteration
					}
				}

				result[title] = parseLinks(content)
				*count++
			}
		}
	}

	return result, nil
}

func parseLinks(content string) []string {
	urlMatcher := regexp.MustCompile(linkExpression)

	matches := urlMatcher.FindAllString(content, -1)

	references := make([]string, len(matches))

	for i, match := range matches {
		match = strings.TrimLeft(match, "[[")
		match = strings.TrimRight(match, "]]")

		link := match

		if strings.Contains(match, "|") {
			link = strings.Split(match, "|")[0]
		}

		references[i] = link
	}

	return references
}
