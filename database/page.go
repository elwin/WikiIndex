package database

type Pageable interface {
	Title() string
	ReferencesTo() []Pageable
	ReferencedBy() []Pageable
	AddReferenceTo(page Pageable)
	AddReferenceBy(page Pageable)
	String() string
}

type Page struct {
	title        string
	referencesTo map[string]bool
	referencedBy map[string]bool
	index        Index
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) String() string {
	return p.title
}

func (p *Page) ReferencesTo() []Pageable {
	result := make([]Pageable, 0)

	for title := range p.referencesTo {
		page, ok := p.index.Get(title)
		if !ok {
			continue
		}

		result = append(result, page)
	}

	return result
}

func (p *Page) AddReferenceTo(page Pageable) {
	p.referencesTo[page.Title()] = true
}

func (p *Page) ReferencedBy() []Pageable {
	result := make([]Pageable, 0)

	for title := range p.referencedBy {
		page, ok := p.index.Get(title)
		if !ok {
			continue
		}

		result = append(result, page)
	}

	return result
}

func (p *Page) AddReferenceBy(page Pageable) {
	p.referencedBy[page.Title()] = true
}