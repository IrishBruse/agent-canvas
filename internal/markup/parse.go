package markup

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Attr struct {
	Name  string
	Value string
}

type Kid struct {
	Text string
	El   *Node
}

type Node struct {
	Name  string
	Attrs []Attr
	Kids  []Kid
}

func Parse(r io.Reader) (*Node, error) {
	dec := xml.NewDecoder(r)
	var stack []*Node
	var root *Node
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
			n := &Node{Name: t.Name.Local}
			for _, a := range t.Attr {
				n.Attrs = append(n.Attrs, Attr{Name: a.Name.Local, Value: a.Value})
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Kids = append(parent.Kids, Kid{El: n})
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
			if parent.Name != "markdown" {
				text = strings.TrimSpace(text)
				if text == "" {
					continue
				}
			}
			parent.Kids = append(parent.Kids, Kid{Text: text})
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

func (n *Node) Attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name == name {
			return a.Value
		}
	}
	return ""
}

func (n *Node) Truthy(name string) bool {
	switch strings.ToLower(n.Attr(name)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func (n *Node) Child(name string) *Node {
	for _, k := range n.Kids {
		if k.El != nil && k.El.Name == name {
			return k.El
		}
	}
	return nil
}

func (n *Node) Elements(name string) []*Node {
	var out []*Node
	for _, k := range n.Kids {
		if k.El != nil && k.El.Name == name {
			out = append(out, k.El)
		}
	}
	return out
}

func (n *Node) Text() string {
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
