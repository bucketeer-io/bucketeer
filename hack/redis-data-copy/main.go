package main

import (
	"log"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
)

var (
	name    = "redis-data-copy"
	version = ""
	build   = ""
)

func main() {
	app := cli.NewApp(name, "Bucketeer Redis data copy tool", version, build)
	registerCommand(app, app)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
