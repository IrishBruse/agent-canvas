package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func renderTodoItem(n *markup.Node) string {
	return fmt.Sprintf(`<li class="todo-item %s"><span class="todo-mark"></span>%s</li>`, html.EscapeString(n.Attr("status")), html.EscapeString(first(n.Attr("content"), n.Text())))
}

func renderTodoList(n *markup.Node, card bool) string {
	var items strings.Builder
	items.WriteString(`<ul class="todo">`)
	host := n.Child("todos")
	if host == nil {
		host = n
	}
	done := 0
	total := 0
	for _, item := range append(host.Elements("todo"), host.Elements("TodoItem")...) {
		total++
		if item.Attr("status") == "completed" {
			done++
		}
		items.WriteString(renderTodoItem(item))
	}
	items.WriteString(`</ul>`)
	if !card {
		return items.String()
	}
	open := ""
	if n.Truthy("defaultExpanded") {
		open = " open"
	}
	return fmt.Sprintf(`<details class="todo-card"%s><summary>%d of %d Done</summary><div class="card-body">%s</div></details>`, open, done, total, items.String())
}
