package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func renderCard(n *markup.Node) string {
	className := join("card", n.Attr("variant"), n.Attr("size"))
	if className == "card default" {
		className = "card"
	}
	inner := renderKids(n, nil)
	if n.Truthy("collapsible") {
		open := ""
		if n.Truthy("open") || n.Truthy("defaultOpen") {
			open = " open"
		}
		return fmt.Sprintf(`<details class="%s"%s>%s</details>`, html.EscapeString(className), open, inner)
	}
	return el("section", className, "", inner)
}

func renderCardHeader(n *markup.Node) string {
	trailing := n.Child("trailing")
	body := renderKids(n, map[string]bool{"trailing": true})
	if trailing != nil {
		body += el("span", "card-trailing", "", renderKids(trailing, nil))
	}
	return el("div", "card-header", "", body)
}

func renderSection(n *markup.Node) string {
	open := ""
	if n.Truthy("defaultOpen") {
		open = " open"
	}
	var summary strings.Builder
	if lead := n.Child("leading"); lead != nil {
		summary.WriteString(renderKids(lead, nil))
	}
	summary.WriteString(html.EscapeString(n.Attr("title")))
	if n.Attr("count") != "" {
		summary.WriteString(el("span", "section-count", "", html.EscapeString(n.Attr("count"))))
	}
	if trail := n.Child("trailing"); trail != nil {
		summary.WriteString(el("span", "card-trailing", "", renderKids(trail, nil)))
	}
	body := renderKids(n, map[string]bool{"leading": true, "trailing": true})
	return fmt.Sprintf(`<details class="section"%s><summary>%s</summary><div class="section-body">%s</div></details>`, open, summary.String(), body)
}

func renderCallout(n *markup.Node) string {
	className := "callout"
	if tone := n.Attr("tone"); tone != "" {
		className += " " + tone
	}
	title := ""
	if n.Attr("title") != "" {
		title = el("div", "callout-title", "", html.EscapeString(n.Attr("title")))
	}
	return el("aside", className, "", title+el("p", "text", "", renderKids(n, nil)))
}

func renderStat(n *markup.Node) string {
	className := "stat"
	if tone := n.Attr("tone"); tone != "" {
		className += " " + tone
	}
	return el("div", className, "", el("div", "stat-value", "", html.EscapeString(n.Attr("value")))+el("div", "stat-label", "", html.EscapeString(n.Attr("label"))))
}
