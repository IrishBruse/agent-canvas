package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: agent-canvas <input.canvas> [output.html]")
		os.Exit(2)
	}
	in := os.Args[1]
	out := strings.TrimSuffix(in, filepath.Ext(in)) + ".html"
	if len(os.Args) == 3 {
		out = os.Args[2]
	}
	raw, err := os.ReadFile(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	root, err := parse(strings.NewReader(string(raw)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, []byte(page(root)), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}
