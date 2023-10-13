package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/leandrobraga/testing-course-golang/webapp/pkg/data"
	"github.com/leandrobraga/testing-course-golang/webapp/pkg/repository/dbrepo"

	"github.com/leandrobraga/testing-course-golang/webapp/pkg/repository"
)

type application struct {
	Session *scs.SessionManager
	DSN     string
	DB      repository.DatabaseRepo
}

func main() {
	gob.Register(data.User{})
	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.Session = getSession()
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	mux := app.routes()

	log.Println("Starting server on port 8080...")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
