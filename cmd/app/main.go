package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typimain"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	app := typimain.NewTypicalApplication(typical.Context)
	err := app.Cli().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
