package installer

import (
	"context"
	"os/exec"

	processctl "github.com/uvwt/agentdock/internal/process"
)

// Every native control command invoked by Setup inherits the no-console
// policy, including status probes and rollback commands.
func installerCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	processctl.ConfigureBackground(cmd)
	return cmd
}
