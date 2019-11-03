package main

import (
	"WikiIndex/app"
	"fmt"
)

func main() {
	app := app.New()
	app.Index.BatchProcess(map[string][]string{
		"a": {"b yo", "c"},
		"b yo": {"c", "d"},
		"c": {"d", "a"},
		"d": {},
	})

	if err := app.Serve(); err != nil {
		fmt.Println(err)
	}
}
