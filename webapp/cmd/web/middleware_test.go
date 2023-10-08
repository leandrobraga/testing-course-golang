package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:world", false},
	}

	// create a dummy handler that we'll use to check the context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure that value exists  in the context
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "not present")
		}

		ip, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		// create the handler to test
		handlerToTest := app.addIPToContext(nextHandler)

		req := httptest.NewRequest("GET", "http://trsting", nil)

		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_application_ipFromCOntext(t *testing.T) {
	var ctx = context.Background()

	ctx = context.WithValue(ctx, contextUserKey, "whatever")

	ip := app.ipFromContext(ctx)

	if !strings.EqualFold("whatever", ip) {
		t.Error("wrong value returned from context")
	}

}
