package database

import (
	"errors"
	"fmt"
	_ "github.com/pkg/errors"
	"strings"
)

type Index interface {
	Get(title string) (*Page, error)
	BatchProcess(map[string][]string)
	Size() int
}

type MapIndex struct {
	index map[string]*Page
}

func New() *MapIndex {
	return &MapIndex{
		map[string]*Page{},
	}
}

func (i *MapIndex) Get(title string) (*Page, error) {
	title = i.normalize(title)

	page, ok := i.index[title]
	if !ok {
		return nil, errors.New("page not found")
	}

	return page, nil
}

func (i *MapIndex) add(title string) {
	title = i.normalize(title)

	page := &Page{
		title,
		make(map[*Page]bool),
	}

	i.index[title] = page
}

//func (i *MapIndex) Reindex() {
//	for _, page := range i.index {
//		i.process(page)
//	}
//}
//
//func (i *MapIndex) process(p *Page) {
//	for _, title := range p.parseLinks() {
//		if ref, ok := i.index[title]; ok {
//			p.references = append(p.references, ref)
//		}
//	}
//}

func (i *MapIndex) Size() int {
	return len(i.index)
}

func (i *MapIndex) BatchProcess(data map[string][]string) {

	// Add all entries to database
	fmt.Println("Adding entries")
	for title := range data {
		i.add(title)
	}

	// Add references to entries
	fmt.Println("Adding References")
	for title, references := range data {
		p, err := i.Get(title)
		if err != nil {
			continue
		}

		for _, reference := range references {
			r, err := i.Get(reference)
			if err != nil {
				continue
			}

			p.References[r] = true
		}
	}
}

func (i *MapIndex) normalize(title string) string {
	title = strings.TrimSpace(title)
	return strings.ToLower(title)
}