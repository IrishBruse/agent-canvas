package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type attr struct {
	Name  string
	Value string
}

type kid struct {
	Text string
	El   *node
}

type node struct {
	Name  string
	Attrs []attr
	Kids  []kid
}

func parse(r io.Reader) (*node, error) {
	dec := xml.NewDecoder(r)
	var stack []*node
	var root *node
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			n := &node{Name: t.Name.Local}
			for _, a := range t.Attr {
				n.Attrs = append(n.Attrs, attr{Name: a.Name.Local, Value: a.Value})
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Kids = append(parent.Kids, kid{El: n})
			}
			stack = append(stack, n)
			if root == nil {
				root = n
			}
		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("unexpected end tag %s", t.Name.Local)
			}
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			parent := stack[len(stack)-1]
			text := string(t)
			if parent.Name != "MD" {
				text = strings.TrimSpace(text)
				if text == "" {
					continue
				}
			}
			parent.Kids = append(parent.Kids, kid{Text: text})
		}
	}
	if root == nil {
		return nil, fmt.Errorf("empty markup")
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("unclosed %s", stack[len(stack)-1].Name)
	}
	return root, nil
}

func (n *node) attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name == name {
			return a.Value
		}
	}
	return ""
}

func (n *node) truthy(name string) bool {
	switch strings.ToLower(n.attr(name)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func (n *node) child(name string) *node {
	for _, k := range n.Kids {
		if k.El != nil && k.El.Name == name {
			return k.El
		}
	}
	return nil
}

func (n *node) elements(name string) []*node {
	var out []*node
	for _, k := range n.Kids {
		if k.El != nil && k.El.Name == name {
			out = append(out, k.El)
		}
	}
	return out
}

func (n *node) text() string {
	var b strings.Builder
	for _, k := range n.Kids {
		if k.El == nil && k.Text != "" {
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(k.Text)
		}
	}
	return b.String()
}
