package main

import (
	"WikiIndex/app"
	"fmt"
	"log"
	"os"
)


const filename = "testdata/simplewiki-20170820-pages-meta-current.xml"
//const filename = "testdata/sample.xml"


func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := app.New()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	go func() {
		if err = app.AddWikiIndex(f); err != nil {
			fmt.Println(err)
		}
	}()

	return app.Serve()
}