package main

import (
	"github.com/nktauserum/anonymous-messages/internal/app"
	"log"
)

func main() {
	application := app.NewApplication()

	if err := application.Run(); err != nil {
		log.Fatalln(err)
	}
}
