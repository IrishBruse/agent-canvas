---
name: agent-canvas
description: Turn agent-canvas markup into one HTML file. Use when rendering a .canvas file, running agent-canvas, or writing canvas markup and plain text.
---

# Agent canvas

`agent-canvas` reads one markup file and writes one HTML file.
The markup is XML-shaped component tags.
One file per component lives under `components/`.
The folders are `layout`, `text`, `status`, `chart`, `diff`, and `controls`.
Combined pages live in `examples/`.
The page inlines its CSS and has no other assets.

## Render

1. Run the CLI on the markup file.
   The first argument is the `.canvas` file.
   The second argument is the HTML path.
   When the second argument is omitted, the HTML file uses the markup path with `.html` in place of `.canvas`.

```bash
go run ./cmd/agent-canvas examples/release-health.canvas /tmp/release-health.html
```

Run that from the repo root.
Done when the command prints the HTML path and that file is one document.

With the Agent Canvas extension installed, VS Code and Cursor open `.canvas` in a live preview.
The preview renders again when the file changes.

## Markup

Write one root element.
Put a parameter on an attribute.
Put nested components and slot content inside the element.
Use an HTML tag when HTML has that element.
A heading is `h1`, `h2`, or `h3`.
A link is `link`.
Set `href` to the URL.
A table uses `table`, `thead`, `tbody`, `tr`, `th`, and `td`.
Put `align` on a `th` to align that column.
Put `tone` on a `tr` to mark that row.
Open the file with that name under `components/` when the shape of a component is the question.
Open a file in `examples/` when the question is how components sit together.
A component with no parameters is an empty tag.

## Plain text

Write prose, lists, tables, and inline code as markdown inside `markdown`.
Indent the markdown with the surrounding markup.
The CLI strips the shared indent before rendering.

```xml
<markdown>
  Payments **p95** crossed the budget.

  - checkout
  - catalog
</markdown>
```

Use `p`, `h1`, or `code` when that element is the whole component.
