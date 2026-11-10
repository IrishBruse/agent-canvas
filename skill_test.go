package canvas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallSkill(t *testing.T) {
	home := t.TempDir()
	dest, err := InstallSkill(home)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".agents", "skills", "agent-canvas", "SKILL.md")
	if dest != want {
		t.Fatalf("dest = %s", dest)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != skillMarkdown {
		t.Fatal("installed skill does not match the embedded file")
	}
	if len(got) == 0 {
		t.Fatal("embedded skill is empty")
	}
}
