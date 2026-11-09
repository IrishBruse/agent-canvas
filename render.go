package main

import (
	"bytes"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

func page(root *node) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	fmt.Fprintf(&b, "<title>%s</title>\n<style>\n%s\n</style>\n</head>\n<body>\n<main class=\"page stack\">\n", html.EscapeString(root.Name), canvasCSS)
	b.WriteString(render(root))
	b.WriteString("\n</main>\n</body>\n</html>\n")
	return b.String()
}

func render(n *node) string {
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
		return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(n.attr("href")), renderKids(n, nil))
	case "Card":
		return renderCard(n)
	case "CardHeader":
		return renderCardHeader(n)
	case "CardBody":
		return el("div", "card-body", "", renderKids(n, nil))
	case "CanvasChevron":
		mark := "▸"
		if n.truthy("expanded") {
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
		return fmt.Sprintf(`<span class="swatch c-%s" title="%s"></span>`, safeColor(n.attr("tone")), html.EscapeString(n.attr("color")))
	case "DiffStats":
		return renderDiffStats(n)
	case "DiffLineData":
		return renderDiffLine(n.attr("type"), n.attr("content"), n.attr("lineNumber"))
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
		return renderPie(&node{Name: "PieChart", Kids: []kid{{El: n}}})
	case "ChartSeries", "ChartReferenceLine", "UsageBarSegment":
		return renderRecord(n)
	case "BarChart":
		return renderBars(n, false)
	case "LineChart":
		return renderLine(n)
	case "PieChart":
		return renderPie(n)
	case "UsageBar":
		return renderUsage(n)
	case "Swatch":
		return fmt.Sprintf(`<span class="swatch c-%s"></span>`, safeColor(n.attr("color")))
	default:
		return renderRecord(n)
	}
}

func renderMD(n *node) string {
	var src strings.Builder
	for _, k := range n.Kids {
		if k.El == nil {
			src.WriteString(k.Text)
		}
	}
	body := dedent(src.String())
	var out bytes.Buffer
	if err := markdown.Convert([]byte(body), &out); err != nil {
		return el("pre", "", "", html.EscapeString(body))
	}
	return `<div class="md">` + out.String() + `</div>`
}

func dedent(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	prefix := ""
	seen := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		i := 0
		for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		ws := line[:i]
		if !seen {
			prefix = ws
			seen = true
			continue
		}
		for prefix != "" && !strings.HasPrefix(ws, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines[i] = ""
			continue
		}
		lines[i] = strings.TrimPrefix(line, prefix)
	}
	return strings.Join(lines, "\n")
}

