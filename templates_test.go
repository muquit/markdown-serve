package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderDocumentDoesNotBlockOnHighlightJS(t *testing.T) {
	var output bytes.Buffer
	renderDocument(&output, "example.md", []byte("<p>hello</p>"), "", false, "", false)

	html := output.String()
	headEnd := strings.Index(html, "</head>")
	if headEnd == -1 {
		t.Fatal("rendered document has no closing head tag")
	}
	if strings.Contains(html[:headEnd], "cdnjs.cloudflare.com") {
		t.Fatal("highlight.js CDN resource in head would block initial page load")
	}
	if !strings.Contains(html, "window.addEventListener('load'") {
		t.Fatal("highlight.js is not deferred until after the initial page load")
	}
	if !strings.Contains(html, "github.min.css") {
		t.Fatal("light highlight.js theme is missing")
	}
}

func TestRenderDocumentUsesDarkHighlightTheme(t *testing.T) {
	var output bytes.Buffer
	renderDocument(&output, "example.md", nil, "", false, "", true)

	if !strings.Contains(output.String(), "github-dark.min.css") {
		t.Fatal("dark highlight.js theme is missing")
	}
}

func TestWatchScriptDoesNotOpenPermanentConnection(t *testing.T) {
	var output bytes.Buffer
	renderDocument(&output, "example.md", nil, "", true, "", false)

	html := output.String()
	if strings.Contains(html, "EventSource") {
		t.Fatal("watch mode must not consume a permanent browser connection")
	}
	if !strings.Contains(html, "setInterval(checkForChanges, 1000)") {
		t.Fatal("watch mode revision polling is missing")
	}
}
