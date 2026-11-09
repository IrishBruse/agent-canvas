package render

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func el(tag, className, style, inner string) string {
	var b strings.Builder
	b.WriteByte('<')
	b.WriteString(tag)
	if className != "" {
		fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(className))
	}
	if style != "" {
		fmt.Fprintf(&b, ` style="%s"`, html.EscapeString(style))
	}
	b.WriteByte('>')
	b.WriteString(inner)
	fmt.Fprintf(&b, "</%s>", tag)
	return b.String()
}

func join(parts ...string) string {
	var out []string
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}

func gapStyle(n *markup.Node) string {
	if n.Attr("gap") == "" {
		return ""
	}
	return "gap:" + px(n.Attr("gap"))
}

func px(v string) string {
	if _, err := strconv.Atoi(v); err == nil {
		return v + "px"
	}
	return v
}

func first(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func safeColor(v string) string {
	switch v {
	case "gray", "purple", "green", "yellow", "cyan", "pink", "blue", "orange", "red", "accent",
		"success", "danger", "warning", "info", "neutral":
		return v
	default:
		return "gray"
	}
}

func number(v string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return 0
	}
	return f
}
