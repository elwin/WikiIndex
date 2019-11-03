package database

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type Index interface {
	Get(title string) (Pageable, bool)
	Path(from, to Pageable) ([]Pageable, error)
	LongestPath(from Pageable) (Pageable, int)
	LongestTotalPath() (from, to Pageable, cost int)
	BatchProcess(map[string][]string)
	Size() int
	Slugify(string) string
}

type MapIndex struct {
	index map[string]*Page
}

func New() *MapIndex {
	return &MapIndex{
		map[string]*Page{},
	}
}

func (i *MapIndex) Get(title string) (Pageable, bool) {
	title = i.Slugify(title)

	page, ok := i.index[title]
	if !ok {
		return nil, false
	}

	return page, true
}

func (i *MapIndex) add(title string) {
	page := NewPage(title, i)

	i.index[page.Slug()] = page
}

func (i *MapIndex) Size() int {
	return len(i.index)
}

func (i *MapIndex) Path(from, to Pageable) ([]Pageable, error) {
	cost := map[Pageable]int{}

	queue := make([]Pageable, 0)
	queue = append(queue, from)

	for {
		if len(queue) == 0 {
			return nil, errors.New("no path found")
		}

		current := queue[0]
		queue = queue[1:]

		// Found neighbour
		if current == to {
			break
		}

		currentCost := cost[current]

		for _, neighbour := range current.ReferencesTo() {

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

	path := make([]Pageable, 0)
	path = append(path, to)

	for {

		current := path[len(path) - 1]
		nextCost := cost[current] - 1

		if nextCost == 0 {
			path = append(path, from)
			break
		}

		for _, backref := range current.ReferencedBy() {
			if cost[backref] == nextCost {
				path = append(path, backref)
				break
			}
		}
	}

	// Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

func (i *MapIndex) LongestPath(from Pageable) (Pageable, int) {
	cost := map[Pageable]int{}
	cost[from] = 0

	queue := make([]Pageable, 0)
	queue = append(queue, from)

	for {
		if len(queue) == 0 {
			break
		}

		current := queue[0]
		queue = queue[1:]

		currentCost := cost[current]

		for _, neighbour := range current.ReferencesTo() {

			// Already visited
			if _, ok := cost[neighbour]; ok {
				continue
			}

			// Assign cost
			cost[neighbour] = currentCost + 1

			// Add the queue
			queue = append(queue, neighbour)
		}
	}

	maxCost := 0
	var node Pageable

	for n, c := range cost {
		if c > maxCost {
			maxCost = c
			node = n
		}
	}

	return node, maxCost
}

func (i *MapIndex) LongestTotalPath() (from, to Pageable, cost int) {
	maxCost := 0
	var maxFrom Pageable
	var maxTo Pageable

	for _, from := range i.index {
		to, cost := i.LongestPath(from)
		if cost > maxCost {
			maxCost = cost
			maxFrom = from
			maxTo = to
		}
	}

	return maxFrom, maxTo, maxCost
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
		p, ok := i.Get(title)
		if !ok {
			continue
		}

		for _, referenceTitle := range references {
			reference, ok := i.Get(referenceTitle)
			if !ok {
				continue
			}

			// Add reference
			p.AddReferenceTo(reference)

			// Add back reference
			reference.AddReferenceBy(p)
		}

	}
}

func (i *MapIndex) Slugify(title string) string {
	//title = strings.TrimSpace(title)
	//title = strings.ToLower(title)
	return slug.Make(title)
}