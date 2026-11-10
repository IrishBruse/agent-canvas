package canvas

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed SKILL.md
var skillMarkdown string

// InstallSkill writes the embedded skill to ~/.agents/skills/agent-canvas/SKILL.md under home.
func InstallSkill(home string) (string, error) {
	dest := filepath.Join(home, ".agents", "skills", "agent-canvas", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(dest, []byte(skillMarkdown), 0o644); err != nil {
		return "", err
	}
	return dest, nil
}
