package wiki

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const expression = `\[\[([^\]\[:]+)\|([^\]\[:]+)\]\]`


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

type Result map[string][]string

func Process(r io.ReadCloser) (Result, error) {
	defer r.Close()

	result := Result{}

	decoder := xml.NewDecoder(r)
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

				if len(page.Revision) > 0 {
					title := page.Title
					content := page.Revision[0].Text
					if page.Title == "Architecture" {
						fmt.Println("yoooo")
					}
					result[title] = parseLinks(content)
				}
			}
		}
	}

	return result, nil
}


func parseLinks(content string) []string {
	urlMatcher := regexp.MustCompile(expression)

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


