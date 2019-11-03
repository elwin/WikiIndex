package main

import (
	"WikiIndex/app"
	"fmt"
)

func main() {
	app := app.New()
	app.Index.BatchProcess(map[string][]string{
		"a": {"b", "c"},
		"b": {"c", "d"},
		"c": {"d", "a"},
		"d": {},
	})

	if err := app.Serve(); err != nil {
		fmt.Println(err)
	}
}
