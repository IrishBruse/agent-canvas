# agent-canvas

Turn agent-canvas markup into one self-contained HTML file.

The CLI reads XML-shaped component tags, inlines CSS, and writes a single document with no other assets.
Component definitions live in `components/`.
Combined pages live in `examples/`.

```bash
just install
agent-canvas examples/release-health.canvas /tmp/release-health.html
```

`just install` runs `go install` and writes the binary to `GOPATH/bin`.

From the repo, without installing:

```bash
go run ./cmd/agent-canvas examples/release-health.canvas /tmp/release-health.html
```

The command entry point is `cmd/agent-canvas`.
Markup parsing lives in `internal/markup`.
HTML rendering lives in `internal/render`, split by component family.

## Preview

The VS Code extension in `editors/vscode` opens `.canvas` in a custom editor.
The preview renders again when the document changes or the file changes on disk.

```bash
just install-extension
```

Reload the window after installing.
Open a `.canvas` file to see the preview.
**Agent Canvas: Open Source** opens the markup beside that preview.

See [SKILL.md](./SKILL.md) for markup rules and component layout.
