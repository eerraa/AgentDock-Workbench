//go:build windows

package desktopruntime

import (
	"os"
	"strings"
	"testing"
)

func TestWindowsDesktopRuntimeBackgroundCommandsUseNoConsoleConfigure(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		file   string
		anchor string
	}{
		{
			name:   "scheduled task query",
			file:   "service_startup_windows.go",
			anchor: "exec.CommandContext(ctx, \"schtasks.exe\", \"/Query\"",
		},
		{
			name:   "scheduled task mutation",
			file:   "service_startup_windows.go",
			anchor: "func runScheduledTaskCommand(ctx context.Context, args ...string) error",
		},
		{
			name:   "core detached launch",
			file:   "service_windows.go",
			anchor: "exec.Command(coreBinary, \"service\", \"launch-core\", \"--runtime-root\", runtimeRoot)",
		},
		{
			name:   "tunnel supervisor launch",
			file:   "tunnel_windows.go",
			anchor: "exec.Command(supervisorBinary, \"tunnel\", \"launch\", \"--runtime-root\", runtime.root)",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			if testCase.name == "scheduled task query" {
				source, err := os.ReadFile(testCase.file)
				if err != nil {
					t.Fatal(err)
				}
				text := string(source)
				start := strings.Index(text, "func coreAutostartEnabled(")
				if start < 0 {
					t.Fatal("native task query missing")
				}
				end := strings.Index(text[start:], "func parseScheduledTaskXML(")
				if start < 0 || end < 0 {
					t.Fatal("native task query boundary missing")
				}
				body := text[start : start+end]
				if !strings.Contains(body, `nativeRuntimeTaskAction(ctx, manifest.InstallRoot, manifest, "status")`) ||
					!strings.Contains(body, "return state.Enabled, err") || strings.Contains(body, "exec.Command") || strings.Contains(body, "schtasks.exe") {
					t.Fatal("task query must remain native COM with explicit errors and no console process")
				}
				return
			}
			assertDesktopRuntimeCommandConfiguredNearby(t, testCase.file, testCase.anchor)
		})
	}
}

func assertDesktopRuntimeCommandConfiguredNearby(t *testing.T, fileName, anchor string) {
	t.Helper()

	source, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("读取 %s 失败: %v", fileName, err)
	}
	text := string(source)
	anchorIndex := strings.Index(text, anchor)
	if anchorIndex < 0 {
		t.Fatalf("%s 未找到启动点 %q", fileName, anchor)
	}

	const nearbyBytes = 700
	end := anchorIndex + nearbyBytes
	if end > len(text) {
		end = len(text)
	}
	if !strings.Contains(text[anchorIndex:end], "processcontrol.Configure(command)") {
		t.Fatalf("%s 的启动点 %q 未通过 processcontrol.Configure 收口", fileName, anchor)
	}
}

func TestWindowsNativeHostUsesSharedHealthBudgetAndSpawnDiagnostics(t *testing.T) {
	state, err := os.ReadFile("runtime_host_state_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(string(state), "func waitCapturedCoreHealth(")
	if start < 0 || !strings.Contains(string(state)[start:], "context.WithTimeout(ctx, WindowsCoreStartTimeout)") {
		t.Fatal("captured native Core must use the same startup budget")
	}
	host, err := os.ReadFile("runtime_host_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	body := string(host)
	spawn := strings.Index(body, "child, err := job.Start(coreBinary")
	diagnostic := strings.Index(body, `slog.Int("core_pid", int(child.PID))`)
	if spawn < 0 || diagnostic < spawn || !strings.Contains(body, `"task_core_host", "core_spawn"`) || strings.Contains(body, "exec.Command") {
		t.Fatal("spawn diagnostics must follow the actual creation-time Job owner, not introduce another launcher")
	}
}
