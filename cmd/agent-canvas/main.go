package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/IrishBruse/agent-canvas"
	"github.com/IrishBruse/agent-canvas/internal/markup"
	"github.com/IrishBruse/agent-canvas/internal/render"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "install-skill":
		if len(os.Args) != 2 {
			usage()
			os.Exit(2)
		}
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		dest, err := canvas.InstallSkill(home)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(dest)
		return
	case "skills":
		os.Exit(canvas.Skills(os.Stdout, os.Stderr, os.Args[1:]))
	}
	if len(os.Args) > 3 {
		usage()
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
	root, err := markup.Parse(strings.NewReader(string(raw)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, []byte(render.Page(root)), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: agent-canvas <input.canvas> [output.html]")
	fmt.Fprintln(os.Stderr, "       agent-canvas skills [list]")
	fmt.Fprintln(os.Stderr, "       agent-canvas skills get <name> [--ref <file>]... [--full]")
	fmt.Fprintln(os.Stderr, "       agent-canvas skills get --all")
	fmt.Fprintln(os.Stderr, "       agent-canvas install-skill")
}
