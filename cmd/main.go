package main

import (
	"log"

	"github.com/caseapia/goproject-flush/internal/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Listen(":8080"))
}
