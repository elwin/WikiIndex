package database

type Page struct {
	Title      string
	References map[*Page]bool
}

func (p *Page) String() string {
	return p.Title
}