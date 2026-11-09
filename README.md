# agent-canvas

Turn agent-canvas markup into one self-contained HTML file.

The CLI reads XML-shaped component tags, inlines CSS, and writes a single document with no other assets.
Component definitions live in `components/`; combined pages are in `examples/`.

```bash
go run . examples/release-health.canvas /tmp/release-health.html
```

See [SKILL.md](./SKILL.md) for markup rules and component layout.
