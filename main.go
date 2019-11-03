package main

import (
	"WikiIndex/app"
	"compress/bzip2"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var options struct {
	Filename string `short:"i" long:"input" description:"Wikipedia XML dump (something ending in .xml.bz2)" required:"true"`
	Address  string `short:"h" long:"http" description:"HTTP address" default:":8080"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	_, err := flags.Parse(&options)
	if err != nil {
		os.Exit(1)
	}

	app := app.New()

	f, err := os.Open(options.Filename)
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

	return app.Serve(options.Address)
}
