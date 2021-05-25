package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrectStatus(t *testing.T) {

	// Create request
	req, err := http.NewRequest("GET", "http://localhost:5000/api/v1/courses", nil)
	if err != nil {
		t.Fatal(err)
	}
	// create recorder
	newRecorder := httptest.NewRecorder()

	// spawn http.Handler
	handler := createServer()

	// call req with record + handler
	handler.ServeHTTP(newRecorder, req)

	if status := newRecorder.Code; status != http.StatusAccepted {
		t.Errorf("Expect %d got %d\n", http.StatusAccepted, status)
	}
}
