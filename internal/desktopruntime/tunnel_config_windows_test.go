//go:build windows

package desktopruntime

import (
	"context"
	"testing"
)

func TestElevatedConfigureLeavesSupervisorToScheduledTask(t *testing.T) {
	runtime := tunnelRuntime{
		manifest: Manifest{
			PrivilegeMode:     "elevated",
			AgentDockTaskName: "AgentDock",
		},
		mode: "named",
	}
	if err := startConfiguredTunnel(context.Background(), runtime); err != nil {
		t.Fatal(err)
	}
}
