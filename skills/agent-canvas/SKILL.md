---
name: agent-canvas
description: Turn agent-canvas markup into one HTML file. Use when rendering a .canvas file, running agent-canvas, or writing canvas markup and plain text.
user-invocable: false
---

# agent-canvas

`agent-canvas` reads one markup file and writes one HTML file.

## Start here

This file is a discovery stub, not the usage guide.
Load the workflow from the CLI before you write markup or run `agent-canvas`:

```bash
agent-canvas skills get core
```

When the binary is not on `PATH`, run `go run ./cmd/agent-canvas skills get core` from the repo root.

That guide ends with the component names to load next.
Use `agent-canvas skills get core --full` only when you need every component at once.

The CLI serves skill content that matches the installed version.
This stub only points at `skills get core`.

Run `agent-canvas skills list` to see the skills on the installed version.

## Why agent-canvas

- One markup file becomes one HTML file
- The CSS is inside that file
- The CLI, this skill, and the editor preview share that file
