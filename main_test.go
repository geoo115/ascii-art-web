package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleMainPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDefault)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned unexpected content type: got %v, want %v", contentType, "text/html; charset=utf-8")
	}
}

func TestGenerateAsciiArt(t *testing.T) {
	// Call the generateAsciiArt function with sample input
	text := "h"
	banner := "thinkertoy"
	result := generateAsciiArt(text, banner)

	// Perform assertions on the generated ASCII art
	expectedResult := `     
o    
|    
O--o 
|  | 
o  o 
     
     
`
	if result != expectedResult {
		t.Errorf("generateAsciiArt returned unexpected result: got %v, want %v", result, expectedResult)
	}
}

func TestHandleAsciiArt(t *testing.T) {

	form := strings.NewReader("text=Hello&banner=thinkertoy")
	req, err := http.NewRequest("POST", "/ascii-art", form)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAsciiArt)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned unexpected content type: got %v, want %v", contentType, "text/html; charset=utf-8")
	}

}
