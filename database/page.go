package database

import (
	"fmt"
	"strings"
)

type Pageable interface {
	Title() string
	Slug() string
	ReferencesTo() []Pageable
	ReferencedBy() []Pageable
	AddReferenceTo(page Pageable)
	AddReferenceBy(page Pageable)
	WikipediaUrl() string
}

type Page struct {
	title        string
	slug         string
	referencesTo map[string]bool
	referencedBy map[string]bool
	index        Index
}

func NewPage(title string, i Index) *Page {
	return &Page{
		title,
		i.UniqueSlug(title),
		make(map[string]bool),
		make(map[string]bool),
		i,
	}
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) Slug() string {
	return p.slug
}

func (p *Page) ReferencesTo() []Pageable {
	result := make([]Pageable, 0)

	for slug := range p.referencesTo {
		page, ok := p.index.Get(slug)
		if !ok {
			continue
		}

		result = append(result, page)
	}

	return result
}

func (p *Page) AddReferenceTo(page Pageable) {
	p.referencesTo[page.Slug()] = true
}

func (p *Page) ReferencedBy() []Pageable {
	result := make([]Pageable, 0)

	for slug := range p.referencedBy {
		page, ok := p.index.Get(slug)
		if !ok {
			continue
		}

		result = append(result, page)
	}

	return result
}

func (p *Page) AddReferenceBy(page Pageable) {
	p.referencedBy[page.Slug()] = true
}

func (p *Page) WikipediaUrl() string {
	title := strings.Replace(p.title, " ", "_", -1)
	return fmt.Sprintf("https://simple.wikipedia.org/wiki/%s", title)
}
