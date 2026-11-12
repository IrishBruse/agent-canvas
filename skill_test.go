package canvas

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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
	if !strings.Contains(string(got), "agent-canvas skills get core") {
		t.Fatal("installed skill is not the discovery stub")
	}
}

func TestSkillsGetCore(t *testing.T) {
	var out, err bytes.Buffer
	if code := Skills(&out, &err, []string{"skills", "get", "core"}); code != 0 {
		t.Fatalf("code %d: %s", code, err.String())
	}
	if !strings.Contains(out.String(), "name: core") || !strings.Contains(out.String(), "Write one `.canvas` file.") {
		t.Fatalf("core skill missing: %s", out.String())
	}
}

func TestSkillsGetRefMatchesFile(t *testing.T) {
	want, err := os.ReadFile("components/controls/button.canvas")
	if err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	if code := Skills(&out, &stderr, []string{"skills", "get", "core", "--ref", "button"}); code != 0 {
		t.Fatalf("code %d: %s", code, stderr.String())
	}
	if strings.TrimRight(out.String(), "\n") != strings.TrimRight(string(want), "\n") {
		t.Fatalf("ref = %q", out.String())
	}
}

func TestSkillsList(t *testing.T) {
	var out, err bytes.Buffer
	if code := Skills(&out, &err, []string{"skills", "list"}); code != 0 {
		t.Fatal(err.String())
	}
	if !strings.Contains(out.String(), "core") {
		t.Fatalf("list = %s", out.String())
	}
	if strings.Contains(out.String(), "agent-canvas") {
		t.Fatalf("stub listed: %s", out.String())
	}
}
