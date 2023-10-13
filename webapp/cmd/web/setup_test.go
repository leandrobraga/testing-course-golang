package main

import (
	"os"
	"testing"

	"github.com/leandrobraga/testing-course-golang/webapp/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	// Use the lines below in case use real database for test. It's not common for unit test

	// app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	// conn, err := app.connectToDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()

	// app.Session = getSession()
	// app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
