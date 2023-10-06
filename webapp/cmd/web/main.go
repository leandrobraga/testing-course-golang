package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	app := application{}

	mux := app.routes()

	log.Println("Starting server on port 8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
