//service_handler_test
package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTextHandler(t *testing.T) {
	handler := new(TextHandler)
	expectedBody := `
John Smith is 22 years old.

Alice Smith is 25 years old.

Bob Baker is 24 years old.
`

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost/hello?say=%s", expectedBody), nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	handler.ServeHTTP(recorder, req)

	switch recorder.Body.String() {
	case expectedBody:
		// body is equal so no need to do anything
	default:
		t.Errorf("Body (%s) did not match expectation (%s).",
			recorder.Body.String(),
			expectedBody)
	}
}

func TestEchosContent(t *testing.T) {
	handler := new(HelloHandler)
	expectedBody := "hellooo!"

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost/hello?say=%s", expectedBody), nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	handler.ServeHTTP(recorder, req)

	switch recorder.Body.String() {
	case expectedBody:
		// body is equal so no need to do anything
	default:
		t.Errorf("Body (%s) did not match expectation (%s).",
			recorder.Body.String(),
			expectedBody)
	}
}

func TestReturns404IfYouSayNothing(t *testing.T) {
	handler := new(HelloHandler)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://example.com/echo?say=Nothing", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Code != 404 {
		t.Errorf("Did not get a 404.")
	}
}

func TestClient(t *testing.T) {
	server := httptest.NewServer(new(HelloHandler))
	defer server.Close()

	// Pretend this is some sort of Go client...
	url := fmt.Sprintf("%s?say=Nothing", server.URL)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Errorf("Error performing request.")
	}

	if resp.StatusCode != 404 {
		t.Errorf("Did not get a 404.")
	}
}
