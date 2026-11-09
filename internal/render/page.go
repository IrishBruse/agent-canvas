package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func Page(root *markup.Node) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	fmt.Fprintf(&b, "<title>%s</title>\n<style>\n%s\n</style>\n</head>\n<body>\n<main class=\"page stack\">\n", html.EscapeString(root.Name), canvasCSS)
	b.WriteString(render(root))
	b.WriteString("\n</main>\n</body>\n</html>\n")
	return b.String()
}
