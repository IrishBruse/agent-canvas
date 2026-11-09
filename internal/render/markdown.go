package render

import (
	"bytes"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

func renderMD(n *markup.Node) string {
	var src strings.Builder
	for _, k := range n.Kids {
		if k.El == nil {
			src.WriteString(k.Text)
		}
	}
	body := dedent(src.String())
	var out bytes.Buffer
	if err := markdown.Convert([]byte(body), &out); err != nil {
		return el("pre", "", "", html.EscapeString(body))
	}
	return `<div class="md">` + out.String() + `</div>`
}

func dedent(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	prefix := ""
	seen := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		i := 0
		for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		ws := line[:i]
		if !seen {
			prefix = ws
			seen = true
			continue
		}
		for prefix != "" && !strings.HasPrefix(ws, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines[i] = ""
			continue
		}
		lines[i] = strings.TrimPrefix(line, prefix)
	}
	return strings.Join(lines, "\n")
}
