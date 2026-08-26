package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWatchHandlerReturnsAndReportsRevision(t *testing.T) {
	b := newBroadcaster()
	handler := watchHandler(b)

	request := httptest.NewRequest("GET", "/events", nil)
	response := httptest.NewRecorder()
	handler(response, request)
	if response.Body.String() != "0\n" {
		t.Fatalf("initial revision = %q, want %q", response.Body.String(), "0\\n")
	}
	if strings.Contains(response.Header().Get("Content-Type"), "text/event-stream") {
		t.Fatal("watch handler must not keep an SSE connection open")
	}

	b.broadcast()
	response = httptest.NewRecorder()
	handler(response, request)
	if response.Body.String() != "1\n" {
		t.Fatalf("updated revision = %q, want %q", response.Body.String(), "1\\n")
	}
}
