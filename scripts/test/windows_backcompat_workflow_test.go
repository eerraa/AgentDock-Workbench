package scripts

import (
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

// Compatibility validation must consume the verified payload and stay behind
// the explicit installation opt-in. A static-only job is not migration proof.
func TestWindowsLegacyCompatibilityGatesUseVerifiedPayload(t *testing.T) {
	var workflow struct {
		Jobs map[string]struct {
			Steps []struct {
				Name string         `yaml:"name"`
				Uses string         `yaml:"uses"`
				If   string         `yaml:"if"`
				Run  string         `yaml:"run"`
				With map[string]any `yaml:"with"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal([]byte(readWorkflow(t, "windows-installer.yml")), &workflow); err != nil {
		t.Fatal(err)
	}
	job, ok := workflow.Jobs["validate"]
	if !ok {
		t.Fatal("missing Windows validation job")
	}
	const optIn = "github.event_name == 'workflow_dispatch' && inputs.installation_tests == true"
	prepared, migration, historical := -1, -1, -1
	fullHistory := false
	for index, step := range job.Steps {
		if strings.HasPrefix(step.Uses, "actions/checkout@") && step.With["fetch-depth"] == 0 {
			fullHistory = true
		}
		if step.Name == "Prepare offline AMD64 payload" {
			prepared = index
			if step.If != optIn {
				t.Fatal("payload construction lost explicit installation opt-in")
			}
			for _, required := range []string{
				`Copy-Item .\packaging\windows\compat\manage-windows.ps1 .\dist\manage-windows.ps1 -Force`,
				`.\dist\manage-windows.ps1, .\dist\share`,
			} {
				if !strings.Contains(step.Run, required) {
					t.Errorf("verified payload missing compatibility member: %s", required)
				}
			}
		}
		if strings.Contains(step.Run, "test-windows-release-backcompat.ps1 -ArchivePath $archivePath") {
			historical = index
			if step.If != optIn {
				t.Fatal("published-updater fixture must remain explicitly opted in")
			}
			archive := strings.Index(step.Run, "Compress-Archive")
			check := strings.Index(step.Run, "test-windows-release-backcompat.ps1 -ArchivePath $archivePath")
			if archive < 0 || check <= archive {
				t.Fatal("updater check must inspect the constructed current payload")
			}
		}
		if strings.Contains(step.Run, "test-windows-legacy-online-migration.ps1") {
			migration = index
			if step.If != optIn {
				t.Fatal("native migration lost explicit installation opt-in")
			}
			if strings.Count(step.Run, "test-windows-legacy-online-migration.ps1") != 2 ||
				strings.Count(step.Run, "-ArchivePath $env:AMD64_ARCHIVE") != 2 ||
				strings.Count(step.Run, "-ChecksumPath $env:AMD64_CHECKSUM") != 2 ||
				!strings.Contains(step.Run, "-InjectTrayReplaceFailure") {
				t.Fatal("migration must test verified payload in normal and injected-failure recovery paths")
			}
		}
	}
	if !fullHistory {
		t.Error("historical updater sources require a full-history checkout")
	}
	if prepared < 0 || historical != prepared {
		t.Error("published-updater validation is missing from payload preparation")
	}
	if migration <= prepared {
		t.Error("native migration must follow verified payload preparation")
	}
}
