package render

import (
	"strconv"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func rowClass(n *markup.Node) string {
	className := "row"
	if n.Truthy("wrap") {
		className += " wrap"
	}
	if n.Attr("justify") == "space-between" {
		className += " spread"
	}
	if n.Attr("align") == "start" {
		className += " top"
	}
	return className
}

func gridStyle(n *markup.Node) string {
	cols := n.Attr("columns")
	style := ""
	if cols != "" {
		if _, err := strconv.Atoi(cols); err == nil {
			style = "--cols:" + cols
		} else {
			style = "grid-template-columns:" + cols
		}
	}
	if n.Attr("gap") != "" {
		if style != "" {
			style += ";"
		}
		style += "gap:" + px(n.Attr("gap"))
	}
	return style
}

func textClass(n *markup.Node) string {
	className := "text"
	switch n.Attr("tone") {
	case "secondary", "tertiary", "quaternary":
		className += " " + n.Attr("tone")
	}
	if n.Attr("size") == "small" {
		className += " small"
	}
	switch n.Attr("weight") {
	case "medium", "semibold", "bold":
		className += " " + n.Attr("weight")
	}
	if n.Truthy("italic") {
		className += " italic"
	}
	if n.Attr("truncate") != "" {
		className += " truncate"
	}
	return className
}
