package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func renderTable(n *markup.Node) string {
	headers := texts(n.Child("headers"), "header")
	aligns := texts(n.Child("columnAlign"), "align")
	tones := texts(n.Child("rowTone"), "tone")
	var rows [][]string
	if body := n.Child("rows"); body != nil {
		for _, row := range body.Elements("row") {
			rows = append(rows, texts(row, "cell"))
		}
	}
	className := "data"
	if n.Truthy("striped") {
		className += " striped"
	}
	var b strings.Builder
	b.WriteString("<table class=\"" + className + "\"><thead><tr>")
	for i, h := range headers {
		b.WriteString("<th" + alignAttr(aligns, i) + ">" + html.EscapeString(h) + "</th>")
	}
	b.WriteString("</tr></thead><tbody>")
	if len(rows) == 0 {
		msg := "Empty"
		if empty := n.Child("emptyMessage"); empty != nil {
			msg = empty.Text()
		}
		fmt.Fprintf(&b, `<tr><td colspan="%d">%s</td></tr>`, max(1, len(headers)), html.EscapeString(msg))
	}
	for i, row := range rows {
		tone := ""
		if i < len(tones) && tones[i] != "" {
			tone = ` class="tone-` + html.EscapeString(tones[i]) + `"`
		}
		b.WriteString("<tr" + tone + ">")
		for c, cell := range row {
			b.WriteString("<td" + alignAttr(aligns, c) + ">" + html.EscapeString(cell) + "</td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table>")
	frame := "table-frame"
	if n.Truthy("stickyHeader") {
		frame += " sticky"
	}
	if n.Attr("framed") == "false" {
		return b.String()
	}
	return el("div", frame, "", b.String())
}

func alignAttr(aligns []string, i int) string {
	if i < len(aligns) && aligns[i] == "right" {
		return ` class="num"`
	}
	if i < len(aligns) && aligns[i] == "center" {
		return ` style="text-align:center"`
	}
	return ""
}

func texts(n *markup.Node, name string) []string {
	if n == nil {
		return nil
	}
	var out []string
	for _, el := range n.Elements(name) {
		out = append(out, el.Text())
	}
	return out
}
