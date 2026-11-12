package canvas

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
)

//go:embed skill-data
var skillData embed.FS

//go:embed all:components
var componentFiles embed.FS

//go:embed all:examples
var exampleFiles embed.FS

type skillFile struct {
	path    string
	content string
}

type skillInfo struct {
	name        string
	description string
	hidden      bool
	body        string
	files       []skillFile
}

// Skills serves the skill text embedded in this binary.
// args starts at the "skills" command.
func Skills(stdout, stderr io.Writer, args []string) int {
	skills, err := loadSkills()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if len(args) == 1 || args[1] == "list" {
		if len(args) > 2 {
			fmt.Fprintln(stderr, "usage: agent-canvas skills [list]")
			return 2
		}
		writeSkillList(stdout, skills)
		return 0
	}
	if args[1] != "get" {
		fmt.Fprintf(stderr, "Unknown skills subcommand: %s\n", args[1])
		return 1
	}
	names, refs, full, all, code := parseSkillGet(stderr, args[2:])
	if code != 0 {
		return code
	}
	return writeSkillGet(stdout, stderr, skills, names, refs, full, all)
}

func loadSkills() ([]skillInfo, error) {
	entries, err := fs.ReadDir(skillData, "skill-data")
	if err != nil {
		return nil, err
	}
	var skills []skillInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := "skill-data/" + entry.Name()
		body, err := fs.ReadFile(skillData, dir+"/SKILL.md")
		if err != nil {
			continue
		}
		name, description, hidden := parseFrontmatter(string(body))
		if name == "" {
			name = entry.Name()
		}
		info := skillInfo{
			name:        name,
			description: description,
			hidden:      hidden,
			body:        string(body),
		}
		info.files = skillRefs(info.name, dir)
		skills = append(skills, info)
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].name < skills[j].name })
	return skills, nil
}

func skillRefs(name, dir string) []skillFile {
	var files []skillFile
	for _, sub := range []string{"references", "templates"} {
		files = append(files, readSkillFiles(skillData, dir+"/"+sub)...)
	}
	if name == "core" {
		files = append(files, readSkillFiles(componentFiles, "components")...)
		files = append(files, readSkillFiles(exampleFiles, "examples")...)
	}
	return files
}

func readSkillFiles(fsys fs.FS, root string) []skillFile {
	var files []skillFile
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		files = append(files, skillFile{path: path, content: string(b)})
		return nil
	})
	if err != nil {
		return nil
	}
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	return files
}

func parseFrontmatter(content string) (name, description string, hidden bool) {
	content = strings.TrimPrefix(content, "\uFEFF")
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return "", "", false
	}
	rest := content[3:]
	rest = strings.TrimPrefix(rest, "\r\n")
	rest = strings.TrimPrefix(rest, "\n")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", "", false
	}
	lines := strings.Split(rest[:end], "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		switch strings.TrimSpace(key) {
		case "name":
			name = strings.Trim(val, `"'`)
		case "description":
			if val == "" || val == "|" || val == ">" || val == "|-" || val == ">-" {
				var parts []string
				for i+1 < len(lines) {
					next := strings.TrimRight(lines[i+1], "\r")
					if next != "" && !strings.HasPrefix(next, " ") && !strings.HasPrefix(next, "\t") {
						break
					}
					i++
					parts = append(parts, strings.TrimSpace(next))
				}
				description = strings.Join(parts, " ")
			} else {
				description = strings.Trim(val, `"'`)
			}
		case "hidden":
			hidden = val == "true" || val == "yes"
		}
	}
	return name, description, hidden
}

func writeSkillList(w io.Writer, skills []skillInfo) {
	var visible []skillInfo
	for _, s := range skills {
		if !s.hidden {
			visible = append(visible, s)
		}
	}
	if len(visible) == 0 {
		fmt.Fprintln(w, "No skills found")
		return
	}
	width := 0
	for _, s := range visible {
		if len(s.name) > width {
			width = len(s.name)
		}
	}
	for _, s := range visible {
		fmt.Fprintf(w, " %-*s  %s\n", width, s.name, truncateDescription(s.description, 70))
	}
}

func truncateDescription(desc string, max int) string {
	if len(desc) <= max {
		return desc
	}
	cut := desc[:max]
	if i := strings.LastIndex(cut, " "); i > 0 {
		cut = cut[:i]
	}
	return cut + "..."
}

