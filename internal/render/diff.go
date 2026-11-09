package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func renderDiffStats(n *markup.Node) string {
	return fmt.Sprintf(`<span class="diff-stats"><span class="add">+%s</span><span class="del">-%s</span></span>`, html.EscapeString(n.Attr("additions")), html.EscapeString(n.Attr("deletions")))
}

func renderDiffLine(kind, content, line string) string {
	mark := " "
	switch kind {
	case "added":
		mark = "+"
	case "removed":
		mark = "-"
	}
	return fmt.Sprintf(`<div class="diff-line %s"><span class="ln">%s</span><span class="mark">%s</span><span>%s</span></div>`, html.EscapeString(kind), html.EscapeString(line), mark, html.EscapeString(content))
}

func renderDiffView(n *markup.Node) string {
	var b strings.Builder
	b.WriteString(`<div class="diff">`)
	lines := n.Child("lines")
	if lines == nil {
		lines = n
	}
	for _, line := range lines.Elements("line") {
		b.WriteString(renderDiffLine(line.Attr("type"), line.Attr("content"), line.Attr("lineNumber")))
	}
	for _, line := range lines.Elements("DiffLineData") {
		b.WriteString(render(line))
	}
	b.WriteString(`</div>`)
	return b.String()
}
