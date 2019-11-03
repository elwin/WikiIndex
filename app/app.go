package app

import (
	"WikiIndex/database"
	"WikiIndex/database/wiki"
	"fmt"
	"io"
)

type App struct {
	Index database.Index
	Count *int
}

func New() *App {
	var i int
	return &App{
		database.New(),
		&i,
	}
}

func (a *App) AddWikiIndex (r io.Reader) error {
	fmt.Println("Processing file")
	x, err := wiki.Process(r, a.Count)
	if err != nil {
		return err
	}

	fmt.Println("Processing entries")
	a.Index.BatchProcess(x)

	return nil
}
