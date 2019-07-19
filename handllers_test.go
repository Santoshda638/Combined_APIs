package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var ctr = NewHandller()

func TestJokeHandller(t *testing.T) {

	// Create a request for handler.
	req, err := http.NewRequest("GET", "/GetNewJoke", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctr.GetNewJokeHandller)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"type":"success"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func BenchmarkHandller(b *testing.B) {
	req, err := http.NewRequest("GET", "/GetNewJoke", nil)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		rw := httptest.NewRecorder()
		ctr.GetNewJokeHandller(rw, req)
		//fmt.Println(rw.Body.String())
	}
}
