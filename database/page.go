package database

import (
	"fmt"
	"strings"
)

type Pageable interface {
	Id() int
	Title() string
	Slug() string
	ReferencesTo() []Pageable
	ReferencedBy() []Pageable
	AddReferenceTo(page Pageable)
	AddReferenceBy(page Pageable)
	WikipediaUrl() string
}

type Page struct {
	id           int
	title        string
	slug         string
	referencesTo map[int]bool
	referencedBy map[int]bool
	index        Index
}

func NewPage(title string, i Index) *Page {
	return &Page{
		i.NewIndex(),
		title,
		i.UniqueSlug(title),
		make(map[int]bool),
		make(map[int]bool),
		i,
	}
}

func (p *Page) Id() int {
	return p.id
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) Slug() string {
	return p.slug
}

func (p *Page) ReferencesTo() []Pageable {
	result := make([]Pageable, 0)

	for id := range p.referencesTo {
		page, ok := p.index.Get(id)
		if !ok {
			continue
		}

		result = append(result, page)
	}

	return result
}

func (p *Page) AddReferenceTo(page Pageable) {
	p.referencesTo[page.Id()] = true
}

func (p *Page) ReferencedBy() []Pageable {
	result := make([]Pageable, 0)

	for id := range p.referencedBy {
		page, ok := p.index.Get(id)
		if !ok {
			continue
		}

		result = append(result, page)
	}

	return result
}

func (p *Page) AddReferenceBy(page Pageable) {
	p.referencedBy[page.Id()] = true
}

func (p *Page) WikipediaUrl() string {
	title := strings.Replace(p.title, " ", "_", -1)
	return fmt.Sprintf("https://simple.wikipedia.org/wiki/%s", title)
}
