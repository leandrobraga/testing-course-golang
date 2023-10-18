package main

import (
	"os"
	"testing"

	"github.com/leandrobraga/testing-course-golang/webapp/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "aksdhadhaksdhakdhakhdasdhkj"

	os.Exit(m.Run())
}
