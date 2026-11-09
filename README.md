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

The skill teaches an agent to write `.canvas` files and to run the CLI.
The text is in [SKILL.md](./SKILL.md).
The agent writes a file in your project.
You run the CLI on that file.
Any agent that reads the skill can do this work.

Copy the skill into the skill folder your agent reads.
This repo uses `~/.agents/skills`:

```bash
mkdir -p ~/.agents/skills/agent-canvas
cp SKILL.md ~/.agents/skills/agent-canvas/SKILL.md
```

Run those commands from a clone.
`just install-skill` copies `SKILL.md` to that folder.
When the agent loads skills from this repo, point it at `SKILL.md`.

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
`MD` holds markdown.
Use it for this text:

- prose
- lists
- tables
- inline code

```xml
<Stack gap="16">
  <H1>Checkout reliability</H1>
  <Stat value="99.2%" label="Success" tone="success"/>
  <MD>
    Payments **p95** crossed the budget.
  </MD>
</Stack>
```

These files are full pages:

- `examples/release-health.canvas` — stats, a line chart, a table, and a usage bar
- `examples/pull-request.canvas` — a diff, pills, and cards
- `examples/context-usage.canvas` — a usage bar and nested sections

Each file in `components/` shows one component.

Layout:

- `Stack`, `Row`, `Grid`, `Card`, `Divider`, `Spacer`, `CollapsibleSection`

Text:

- `H1`, `H2`, `H3`, `Text`, `MD`, `Code`, `Link`

Status and data:

- `Stat`, `Callout`, `Pill`, `Table`, `TodoList`, `TodoListCard`, `Swatch`
- `BarChart`, `LineChart`, `PieChart`, `UsageBar`
- `DiffView`, `DiffStats`

Controls on the page:

- `Button`, `IconButton`, `TextInput`, `TextArea`, `Checkbox`, `Toggle`, `Select`

Read [SKILL.md](./SKILL.md) for the agent markup rules.

## Parts of this repo

- `cmd/agent-canvas` — the CLI
- `SKILL.md` — the skill
- `editors/vscode` — the VS Code and Cursor extension
- `components/` — one file per component
- `examples/` — full pages
- `internal/markup` — the parser
- `internal/render` — the HTML writer