func renderKids(n *node, skip map[string]bool) string {
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

func gapStyle(n *node) string {
	if n.attr("gap") == "" {
		return ""
	}
	return "gap:" + px(n.attr("gap"))
}

func px(v string) string {
	if _, err := strconv.Atoi(v); err == nil {
		return v + "px"
	}
	return v
}

func rowClass(n *node) string {
	className := "row"
	if n.truthy("wrap") {
		className += " wrap"
	}
	if n.attr("justify") == "space-between" {
		className += " spread"
	}
	if n.attr("align") == "start" {
		className += " top"
	}
	return className
}

func gridStyle(n *node) string {
	cols := n.attr("columns")
	style := ""
	if cols != "" {
		if _, err := strconv.Atoi(cols); err == nil {
			style = "--cols:" + cols
		} else {
			style = "grid-template-columns:" + cols
		}
	}
	if n.attr("gap") != "" {
		if style != "" {
			style += ";"
		}
		style += "gap:" + px(n.attr("gap"))
	}
	return style
}

func textClass(n *node) string {
	className := "text"
	switch n.attr("tone") {
	case "secondary", "tertiary", "quaternary":
		className += " " + n.attr("tone")
	}
	if n.attr("size") == "small" {
		className += " small"
	}
	switch n.attr("weight") {
	case "medium", "semibold", "bold":
		className += " " + n.attr("weight")
	}
	if n.truthy("italic") {
		className += " italic"
	}
	if n.attr("truncate") != "" {
		className += " truncate"
	}
	return className
}

func renderCard(n *node) string {
	className := join("card", n.attr("variant"), n.attr("size"))
	if className == "card default" {
		className = "card"
	}
	inner := renderKids(n, nil)
	if n.truthy("collapsible") {
		open := ""
		if n.truthy("open") || n.truthy("defaultOpen") {
			open = " open"
		}
		return fmt.Sprintf(`<details class="%s"%s>%s</details>`, html.EscapeString(className), open, inner)
	}
	return el("section", className, "", inner)
}

func renderCardHeader(n *node) string {
	trailing := n.child("trailing")
	body := renderKids(n, map[string]bool{"trailing": true})
	if trailing != nil {
		body += el("span", "card-trailing", "", renderKids(trailing, nil))
	}
	return el("div", "card-header", "", body)
}

func renderSection(n *node) string {
	open := ""
	if n.truthy("defaultOpen") {
		open = " open"
	}
	var summary strings.Builder
	if lead := n.child("leading"); lead != nil {
		summary.WriteString(renderKids(lead, nil))
	}
	summary.WriteString(html.EscapeString(n.attr("title")))
	if n.attr("count") != "" {
		summary.WriteString(el("span", "section-count", "", html.EscapeString(n.attr("count"))))
	}
	if trail := n.child("trailing"); trail != nil {
		summary.WriteString(el("span", "card-trailing", "", renderKids(trail, nil)))
	}
	body := renderKids(n, map[string]bool{"leading": true, "trailing": true})
	return fmt.Sprintf(`<details class="section"%s><summary>%s</summary><div class="section-body">%s</div></details>`, open, summary.String(), body)
}

func renderTable(n *node) string {
	headers := texts(n.child("headers"), "header")
	aligns := texts(n.child("columnAlign"), "align")
	tones := texts(n.child("rowTone"), "tone")
	var rows [][]string
	if body := n.child("rows"); body != nil {
		for _, row := range body.elements("row") {
			rows = append(rows, texts(row, "cell"))
		}
	}
	className := "data"
	if n.truthy("striped") {
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
		if empty := n.child("emptyMessage"); empty != nil {
			msg = empty.text()
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
	if n.truthy("stickyHeader") {
		frame += " sticky"
	}
	if n.attr("framed") == "false" {
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

func texts(n *node, name string) []string {
	if n == nil {
		return nil
	}
	var out []string
	for _, el := range n.elements(name) {
		out = append(out, el.text())
	}
	return out
}

func renderButton(n *node) string {
	className := "btn"
	if v := n.attr("variant"); v != "" && v != "secondary" {
		className += " " + v
	}
	disabled := ""
	if n.truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<button class="%s" type="button"%s>%s</button>`, className, disabled, renderKids(n, nil))
}

func renderPill(n *node) string {
	className := "pill"
	if n.truthy("active") {
		className += " active"
	}
	if n.attr("size") == "sm" {
		className += " sm"
	}
	lead := ""
	if child := n.child("leadingContent"); child != nil {
		lead = renderKids(child, nil)
	}
	hint := ""
	if n.attr("keyboardHint") != "" {
		hint = "<kbd>" + html.EscapeString(n.attr("keyboardHint")) + "</kbd>"
	}
	inner := lead + renderKids(n, map[string]bool{"leadingContent": true}) + hint
	if n.attr("onClick") != "" {
		return fmt.Sprintf(`<button class="%s" type="button" title="%s">%s</button>`, className, html.EscapeString(n.attr("title")), inner)
	}
	return el("span", className, "", inner)
}

func renderIconButton(n *node) string {
	className := "icon-btn"
	if n.attr("variant") == "circle" {
		className += " circle"
	}
	if n.attr("size") == "sm" {
		className += " sm"
	}
	disabled := ""
	if n.truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<button class="%s" type="button" title="%s"%s>%s</button>`, className, html.EscapeString(n.attr("title")), disabled, renderKids(n, nil))
}

func renderStat(n *node) string {
	className := "stat"
	if tone := n.attr("tone"); tone != "" {
		className += " " + tone
	}
	return el("div", className, "", el("div", "stat-value", "", html.EscapeString(n.attr("value")))+el("div", "stat-label", "", html.EscapeString(n.attr("label"))))
}

func renderCallout(n *node) string {
	className := "callout"
	if tone := n.attr("tone"); tone != "" {
		className += " " + tone
	}
	title := ""
	if n.attr("title") != "" {
		title = el("div", "callout-title", "", html.EscapeString(n.attr("title")))
	}
	return el("aside", className, "", title+el("p", "text", "", renderKids(n, nil)))
}

func renderDiffStats(n *node) string {
	return fmt.Sprintf(`<span class="diff-stats"><span class="add">+%s</span><span class="del">-%s</span></span>`, html.EscapeString(n.attr("additions")), html.EscapeString(n.attr("deletions")))
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

func renderDiffView(n *node) string {
	var b strings.Builder
	b.WriteString(`<div class="diff">`)
	lines := n.child("lines")
	if lines == nil {
		lines = n
	}
	for _, line := range lines.elements("line") {
		b.WriteString(renderDiffLine(line.attr("type"), line.attr("content"), line.attr("lineNumber")))
	}
	for _, line := range lines.elements("DiffLineData") {
		b.WriteString(render(line))
	}
	b.WriteString(`</div>`)
	return b.String()
}

func renderTodoItem(n *node) string {
	return fmt.Sprintf(`<li class="todo-item %s"><span class="todo-mark"></span>%s</li>`, html.EscapeString(n.attr("status")), html.EscapeString(first(n.attr("content"), n.text())))
}

func renderTodoList(n *node, card bool) string {
	var items strings.Builder
	items.WriteString(`<ul class="todo">`)
	host := n.child("todos")
	if host == nil {
		host = n
	}
	done := 0
	total := 0
	for _, item := range append(host.elements("todo"), host.elements("TodoItem")...) {
		total++
		if item.attr("status") == "completed" {
			done++
		}
		items.WriteString(renderTodoItem(item))
	}
	items.WriteString(`</ul>`)
	if !card {
		return items.String()
	}
	open := ""
	if n.truthy("defaultExpanded") {
		open = " open"
	}
	return fmt.Sprintf(`<details class="todo-card"%s><summary>%d of %d Done</summary><div class="card-body">%s</div></details>`, open, done, total, items.String())
}

func renderInput(n *node, area bool) string {
	disabled := ""
	if n.truthy("disabled") {
		disabled = " disabled"
	}
	if area {
		return fmt.Sprintf(`<textarea class="input" rows="%s" placeholder="%s"%s>%s</textarea>`, html.EscapeString(first(n.attr("rows"), "3")), html.EscapeString(n.attr("placeholder")), disabled, html.EscapeString(n.attr("value")))
	}
	return fmt.Sprintf(`<input class="input" type="text" value="%s" placeholder="%s"%s>`, html.EscapeString(n.attr("value")), html.EscapeString(n.attr("placeholder")), disabled)
}

func renderCheckbox(n *node) string {
	checked := ""
	if n.truthy("checked") {
		checked = " checked"
	}
	disabled := ""
	if n.truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<label class="check"><input type="checkbox"%s%s> %s</label>`, checked, disabled, html.EscapeString(n.attr("label")))
}

func renderToggle(n *node) string {
	checked := ""
	if n.truthy("checked") {
		checked = " checked"
	}
	className := "toggle"
	if n.attr("size") == "md" {
		className += " md"
	}
	return fmt.Sprintf(`<label class="%s"><input type="checkbox"%s><span class="toggle-track"></span></label>`, className, checked)
}

func renderOption(n *node) string {
	disabled := ""
	if n.truthy("disabled") {
		disabled = " disabled"
	}
	return fmt.Sprintf(`<option value="%s"%s>%s</option>`, html.EscapeString(n.attr("value")), disabled, html.EscapeString(n.attr("label")))
}

func renderSelect(n *node) string {
	disabled := ""
	if n.truthy("disabled") {
		disabled = " disabled"
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<select class="input"%s>`, disabled)
	if n.attr("placeholder") != "" {
		fmt.Fprintf(&b, `<option value="">%s</option>`, html.EscapeString(n.attr("placeholder")))
	}
	host := n.child("options")
	if host == nil {
		host = n
	}
	for _, opt := range append(host.elements("option"), host.elements("SelectOption")...) {
		selected := ""
		if opt.attr("value") == n.attr("value") {
			selected = " selected"
		}
		disabled := ""
		if opt.truthy("disabled") {
			disabled = " disabled"
		}
		fmt.Fprintf(&b, `<option value="%s"%s%s>%s</option>`, html.EscapeString(opt.attr("value")), selected, disabled, html.EscapeString(opt.attr("label")))
	}
	b.WriteString(`</select>`)
	return b.String()
}

type series struct {
	name  string
	tone  string
	value []float64
}

func readSeries(n *node) []series {
	host := n.child("series")
	if host == nil {
		host = n
	}
	var out []series
	for _, s := range append(host.elements("series"), host.elements("ChartSeries")...) {
		item := series{name: s.attr("name"), tone: s.attr("tone")}
		for _, v := range s.elements("value") {
			item.value = append(item.value, number(v.text()))
		}
		if data := s.child("data"); data != nil {
			for _, v := range data.elements("value") {
				item.value = append(item.value, number(v.text()))
			}
		}
		out = append(out, item)
	}
	return out
}

func renderBars(n *node, horizontal bool) string {
	_ = horizontal
	cats := texts(n.child("categories"), "category")
	series := readSeries(n)
	maxV := 1.0
	for _, s := range series {
		for _, v := range s.value {
			if v > maxV {
				maxV = v
			}
		}
	}
	var b strings.Builder
	b.WriteString(`<figure class="diagram"><div class="chart-frame"><div class="bars">`)
	count := len(cats)
	if count == 0 && len(series) > 0 {
		count = len(series[0].value)
	}
	for i := 0; i < count; i++ {
		label := ""
		if i < len(cats) {
			label = cats[i]
		}
		b.WriteString(`<div class="bar-col"><div class="bar-stack">`)
		for _, s := range series {
			v := 0.0
			if i < len(s.value) {
				v = s.value[i]
			}
			h := v / maxV * 100
			fmt.Fprintf(&b, `<div class="bar %s" style="height:%.1f%%" title="%s"></div>`, html.EscapeString(s.tone), h, html.EscapeString(s.name))
		}
		fmt.Fprintf(&b, `</div><div class="bar-label">%s</div></div>`, html.EscapeString(label))
	}
	b.WriteString(`</div></div>`)
	b.WriteString(legend(series))
	b.WriteString(`</figure>`)
	return b.String()
}

func renderLine(n *node) string {
	series := readSeries(n)
	cats := texts(n.child("categories"), "category")
	maxV := 1.0
	minV := 0.0
	if !n.truthy("beginAtZero") && n.attr("beginAtZero") == "false" {
		minV = maxV
		for _, s := range series {
			for _, v := range s.value {
				if v < minV {
					minV = v
				}
				if v > maxV {
					maxV = v
				}
			}
		}
	} else {
		for _, s := range series {
			for _, v := range s.value {
				if v > maxV {
					maxV = v
				}
			}
		}
	}
	colors := []string{"#4ca0f5", "#3fb950", "#f85149", "#e3814b", "#a371f7"}
	var b strings.Builder
	b.WriteString(`<figure class="diagram"><div class="chart-frame"><svg class="line-chart" viewBox="0 0 100 40" preserveAspectRatio="none">`)
	for i, s := range series {
		if len(s.value) == 0 {
			continue
		}
		var pts []string
		for j, v := range s.value {
			x := 0.0
			if len(s.value) > 1 {
				x = float64(j) / float64(len(s.value)-1) * 100
			}
			span := maxV - minV
			if span == 0 {
				span = 1
			}
			y := 38 - ((v-minV)/span)*34
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", x, y))
		}
		color := colors[i%len(colors)]
		fmt.Fprintf(&b, `<polyline fill="none" stroke="%s" stroke-width="1.5" points="%s"/>`, color, strings.Join(pts, " "))
	}
	b.WriteString(`</svg></div>`)
	if len(cats) > 0 {
		fmt.Fprintf(&b, `<p class="caption">%s</p>`, html.EscapeString(strings.Join(cats, " · ")))
	}
	b.WriteString(legend(series))
	b.WriteString(`</figure>`)
	return b.String()
}

func renderPie(n *node) string {
	var points []*node
	if data := n.child("data"); data != nil {
		points = append(data.elements("point"), data.elements("ChartDataPoint")...)
	}
	if n.Name == "ChartDataPoint" || n.attr("label") != "" {
		points = []*node{n}
	}
	total := 0.0
	for _, p := range points {
		total += number(first(p.attr("value"), p.text()))
	}
	if total == 0 {
		total = 1
	}
	colors := map[string]string{
		"success": "#3fb950",
		"danger":  "#f85149",
		"warning": "#e3814b",
		"info":    "#4ca0f5",
		"neutral": "#6e7681",
	}
	fallback := []string{"#4ca0f5", "#a371f7", "#3fb950", "#e3814b", "#db61a2"}
	var stops []string
	acc := 0.0
	var legend strings.Builder
	legend.WriteString(`<div class="legend">`)
	for i, p := range points {
		v := number(first(p.attr("value"), p.text()))
		start := acc / total * 100
		acc += v
		end := acc / total * 100
		color := colors[p.attr("tone")]
		if color == "" {
			color = fallback[i%len(fallback)]
		}
		stops = append(stops, fmt.Sprintf("%s %.1f%% %.1f%%", color, start, end))
		fmt.Fprintf(&legend, `<span>%s %s</span>`, html.EscapeString(p.attr("label")), html.EscapeString(first(p.attr("value"), p.text())))
	}
	legend.WriteString(`</div>`)
	className := "pie"
	if n.truthy("donut") {
		className += " donut"
	}
	size := first(n.attr("size"), "200")
	style := fmt.Sprintf("width:%spx;height:%spx;background:conic-gradient(%s)", size, size, strings.Join(stops, ","))
	return el("figure", "diagram", "", el("div", className, style, "")+legend.String())
}

func renderUsage(n *node) string {
	host := n.child("segments")
	if host == nil {
		host = n
	}
	segs := append(host.elements("segment"), host.elements("UsageBarSegment")...)
	total := number(n.attr("total"))
	if total <= 0 {
		total = 1
	}
	var b strings.Builder
	b.WriteString(`<div class="usage"><div class="usage-labels">`)
	fmt.Fprintf(&b, `<span>%s</span><span>%s</span></div><div class="usage-track">`, html.EscapeString(n.attr("topLeftLabel")), html.EscapeString(n.attr("topRightLabel")))
	used := 0.0
	for _, s := range segs {
		v := number(s.attr("value"))
		used += v
		fmt.Fprintf(&b, `<span class="seg c-%s" style="width:%.2f%%"></span>`, safeColor(s.attr("color")), v/total*100)
	}
	rest := total - used
	if rest < 0 {
		rest = 0
	}
	fmt.Fprintf(&b, `<span class="rest" style="width:%.2f%%"></span></div></div>`, rest/total*100)
	return b.String()
}

func legend(series []series) string {
	if len(series) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="legend">`)
	for _, s := range series {
		fmt.Fprintf(&b, `<span>%s</span>`, html.EscapeString(s.name))
	}
	b.WriteString(`</div>`)
	return b.String()
}

func renderRecord(n *node) string {
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
