package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func render(n *markup.Node) string {
	if n == nil {
		return ""
	}
	switch n.Name {
	case "canvas", "Stack":
		return el("div", join("stack"), gapStyle(n), renderKids(n, nil))
	case "Row":
		return el("div", rowClass(n), gapStyle(n), renderKids(n, nil))
	case "Grid":
		return el("div", "grid", gridStyle(n), renderKids(n, nil))
	case "Divider":
		return "<hr>"
	case "Spacer":
		return `<span class="spacer"></span>`
	case "Text":
		return el("p", textClass(n), "", renderKids(n, nil))
	case "H1", "H2", "H3":
		return el(strings.ToLower(n.Name), "", "", renderKids(n, nil))
	case "MD":
		return renderMD(n)
	case "Code":
		return el("code", "", "", renderKids(n, nil))
	case "Link":
		return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(n.Attr("href")), renderKids(n, nil))
	case "Card":
		return renderCard(n)
	case "CardHeader":
		return renderCardHeader(n)
	case "CardBody":
		return el("div", "card-body", "", renderKids(n, nil))
	case "CanvasChevron":
		mark := "▸"
		if n.Truthy("expanded") {
			mark = "▾"
		}
		return el("span", "text", "", html.EscapeString(mark))
	case "CollapsibleSection":
		return renderSection(n)
	case "Table":
		return renderTable(n)
	case "Button":
		return renderButton(n)
	case "Pill":
		return renderPill(n)
	case "IconButton":
		return renderIconButton(n)
	case "Stat":
		return renderStat(n)
	case "Callout":
		return renderCallout(n)
	case "CalloutToneIcon":
		return fmt.Sprintf(`<span class="swatch c-%s" title="%s"></span>`, safeColor(n.Attr("tone")), html.EscapeString(n.Attr("color")))
	case "DiffStats":
		return renderDiffStats(n)
	case "DiffLineData":
		return renderDiffLine(n.Attr("type"), n.Attr("content"), n.Attr("lineNumber"))
	case "DiffView":
		return renderDiffView(n)
	case "TodoItem":
		return renderTodoItem(n)
	case "TodoList":
		return renderTodoList(n, false)
	case "TodoListCard":
		return renderTodoList(n, true)
	case "TextInput":
		return renderInput(n, false)
	case "TextArea":
		return renderInput(n, true)
	case "Checkbox":
		return renderCheckbox(n)
	case "Toggle":
		return renderToggle(n)
	case "SelectOption":
		return renderOption(n)
	case "Select":
		return renderSelect(n)
	case "ChartDataPoint":
		return renderPie(&markup.Node{Name: "PieChart", Kids: []markup.Kid{{El: n}}})
	case "ChartSeries", "ChartReferenceLine", "UsageBarSegment":
		return renderRecord(n)
	case "BarChart":
		return renderBars(n)
	case "LineChart":
		return renderLine(n)
	case "PieChart":
		return renderPie(n)
	case "UsageBar":
		return renderUsage(n)
	case "Swatch":
		return fmt.Sprintf(`<span class="swatch c-%s"></span>`, safeColor(n.Attr("color")))
	default:
		return renderRecord(n)
	}
}

func renderKids(n *markup.Node, skip map[string]bool) string {
	var b strings.Builder
	for _, k := range n.Kids {
		if k.El != nil {
			if skip[k.El.Name] {
				continue
			}
			b.WriteString(render(k.El))
			continue
		}
		if k.Text != "" {
			b.WriteString(html.EscapeString(k.Text))
		}
	}
	return b.String()
}

func renderRecord(n *markup.Node) string {
	var body strings.Builder
	if len(n.Attrs) > 0 {
		body.WriteString(`<p class="kv">`)
		for i, a := range n.Attrs {
			if i > 0 {
				body.WriteString(" · ")
			}
			fmt.Fprintf(&body, "%s=%s", html.EscapeString(a.Name), html.EscapeString(a.Value))
		}
		body.WriteString(`</p>`)
	}
	body.WriteString(renderKids(n, nil))
	return el("section", "card", "", el("div", "card-header", "", html.EscapeString(n.Name))+el("div", "card-body", "", body.String()))
}
