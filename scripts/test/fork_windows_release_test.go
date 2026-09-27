package scripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWindowsForkDistributionTargets(t *testing.T) {
	const repository = "eerraa/AgentDock-Workbench"
	for _, item := range []struct {
		path    string
		markers []string
	}{
		{"internal/selfupdate/update.go", []string{"https://api.github.com/repos/" + repository + "/releases/latest"}},
		{"scripts/install/install.ps1", []string{"https://github.com/" + repository + "/releases/latest/download", "https://github.com/" + repository + "/releases/download/$normalizedVersion"}},
		{"packaging/windows/AgentDock.iss", []string{"AppPublisherURL=https://github.com/" + repository, "AppSupportURL=https://github.com/" + repository + "/issues", "AppUpdatesURL=https://github.com/" + repository + "/releases"}},
	} {
		t.Run(item.path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(item.path)))
			if err != nil {
				t.Fatal(err)
			}
			for _, marker := range item.markers {
				if !strings.Contains(string(data), marker) {
					t.Errorf("fork distribution destination missing: %s", marker)
				}
			}
			if strings.Contains(string(data), "github.com/repos/A-m-o-r-F-a-t-i/agentdock/releases") || strings.Contains(string(data), "github.com/A-m-o-r-F-a-t-i/agentdock/releases") {
				t.Error("Windows fork can still resolve an upstream update payload")
			}
		})
	}
}

func TestWindowsTaskRollbackRetainsRuntimeOwner(t *testing.T) {
	data, err := os.ReadFile("../install/install.ps1")
	if err != nil {
		t.Fatal(err)
	}
	source := strings.ReplaceAll(string(data), "\r\n", "\n")
	start := strings.Index(source, "$restoreTaskActionResult = Start-ElevatedAgentDockTaskAction")
	if start < 0 {
		t.Fatal("task rollback call missing")
	}
	end := strings.Index(source[start:], "if (-not $restoreTaskActionResult.Started)")
	if end < 0 {
		t.Fatal("task rollback result must be checked")
	}
	call := source[start : start+end]
	for _, want := range []string{"-Action restore", "-RuntimeRoot $runtimeDir", "-TaskUser $taskUser", "-BackupDirectory $taskBackupDirectory"} {
		if !strings.Contains(call, want) {
			t.Errorf("rollback lost original owner binding: %s", want)
		}
	}
}

func TestFormalWindowsBuildRequiresCleanSourceForEveryVersion(t *testing.T) {
	data, err := os.ReadFile("../../packaging/windows/build-windows-release.ps1")
	if err != nil {
		t.Fatal(err)
	}
	source := string(data)
	if !strings.Contains(source, "if (-not $Candidate -and $sourceChanges.Count -gt 0)") {
		t.Error("formal clean-source gate is missing or limited to a historical version")
	}
	for _, required := range []string{"source_state_changed", "repository='eerraa/AgentDock-Workbench'"} {
		if !strings.Contains(source, required) {
			t.Errorf("formal build provenance guard missing: %s", required)
		}
	}
}
