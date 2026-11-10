package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func renderTable(n *markup.Node) string {
	className := "data"
	if n.Truthy("striped") {
		className += " striped"
	}
	aligns := columnAligns(n)
	var b strings.Builder
	fmt.Fprintf(&b, `<table class="%s">`, className)
	if caption := n.Child("caption"); caption != nil {
		b.WriteString("<caption>" + renderKids(caption, nil) + "</caption>")
	}
	rows := 0
	for _, k := range n.Kids {
		if k.El == nil {
			continue
		}
		switch k.El.Name {
		case "thead", "tbody", "tfoot":
			section, count := sectionHTML(k.El, aligns)
			b.WriteString(section)
			rows += count
		case "tr":
			writeRow(&b, k.El, aligns)
			rows++
		}
	}
	if rows == 0 {
		msg := "Empty"
		if n.Attr("empty") != "" {
			msg = n.Attr("empty")
		}
		fmt.Fprintf(&b, `<tbody><tr><td colspan="%d">%s</td></tr></tbody>`, max(1, len(aligns)), html.EscapeString(msg))
	}
	b.WriteString("</table>")
	if n.Attr("framed") == "false" {
		return b.String()
	}
	frame := "table-frame"
	if n.Truthy("stickyHeader") {
		frame += " sticky"
	}
	return el("div", frame, "", b.String())
}

func columnAligns(n *markup.Node) []string {
	row := (*markup.Node)(nil)
	if head := n.Child("thead"); head != nil {
		row = head.Child("tr")
	}
	if row == nil {
		row = n.Child("tr")
	}
	if row == nil {
		return nil
	}
	var out []string
	for _, cell := range rowCells(row) {
		out = append(out, cell.Attr("align"))
	}
	return out
}

func sectionHTML(n *markup.Node, aligns []string) (string, int) {
	var inner strings.Builder
	count := 0
	for _, tr := range n.Elements("tr") {
		writeRow(&inner, tr, aligns)
		count++
	}
	if count == 0 {
		return "", 0
	}
	return "<" + n.Name + ">" + inner.String() + "</" + n.Name + ">", count
}

func writeRow(b *strings.Builder, tr *markup.Node, aligns []string) {
	b.WriteString("<tr" + toneClass(tr.Attr("tone")) + ">")
	for i, cell := range rowCells(tr) {
		align := cell.Attr("align")
		if align == "" && i < len(aligns) {
			align = aligns[i]
		}
		fmt.Fprintf(b, "<%s%s>%s</%s>", cell.Name, alignAttr(align), renderKids(cell, nil), cell.Name)
	}
	b.WriteString("</tr>")
}

func rowCells(tr *markup.Node) []*markup.Node {
	var out []*markup.Node
	for _, k := range tr.Kids {
		if k.El != nil && (k.El.Name == "th" || k.El.Name == "td") {
			out = append(out, k.El)
		}
	}
	return out
}

func toneClass(v string) string {
	switch v {
	case "success", "danger", "warning", "info", "neutral":
		return ` class="tone-` + v + `"`
	default:
		return ""
	}
}

func alignAttr(align string) string {
	switch align {
	case "right":
		return ` class="num"`
	case "center":
		return ` style="text-align:center"`
	default:
		return ""
	}
}
