package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/IrishBruse/agent-canvas/internal/markup"
)

type series struct {
	name  string
	tone  string
	value []float64
}

func readSeries(n *markup.Node) []series {
	host := n.Child("series")
	if host == nil {
		host = n
	}
	var out []series
	for _, s := range append(host.Elements("series"), host.Elements("ChartSeries")...) {
		item := series{name: s.Attr("name"), tone: s.Attr("tone")}
		for _, v := range s.Elements("value") {
			item.value = append(item.value, number(v.Text()))
		}
		if data := s.Child("data"); data != nil {
			for _, v := range data.Elements("value") {
				item.value = append(item.value, number(v.Text()))
			}
		}
		out = append(out, item)
	}
	return out
}

func renderBars(n *markup.Node) string {
	cats := texts(n.Child("categories"), "category")
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

func renderLine(n *markup.Node) string {
	series := readSeries(n)
	cats := texts(n.Child("categories"), "category")
	maxV := 1.0
	minV := 0.0
	if !n.Truthy("beginAtZero") && n.Attr("beginAtZero") == "false" {
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

func renderPie(n *markup.Node) string {
	var points []*markup.Node
	if data := n.Child("data"); data != nil {
		points = append(data.Elements("point"), data.Elements("ChartDataPoint")...)
	}
	if n.Name == "ChartDataPoint" || n.Attr("label") != "" {
		points = []*markup.Node{n}
	}
	total := 0.0
	for _, p := range points {
		total += number(first(p.Attr("value"), p.Text()))
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
		v := number(first(p.Attr("value"), p.Text()))
		start := acc / total * 100
		acc += v
		end := acc / total * 100
		color := colors[p.Attr("tone")]
		if color == "" {
			color = fallback[i%len(fallback)]
		}
		stops = append(stops, fmt.Sprintf("%s %.1f%% %.1f%%", color, start, end))
		fmt.Fprintf(&legend, `<span>%s %s</span>`, html.EscapeString(p.Attr("label")), html.EscapeString(first(p.Attr("value"), p.Text())))
	}
	legend.WriteString(`</div>`)
	className := "pie"
	if n.Truthy("donut") {
		className += " donut"
	}
	size := first(n.Attr("size"), "200")
	style := fmt.Sprintf("width:%spx;height:%spx;background:conic-gradient(%s)", size, size, strings.Join(stops, ","))
	return el("figure", "diagram", "", el("div", className, style, "")+legend.String())
}

func renderUsage(n *markup.Node) string {
	host := n.Child("segments")
	if host == nil {
		host = n
	}
	segs := append(host.Elements("segment"), host.Elements("UsageBarSegment")...)
	total := number(n.Attr("total"))
	if total <= 0 {
		total = 1
	}
	var b strings.Builder
	b.WriteString(`<div class="usage"><div class="usage-labels">`)
	fmt.Fprintf(&b, `<span>%s</span><span>%s</span></div><div class="usage-track">`, html.EscapeString(n.Attr("topLeftLabel")), html.EscapeString(n.Attr("topRightLabel")))
	used := 0.0
	for _, s := range segs {
		v := number(s.Attr("value"))
		used += v
		fmt.Fprintf(&b, `<span class="seg c-%s" style="width:%.2f%%"></span>`, safeColor(s.Attr("color")), v/total*100)
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
