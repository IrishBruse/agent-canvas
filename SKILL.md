---
name: agent-canvas
description: Turn agent-canvas markup into one HTML file. Use when rendering a .canvas file, running agent-canvas, or writing canvas markup and plain text.
---

# Agent canvas

`agent-canvas` reads one markup file and writes one HTML file.
The markup is XML-shaped component tags.
One file per component lives in `components/` beside this skill.
Combined pages live in `examples/`.
The page inlines its CSS and has no other assets.

## Render

1. Run the CLI on the markup file.
   The first argument is the `.canvas` file.
   The second argument is the HTML path.
   When the second argument is omitted, the HTML file uses the markup path with `.html` in place of `.canvas`.

```bash
go run . examples/release-health.canvas /tmp/release-health.html
```

Run that from the repo root.
Done when the command prints the HTML path and that file is one document.

## Markup

Write one root element.
Put a parameter on an attribute.
Put nested components and slot content inside the element.
Open the matching file in `components/` when the shape of a component is the question.
Open a file in `examples/` when the question is how components sit together.
A component with no parameters is an empty tag.

## Plain text

Write prose, lists, tables, and inline code as markdown inside `MD`.
Indent the markdown with the surrounding markup.
The CLI strips the shared indent before rendering.

```xml
<MD>
  Payments **p95** crossed the budget.

  - checkout
  - catalog
</MD>
```

Use `Text`, a heading, or `Code` when that single component is the whole element.
