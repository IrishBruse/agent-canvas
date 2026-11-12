# agent-canvas

agent-canvas turns a `.canvas` file into one HTML page.
Three tools share that file.
You can use each tool on its own.

- The CLI renders the file to HTML.
- The skill teaches an agent to write the file and to run the CLI.
- The extension shows a preview in VS Code and in Cursor.

Cursor Canvas keeps the page inside Cursor.
GitHub Copilot canvas keeps the page inside that product.
This repo keeps the page in a file you own.

Open the HTML in any browser.
Give the skill to any agent that reads a skill file.
Open the preview in VS Code or in Cursor.

## CLI

The CLI reads the markup and renders the HTML.
The CSS is inside that file.
That file is the whole page.
Run the CLI in any of these places:

- a terminal
- a script
- CI

A terminal is enough.

The CLI needs Go 1.22 or later.

```bash
go install github.com/IrishBruse/agent-canvas/cmd/agent-canvas@latest
```

`go install` puts `agent-canvas` in `GOPATH/bin`.
Put `GOPATH/bin` on `PATH`.

In a clone of this repo, `just install` runs `go install ./cmd/agent-canvas`.

Render a page with this command:

```bash
agent-canvas examples/release-health.canvas
```

The command renders the page to `examples/release-health.html`.
The command prints that path.
When the HTML path must differ, pass it as the second argument.

```bash
agent-canvas examples/release-health.canvas /tmp/release-health.html
```

Before you install the CLI, run it with `go run` in a clone:

```bash
go run ./cmd/agent-canvas examples/release-health.canvas
```

Open the HTML file in a browser.

## Skill

The installed skill is a short stub.
It tells the agent to load the workflow from the CLI.
The CLI embeds that workflow, so the text matches this binary.

```bash
agent-canvas skills get core
agent-canvas skills list
```

`skills get core` prints the markup rules.
`skills get core --ref <name>` prints one component or one example.
`skills get core --full` prints every component and every example.

The stub is [skills/agent-canvas/SKILL.md](./skills/agent-canvas/SKILL.md).
The workflow is [skill-data/core/SKILL.md](./skill-data/core/SKILL.md).
`agent-canvas install-skill` writes the stub to `~/.agents/skills/agent-canvas/SKILL.md`.

```bash
agent-canvas install-skill
```

In a clone, `just install-skill` runs `go run ./cmd/agent-canvas install-skill`.

## Extension

The Agent Canvas extension shows a preview of `.canvas` in VS Code and in Cursor.
The preview updates when the file changes.
The extension calls the same CLI.
The extension build needs Node.js and npm.

In a clone, `just install-extension` does the build and the install.
The commands below match that recipe.

```bash
go build -o editors/vscode/bin/agent-canvas ./cmd/agent-canvas
cd editors/vscode
npm install
npm test
npx --yes @vscode/vsce package --allow-missing-repository -o agent-canvas.vsix
cursor --install-extension agent-canvas.vsix --force
```

Run the first command from the repo root.
Run the commands after `cd` in `editors/vscode`.
When `cursor` is absent, use `code` for the last command.

Reload the window after the install.
Open a `.canvas` file.
The preview is the default editor for that file.
Run **Agent Canvas: Open Source** to show the markup beside the preview.
Run **Agent Canvas: Open Preview** on a `.canvas` file in the explorer.

An empty `agentCanvas.executable` uses the binary inside the extension.
If that binary is missing, the extension runs `agent-canvas` on `PATH`.
To use another binary, set `agentCanvas.executable` to its path.

## Markup

Write one root element.
Put a parameter on an attribute.
Put nested components inside the element.
`markdown` holds markdown.
Use it for this text:

- prose
- lists
- tables
- inline code

```xml
<Stack gap="16">
  <h1>Checkout reliability</h1>
  <Stat value="99.2%" label="Success" tone="success"/>
  <markdown>
    Payments **p95** crossed the budget.
  </markdown>
</Stack>
```

These files are full pages:

- `examples/release-health.canvas` — stats, a line chart, a table, and a usage bar
- `examples/pull-request.canvas` — a diff, pills, and cards
- `examples/context-usage.canvas` — a usage bar and nested sections

Each file under `components/` shows one component.

Layout (`components/layout/`):

- `Stack`, `Row`, `Grid`, `Card`, `hr`, `Spacer`, `CollapsibleSection`

Text (`components/text/`):

- `h1`, `h2`, `h3`, `p`, `markdown`, `code`, `link`

Status (`components/status/`):

- `Stat`, `Callout`, `Pill`, `table`, `TodoList`, `TodoListCard`, `Swatch`

Charts (`components/chart/`):

- `BarChart`, `LineChart`, `PieChart`, `UsageBar`

Diff (`components/diff/`):

- `DiffView`, `DiffStats`

Controls (`components/controls/`):

- `button`, `IconButton`, `input`, `textarea`, `Checkbox`, `Toggle`, `select`

Read [skill-data/core/SKILL.md](./skill-data/core/SKILL.md) for the agent markup rules.

## Parts of this repo

- `cmd/agent-canvas` — the CLI
- `skills/agent-canvas` — the discovery stub
- `skill-data` — the skill text the CLI serves
- `editors/vscode` — the VS Code and Cursor extension
- `components/` — one file per component, grouped by folder
- `examples/` — full pages
- `internal/markup` — the parser
- `internal/render` — the HTML writer
