package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTest = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/fish", http.StatusNotFound},
	}

	routes := app.routes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	pathToTemplates = "./../../templates/"

	for _, e := range theTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			log.Fatal(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}

}
