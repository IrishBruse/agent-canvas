---
name: core
description: Write a .canvas file and render it to HTML. Use when writing canvas markup or rendering a .canvas file.
---

# Agent canvas

Write one `.canvas` file.
Render it with the `agent-canvas` binary.

```bash
agent-canvas path/to/file.canvas
```

The command prints the HTML path.
That file is one document.
Pass a second argument when the HTML path must differ.

When the binary is not on `PATH`, run it from a clone of this repo:

```bash
go run ./cmd/agent-canvas path/to/file.canvas
```

Run that command from the repo root.

## Markup

Write one root element.
A parameter is an attribute.
Nested components and slot content go inside the element.
Use the HTML tag when HTML has that element.
`h1`, `h2`, and `h3` are headings.
`link` is a link, and `href` is the URL.
`markdown` holds prose, lists, tables, and inline code.
Indent that markdown with the surrounding markup.
The CLI strips the shared indent before rendering.

```xml
<markdown>
  Payments **p95** crossed the budget.

  - checkout
  - catalog
</markdown>
```

Use `p`, `h1`, or `code` when that element is the whole component.

## Load a component

Load one component when its shape is the question:

```bash
agent-canvas skills get core --ref button
agent-canvas skills get core --ref Stack
```

The value is the file stem or the path under `components/`.

- Layout: `Stack`, `Row`, `Grid`, `Card`, `hr`, `Spacer`, `CollapsibleSection`
- Text: `h1`, `h2`, `h3`, `p`, `markdown`, `code`, `link`
- Status: `Stat`, `Callout`, `Pill`, `table`, `TodoList`, `TodoListCard`, `Swatch`
- Charts: `BarChart`, `LineChart`, `PieChart`, `UsageBar`
- Diff: `DiffView`, `DiffStats`
- Controls: `button`, `IconButton`, `input`, `textarea`, `Checkbox`, `Toggle`, `select`

Load a full page the same way: `release-health`, `pull-request`, `context-usage`.
`agent-canvas skills get core --full` prints every component and every example at once.
