package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

func renderButton(n *markup.Node) string {
	className := "btn"
	if v := n.Attr("variant"); v != "" && v != "secondary" {
		className += " " + v
	}
	disabled := ""
	if n.Truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<button class="%s" type="button"%s>%s</button>`, className, disabled, renderKids(n, nil))
}

func renderPill(n *markup.Node) string {
	className := "pill"
	if n.Truthy("active") {
		className += " active"
	}
	if n.Attr("size") == "sm" {
		className += " sm"
	}
	lead := ""
	if child := n.Child("leadingContent"); child != nil {
		lead = renderKids(child, nil)
	}
	hint := ""
	if n.Attr("keyboardHint") != "" {
		hint = "<kbd>" + html.EscapeString(n.Attr("keyboardHint")) + "</kbd>"
	}
	inner := lead + renderKids(n, map[string]bool{"leadingContent": true}) + hint
	if n.Attr("onClick") != "" {
		return fmt.Sprintf(`<button class="%s" type="button" title="%s">%s</button>`, className, html.EscapeString(n.Attr("title")), inner)
	}
	return el("span", className, "", inner)
}

func renderIconButton(n *markup.Node) string {
	className := "icon-btn"
	if n.Attr("variant") == "circle" {
		className += " circle"
	}
	if n.Attr("size") == "sm" {
		className += " sm"
	}
	disabled := ""
	if n.Truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<button class="%s" type="button" title="%s"%s>%s</button>`, className, html.EscapeString(n.Attr("title")), disabled, renderKids(n, nil))
}

func renderInput(n *markup.Node, area bool) string {
	disabled := ""
	if n.Truthy("disabled") {
		disabled = " disabled"
	}
	if area {
		return fmt.Sprintf(`<textarea class="input" rows="%s" placeholder="%s"%s>%s</textarea>`, html.EscapeString(first(n.Attr("rows"), "3")), html.EscapeString(n.Attr("placeholder")), disabled, html.EscapeString(n.Attr("value")))
	}
	return fmt.Sprintf(`<input class="input" type="text" value="%s" placeholder="%s"%s>`, html.EscapeString(n.Attr("value")), html.EscapeString(n.Attr("placeholder")), disabled)
}

func renderCheckbox(n *markup.Node) string {
	checked := ""
	if n.Truthy("checked") {
		checked = " checked"
	}
	disabled := ""
	if n.Truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<label class="check"><input type="checkbox"%s%s> %s</label>`, checked, disabled, html.EscapeString(n.Attr("label")))
}

func renderToggle(n *markup.Node) string {
	checked := ""
	if n.Truthy("checked") {
		checked = " checked"
	}
	className := "toggle"
	if n.Attr("size") == "md" {
		className += " md"
	}
	return fmt.Sprintf(`<label class="%s"><input type="checkbox"%s><span class="toggle-track"></span></label>`, className, checked)
}

func renderOption(n *markup.Node) string {
	disabled := ""
	if n.Truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<option value="%s"%s>%s</option>`, html.EscapeString(n.Attr("value")), disabled, html.EscapeString(n.Attr("label")))
}

func renderSelect(n *markup.Node) string {
	disabled := ""
	if n.Truthy("disabled") {
		disabled = " disabled"
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<select class="input"%s>`, disabled)
	if n.Attr("placeholder") != "" {
		fmt.Fprintf(&b, `<option value="">%s</option>`, html.EscapeString(n.Attr("placeholder")))
	}
	host := n.Child("options")
	if host == nil {
		host = n
	}
	for _, opt := range append(host.Elements("option"), host.Elements("SelectOption")...) {
		selected := ""
		if opt.Attr("value") == n.Attr("value") {
			selected = " selected"
		}
		disabled := ""
		if opt.Truthy("disabled") {
			disabled = " disabled"
		}
		fmt.Fprintf(&b, `<option value="%s"%s%s>%s</option>`, html.EscapeString(opt.Attr("value")), selected, disabled, html.EscapeString(opt.Attr("label")))
	}
	b.WriteString(`</select>`)
	return b.String()
}
