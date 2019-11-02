package database

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Index interface {
	Get(title string) (*Page, error)
	Path(from, to *Page) (int, error)
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
		return nil, errors.Errorf("page \"%s\" not found", title)
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

func (i *MapIndex) Size() int {
	return len(i.index)
}

func (i *MapIndex) Path(a, b *Page) (int, error) {
	cost := map[*Page]int{}

	queue := make([]*Page, 0)
	queue = append(queue, a)

	for {
		if len(queue) == 0 {
			return 0, errors.New("no path found")
		}

		current := queue[0]
		queue = queue[1:]

		// Found neighbour
		if current == b {
			break
		}

		currentCost := cost[current]

		for neighbour := range current.References {

			// Already visited
			if cost[neighbour] != 0 {
				continue
			}

			// Assign cost
			cost[neighbour] = currentCost + 1

			// Add the queue
			queue = append(queue, neighbour)
		}
	}

	path := make([]*Page, 0)
	path = append(path, a, b)

	return cost[b], nil
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