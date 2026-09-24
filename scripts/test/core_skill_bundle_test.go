package scripts

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	skills "github.com/uvwt/agentdock/internal/skill"
	skillbundle "github.com/uvwt/agentdock/internal/skill/bundle"
	skillstate "github.com/uvwt/agentdock/internal/skill/state"
)

func coreSkillBundleBuilder(t *testing.T) (string, string) {
	t.Helper()
	python := ""
	for _, candidate := range []string{"python3", "python"} {
		path, err := exec.LookPath(candidate)
		if err != nil {
			continue
		}
		probe := exec.Command(path, "--version")
		probe.Dir = t.TempDir()
		output, err := probe.CombinedOutput()
		if err == nil && strings.Contains(string(output), "Python 3") {
			python = path
			break
		}
	}
	if python == "" {
		if os.Getenv("CI") != "" {
			t.Fatal("Python 3 is required to test the core Skill bundle builder in CI")
		}
		t.Skip("Python is required to test the core Skill bundle builder")
	}
	script, err := filepath.Abs("../../packaging/build-core-skill-bundle.py")
	if err != nil {
		t.Fatal(err)
	}
	return python, script
}

func TestCoreSkillBundleNormalizesTextLineEndings(t *testing.T) {
	python, script := coreSkillBundleBuilder(t)
	build := func(lineEnding string) string {
		t.Helper()
		repoRoot := t.TempDir()
		for _, name := range []string{"agentdock-user-guide", "skill-authoring", "skill-installation"} {
			skillRoot := filepath.Join(repoRoot, "core-skills", name)
			if err := os.MkdirAll(skillRoot, 0o700); err != nil {
				t.Fatal(err)
			}
			document := "---\nname: " + name + "\ndescription: Test Skill.\nversion: 1.0.0\n---\n\n# Test\n"
			document = strings.ReplaceAll(document, "\n", lineEnding)
			if err := os.WriteFile(filepath.Join(skillRoot, "SKILL.md"), []byte(document), 0o600); err != nil {
				t.Fatal(err)
			}
			scriptBody := strings.ReplaceAll("print('test')\n", "\n", lineEnding)
			if err := os.WriteFile(filepath.Join(skillRoot, "run.py"), []byte(scriptBody), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		output := filepath.Join(t.TempDir(), "bundle")
		command := exec.Command(python, script, "--repo-root", repoRoot, "--output", output)
		command.Dir = t.TempDir()
		if data, err := command.CombinedOutput(); err != nil {
			t.Fatalf("build core Skill bundle: %v\n%s", err, data)
		}
		return output
	}

	lfBundle := build("\n")
	crlfBundle := build("\r\n")
	for _, relative := range []string{
		"manifest.json",
		filepath.Join("packages", "agentdock-user-guide.zip"),
		filepath.Join("packages", "skill-authoring.zip"),
		filepath.Join("packages", "skill-installation.zip"),
	} {
		lfData, err := os.ReadFile(filepath.Join(lfBundle, relative))
		if err != nil {
			t.Fatal(err)
		}
		crlfData, err := os.ReadFile(filepath.Join(crlfBundle, relative))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(lfData, crlfData) {
			t.Fatalf("core Skill bundle differs between LF and CRLF input: %s", relative)
		}
	}
}

func TestCoreSkillBundleBootstraps(t *testing.T) {
	python, script := coreSkillBundleBuilder(t)
	output := filepath.Join(t.TempDir(), "bundle")
	command := exec.Command(python, script, "--output", output)
	command.Dir = t.TempDir()
	if data, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build official core Skill bundle: %v\n%s", err, data)
	}
	state, err := skillstate.New(filepath.Join(t.TempDir(), "skills"))
	if err != nil {
		t.Fatal(err)
	}
	manager, err := skills.New(state)
	if err != nil {
		t.Fatal(err)
	}
	result, err := skillbundle.Bootstrap(t.Context(), state, manager, output)
	if err != nil {
		t.Fatalf("bootstrap builder output: %v", err)
	}
	if len(result.Skills) != 3 {
		t.Fatalf("installed %d core Skills, want 3", len(result.Skills))
	}
	for _, entry := range result.Skills {
		active, err := state.ActiveVersion(entry.Name)
		if err != nil || active != entry.Version {
			t.Fatalf("%s active=%q, want %q: %v", entry.Name, active, entry.Version, err)
		}
	}
}

func TestCoreSkillBundleRejectsMissingVersion(t *testing.T) {
	python, script := coreSkillBundleBuilder(t)
	for _, test := range []struct{ name, frontmatter string }{
		{"missing", ""},
		{"empty", "version: \n"},
		{"quoted-empty", "version: \"\"\n"},
		{"nested-only", "metadata:\n  version: 1.0.0\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			repoRoot := t.TempDir()
			for _, name := range []string{"agentdock-user-guide", "skill-authoring", "skill-installation"} {
				root := filepath.Join(repoRoot, "core-skills", name)
				if err := os.MkdirAll(root, 0o700); err != nil {
					t.Fatal(err)
				}
				version := "version: 1.0.0\n"
				if name == "agentdock-user-guide" {
					version = test.frontmatter
				}
				document := "---\nname: " + name + "\ndescription: Test Skill.\n" + version + "---\n\n# Test\n"
				if err := os.WriteFile(filepath.Join(root, "SKILL.md"), []byte(document), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			output := filepath.Join(t.TempDir(), "bundle")
			command := exec.Command(python, script, "--repo-root", repoRoot, "--output", output)
			command.Dir = t.TempDir()
			data, err := command.CombinedOutput()
			if err == nil || !strings.Contains(string(data), "version is required") {
				t.Fatalf("unversioned core Skill was not rejected: %v\n%s", err, data)
			}
			if _, err := os.Stat(filepath.Join(output, "manifest.json")); !os.IsNotExist(err) {
				t.Fatalf("invalid bundle published a manifest: %v", err)
			}
		})
	}
}