func parseSkillGet(stderr io.Writer, args []string) (names, refs []string, full, all bool, code int) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--full":
			full = true
		case arg == "--all":
			all = true
		case arg == "--ref":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "--ref needs a value. Usage: agent-canvas skills get <name> --ref <file>")
				return nil, nil, false, false, 2
			}
			i++
			refs = append(refs, args[i])
		case strings.HasPrefix(arg, "--ref="):
			refs = append(refs, strings.TrimPrefix(arg, "--ref="))
		case strings.HasPrefix(arg, "-"):
			fmt.Fprintf(stderr, "Unknown flag: %s\n", arg)
			return nil, nil, false, false, 2
		default:
			names = append(names, arg)
		}
	}
	return names, refs, full, all, 0
}

func writeSkillGet(stdout, stderr io.Writer, skills []skillInfo, names, refs []string, full, all bool) int {
	targets, code := selectSkills(stderr, skills, names, all)
	if code != 0 {
		return code
	}
	if len(refs) > 0 {
		if all || full || len(targets) != 1 {
			fmt.Fprintln(stderr, "--ref needs exactly one skill name and cannot be combined with --all or --full")
			return 2
		}
		return writeSkillRefs(stdout, stderr, targets[0], refs)
	}
	for i, s := range targets {
		if i > 0 {
			fmt.Fprint(stdout, "\n---\n\n")
		}
		writeText(stdout, s.body)
		if full {
			for _, file := range s.files {
				fmt.Fprintf(stdout, "\n--- %s ---\n\n", file.path)
				writeText(stdout, file.content)
			}
		}
	}
	return 0
}

func selectSkills(stderr io.Writer, skills []skillInfo, names []string, all bool) ([]skillInfo, int) {
	if all {
		var targets []skillInfo
		for _, s := range skills {
			if !s.hidden {
				targets = append(targets, s)
			}
		}
		if len(targets) == 0 {
			fmt.Fprintln(stderr, "No skills found")
			return nil, 1
		}
		return targets, 0
	}
	if len(names) == 0 {
		fmt.Fprintln(stderr, "No skill name provided. Usage: agent-canvas skills get <name>")
		return nil, 2
	}
	var targets []skillInfo
	for _, name := range names {
		found := false
		for _, s := range skills {
			if s.name == name {
				targets = append(targets, s)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(stderr, "Skill not found: %s. Available: %s\n", name, availableSkills(skills))
			return nil, 1
		}
	}
	return targets, 0
}

func availableSkills(skills []skillInfo) string {
	var names []string
	for _, s := range skills {
		if !s.hidden {
			names = append(names, s.name)
		}
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

func writeSkillRefs(stdout, stderr io.Writer, skill skillInfo, refs []string) int {
	var selected []skillFile
	seen := map[string]bool{}
	for _, wanted := range refs {
		file, err := resolveRef(skill.files, wanted)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if seen[file.path] {
			continue
		}
		seen[file.path] = true
		selected = append(selected, file)
	}
	for i, file := range selected {
		if len(selected) > 1 {
			if i > 0 {
				fmt.Fprintln(stdout)
			}
			fmt.Fprintf(stdout, "--- %s ---\n\n", file.path)
		}
		writeText(stdout, file.content)
	}
	return 0
}

func resolveRef(files []skillFile, wanted string) (skillFile, error) {
	wanted = strings.TrimPrefix(wanted, "./")
	var matches []skillFile
	for _, file := range files {
		base := file.path
		if i := strings.LastIndex(base, "/"); i >= 0 {
			base = base[i+1:]
		}
		stem := base
		if i := strings.LastIndex(stem, "."); i >= 0 {
			stem = stem[:i]
		}
		if file.path == wanted || base == wanted || stem == wanted {
			matches = append(matches, file)
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return skillFile{}, fmt.Errorf("Reference not found: %s", wanted)
	default:
		paths := make([]string, len(matches))
		for i, file := range matches {
			paths[i] = file.path
		}
		return skillFile{}, fmt.Errorf("Ambiguous reference %s: matches %s", wanted, strings.Join(paths, ", "))
	}
}

func writeText(w io.Writer, content string) {
	fmt.Fprint(w, content)
	if content != "" && !strings.HasSuffix(content, "\n") {
		fmt.Fprintln(w)
	}
}
