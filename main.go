package main

import (
	"WikiIndex/app"
	"compress/bzip2"
	"fmt"
	"log"
	"os"
)


const filename = "testdata/enwiki-20190101-pages-articles-multistream.xml.bz2"
//const filename = "testdata/simplewiki-20170820-pages-meta-current.xml"
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

	r := bzip2.NewReader(f)

	go func() {
		if err = app.AddWikiIndex(r); err != nil {
			fmt.Println(err)
		}
		f.Close()
	}()

	return app.Serve()
}