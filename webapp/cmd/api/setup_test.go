package main

import (
	"os"
	"testing"

	"github.com/leandrobraga/testing-course-golang/webapp/pkg/repository/dbrepo"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2OTczNzg5ODYsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.7xsYYKybGqOIyVU3yRqEnd6fAWP5hLkx_Jaw6rsrK8A"

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "aksdhadhaksdhakdhakhdasdhkj"

	os.Exit(m.Run())
}
