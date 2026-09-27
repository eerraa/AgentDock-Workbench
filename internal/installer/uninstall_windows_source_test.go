//go:build windows

package installer

import (
	"context"
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"testing"
)

func TestWindowsOptionalUninstallCommandsUseNoConsolePolicy(t *testing.T) {
	source, err := os.ReadFile("uninstall.go")
	if err != nil {
		t.Fatalf("read uninstall.go: %v", err)
	}
	text := string(source)
	anchor := "func runOptionalCmd(ctx context.Context, name string, args ...string) error"
	start := strings.Index(text, anchor)
	if start < 0 {
		t.Fatalf("uninstall.go missing %q", anchor)
	}
	end := start + 500
	if end > len(text) {
		end = len(text)
	}
	if !strings.Contains(text[start:end], "installerCommand(ctx, name, args...)") {
		t.Fatal("runOptionalCmd must apply the Windows no-console background process policy")
	}
	command := installerCommand(context.Background(), "fixture-never-started.exe", "inspect")
	if command.SysProcAttr == nil || !command.SysProcAttr.HideWindow || command.SysProcAttr.CreationFlags&windows.CREATE_NO_WINDOW == 0 {
		t.Fatal("shared installer command did not suppress console creation")
	}
	if command.Process != nil {
		t.Fatal("constructing a native command must not start it")
	}
}
