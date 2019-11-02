package app

import (
	"WikiIndex/database"
	"WikiIndex/database/wiki"
	"fmt"
	"io"
)

type App struct {
	index database.Index
}

func New() *App {
	return &App{
		database.New(),
	}
}

func (a *App) AddWikiIndex (r io.ReadCloser) error {
	fmt.Println("Processing file")
	x, err := wiki.Process(r)
	if err != nil {
		return err
	}

	fmt.Println("Processing entries")
	a.index.BatchProcess(x)

	return nil
}
