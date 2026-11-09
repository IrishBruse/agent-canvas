package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExamplesRender(t *testing.T) {
	components, err := filepath.Glob("components/*.canvas")
	if err != nil {
		t.Fatal(err)
	}
	examples, err := filepath.Glob("examples/*.canvas")
	if err != nil {
		t.Fatal(err)
	}
	if len(components) < 60 {
		t.Fatalf("components = %d", len(components))
	}
	if len(examples) < 3 {
		t.Fatalf("examples = %d", len(examples))
	}
	matches := append(components, examples...)
	for _, path := range matches {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		root, err := parse(strings.NewReader(string(raw)))
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		html := page(root)
		if !strings.HasPrefix(html, "<!DOCTYPE html>") {
			t.Fatalf("%s: missing doctype", path)
		}
		if !strings.Contains(html, "<style>") || strings.Contains(html, "stylesheet") {
			t.Fatalf("%s: style is not inlined", path)
		}
		if strings.Count(html, "<html") != 1 {
			t.Fatalf("%s: expected one document", path)
		}
	}
}

func TestMarkdown(t *testing.T) {
	root, err := parse(strings.NewReader(`
<Stack>
  <MD>
    # Notes

    Payments **p95** crossed.

    - checkout
      - nested
  </MD>
</Stack>`))
	if err != nil {
		t.Fatal(err)
	}
	html := page(root)
	for _, want := range []string{"<h1>Notes</h1>", "<strong>p95</strong>", "<li>checkout", "<li>nested</li>", `class="md"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %s\n%s", want, html)
		}
	}
	if strings.Contains(html, "<pre>") {
		t.Fatalf("indented markdown rendered as code:\n%s", html)
	}
}

func TestStackMarkup(t *testing.T) {
	root, err := parse(strings.NewReader(`<Stack gap="16"><H1>Release health</H1></Stack>`))
	if err != nil {
		t.Fatal(err)
	}
	html := page(root)
	if !strings.Contains(html, `class="stack"`) || !strings.Contains(html, "<h1>Release health</h1>") {
		t.Fatalf("html = %s", html)
	}
}
