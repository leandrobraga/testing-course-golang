package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestAppHome(t *testing.T) {
	var tests = []struct {
		name            string
		putInSession    string
		expectedSession string
	}{
		{"first visit", "", "<small>From Session:"},
		{"sec visit", "hello, world!", "<small>From Session: hello, world!"},
	}

	for _, e := range tests {
		// create a request
		req, _ := http.NewRequest("GET", "/", nil)

		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)

		handler.ServeHTTP(rr, req)

		// check status code
		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedSession) {
			t.Errorf("%s: did not find %s in response body", e.name, e.expectedSession)
		}
	}
}

func TestApp_renderWithbadTemplate(t *testing.T) {
	// set templatepath to a location with bad template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Errorf("expected error from bad tempalte, but did not get one")
	}

	// Reset to original value for the other test dont have problem
	pathToTemplates = "./../../templates/"

}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "uknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}
