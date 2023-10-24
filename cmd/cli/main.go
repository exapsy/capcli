package main

import "github.com/exapsy/capcli/internal/cli"

func main() {
	var err error
	app := cli.NewApp()
	err = app.Run()
	if err != nil {
		panic(err)
	}
}
