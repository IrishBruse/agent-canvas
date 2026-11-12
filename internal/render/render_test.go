package render

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func TestExamplesRender(t *testing.T) {
	root := repoRoot(t)
	components, err := canvasFiles(filepath.Join(root, "components"))
	if err != nil {
		t.Fatal(err)
	}
	examples, err := filepath.Glob(filepath.Join(root, "examples", "*.canvas"))
	if err != nil {
		t.Fatal(err)
	}
	if len(components) < 30 {
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
		doc, err := markup.Parse(strings.NewReader(string(raw)))
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		html := Page(doc)
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

func canvasFiles(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".canvas") {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

func TestMarkdown(t *testing.T) {
	doc, err := markup.Parse(strings.NewReader(`
<Stack>
  <markdown>
    # Notes

    Payments **p95** crossed.

    - checkout
      - nested
  </markdown>
</Stack>`))
	if err != nil {
		t.Fatal(err)
	}
	html := Page(doc)
	for _, want := range []string{"<h1>Notes</h1>", "<strong>p95</strong>", "<li>checkout", "<li>nested</li>", `class="md"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %s\n%s", want, html)
		}
	}
	if strings.Contains(html, "<pre>") {
		t.Fatalf("indented markdown rendered as code:\n%s", html)
	}
}

func TestSecondaryButton(t *testing.T) {
	doc, err := markup.Parse(strings.NewReader(`
<Row gap="8">
  <button variant="secondary">Disable</button>
  <button variant="secondary" dropdown="true">Uninstall</button>
</Row>`))
	if err != nil {
		t.Fatal(err)
	}
	html := Page(doc)
	for _, want := range []string{
		`<button class="btn secondary" type="button">Disable</button>`,
		`class="btn-split"`,
		`class="btn-split-sep"`,
		`aria-label="More Actions"`,
		`class="chevron"`,
		`>Uninstall</button>`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %s\n%s", want, html)
		}
	}
}

func TestStackMarkup(t *testing.T) {
	doc, err := markup.Parse(strings.NewReader(`<Stack gap="16"><h1>Release health</h1></Stack>`))
	if err != nil {
		t.Fatal(err)
	}
	html := Page(doc)
	if !strings.Contains(html, `class="stack"`) || !strings.Contains(html, "<h1>Release health</h1>") {
		t.Fatalf("html = %s", html)
	}
}

func TestTable(t *testing.T) {
	doc, err := markup.Parse(strings.NewReader(`
<table striped="true" stickyHeader="true">
  <thead>
    <tr>
      <th>Service</th>
      <th align="right">RPS</th>
    </tr>
  </thead>
  <tbody>
    <tr tone="danger">
      <td>api</td>
      <td>3.2k</td>
    </tr>
  </tbody>
</table>`))
	if err != nil {
		t.Fatal(err)
	}
	html := Page(doc)
	for _, want := range []string{
		`class="table-frame sticky"`,
		`<table class="data striped">`,
		`<th class="num">RPS</th>`,
		`<tr class="tone-danger">`,
		`<td class="num">3.2k</td>`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %s\n%s", want, html)
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
