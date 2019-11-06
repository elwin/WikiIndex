package app

import (
	"WikiIndex/database"
	"WikiIndex/database/wiki"
	"fmt"
	"io"
)

type App struct {
	Index           database.Index
	Count           *int
	IndexInProgress bool
}

func New() *App {
	var i int
	return &App{
		database.New(),
		&i,
		true,
	}
}

func (a *App) AddWikiIndex(r io.Reader) error {
	fmt.Println("Processing file")
	x, err := wiki.Process(r, a.Count)
	if err != nil {
		return err
	}

	fmt.Println("Processing entries")
	a.Index.BatchProcess(x)

	fmt.Println("Finished processing")
	a.IndexInProgress = false

	return nil
}
