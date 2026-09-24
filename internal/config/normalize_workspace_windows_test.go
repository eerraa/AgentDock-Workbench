//go:build windows

package config

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

func TestNormalizeSecuresHomeAndLeavesWorkspaceDACL(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(workspace, "project.txt")
	if err := os.WriteFile(child, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeWorkspace, err := daclSDDL(workspace)
	if err != nil {
		t.Fatal(err)
	}
	beforeChild, err := daclSDDL(child)
	if err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		AgentDockHome:       home,
		AgentDockDefaultDir: workspace,
		Host:                "127.0.0.1",
		Port:                8765,
		LogLevel:            "info",
	}
	if err := cfg.Normalize(); err != nil {
		t.Fatal(err)
	}
	if cfg.AgentDockHome != home || cfg.AgentDockDefaultDir != workspace {
		t.Fatalf("paths = %s %s", cfg.AgentDockHome, cfg.AgentDockDefaultDir)
	}

	afterWorkspace, err := daclSDDL(workspace)
	if err != nil {
		t.Fatal(err)
	}
	afterChild, err := daclSDDL(child)
	if err != nil {
		t.Fatal(err)
	}
	if afterWorkspace != beforeWorkspace || afterChild != beforeChild {
		t.Fatalf("workspace DACL changed\nbefore dir %s\nafter dir %s\nbefore child %s\nafter child %s", beforeWorkspace, afterWorkspace, beforeChild, afterChild)
	}

	homeSDDL, err := daclSDDL(home)
	if err != nil {
		t.Fatal(err)
	}
	if homeSDDL == beforeWorkspace {
		t.Fatal("home DACL was not secured")
	}
}

func daclSDDL(path string) (string, error) {
	descriptor, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return "", err
	}
	return descriptor.String(), nil
}
