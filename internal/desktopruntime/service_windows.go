//go:build windows

package desktopruntime

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

const windowsCoreStartTimeout = 2 * time.Minute

func platformServiceStatus(ctx context.Context, runtimeRoot string) (ServiceStatus, error) {
	manifest, _, err := loadDesktopManifest(runtimeRoot)
	if err != nil {
		return ServiceStatus{}, err
	}
	coreBinary := ActiveCoreBinary(runtimeRoot, manifest)
	supervisorPID, err := activeTunnelSupervisorPID(runtimeRoot, coreBinary)
	if err != nil {
		return ServiceStatus{}, err
	}
	excluded := map[uint32]struct{}{}
	if supervisorPID != 0 {
		excluded[supervisorPID] = struct{}{}
	}
	coreProcesses, err := processIDsAtPathExcept(coreBinary, excluded)
	if err != nil {
		return ServiceStatus{}, err
	}
	running := len(coreProcesses) > 0
	healthy := testHealth(ctx, manifest.HealthURL())
	startupEnabled, err := coreAutostartEnabled(ctx, manifest)
	if err != nil {
		return ServiceStatus{}, fmt.Errorf("读取 AgentDock 开机启动状态失败: %w", err)
	}
	return ServiceStatus{Running: running || healthy, Healthy: healthy, StartupEnabled: startupEnabled}, nil
}

func platformServiceAction(ctx context.Context, runtimeRoot, action string) error {
	manifest, root, err := loadDesktopManifest(runtimeRoot)
	if err != nil {
		return err
	}

	switch action {
	case "start":
		return startCore(ctx, manifest, root)
	case "stop":
		return stopCore(ctx, manifest, root)
	case "restart":
		if manifest.UsesScheduledTask() {
			return restartScheduledCore(ctx, manifest, root)
		}
		if err := stopCore(ctx, manifest, root); err != nil {
			return err
		}
		return startCore(ctx, manifest, root)
	default:
		return fmt.Errorf("不支持的 Windows 服务操作：%s", action)
	}
}

func loadDesktopManifest(runtimeRoot string) (Manifest, string, error) {
	root, err := filepath.Abs(strings.TrimSpace(runtimeRoot))
	if err != nil {
		return Manifest{}, "", fmt.Errorf("解析 Windows 运行目录失败: %w", err)
	}
	manifest, err := Load(filepath.Join(root, "runtime.json"))
	if err != nil {
		return Manifest{}, "", fmt.Errorf("读取 Windows 运行清单失败: %w", err)
	}
	return manifest, root, nil
}

func startCore(ctx context.Context, manifest Manifest, runtimeRoot string) error {
	if err := CheckExecutionCompatibility(ctx, runtimeRoot, ActiveCoreBinary(runtimeRoot, manifest)); err != nil {
		return err
	}
	if manifest.UsesScheduledTask() {
		return startScheduledCore(ctx, manifest, runtimeRoot)
	}
	if testHealth(ctx, manifest.HealthURL()) {
		return nil
	}
	recoverAbandonedCoreLocks(manifest, runtimeRoot)
	if err := startDetachedCore(manifest, runtimeRoot); err != nil {
		return err
	}
	return waitForHealth(ctx, manifest.HealthURL(), windowsCoreStartTimeout)
}

// The scheduled task is the only owner of an elevated core. Stop ends that
// task and does not terminate an arbitrary matching image path. Restart fails
// while the previous core still answers health, then starts the same task.
func startScheduledCore(ctx context.Context, manifest Manifest, runtimeRoot string) error {
	if testHealth(ctx, manifest.HealthURL()) {
		return nil
	}
	recoverAbandonedCoreLocks(manifest, runtimeRoot)
	if err := StartInteractiveScheduledTask(ctx, runtimeRoot, manifest.AgentDockTaskName); err != nil {
		return err
	}
	return waitForHealth(ctx, manifest.HealthURL(), windowsCoreStartTimeout)
}

func restartScheduledCore(ctx context.Context, manifest Manifest, runtimeRoot string) error {
	if err := stopScheduledCore(ctx, manifest, runtimeRoot); err != nil {
		return err
	}
	if testHealth(ctx, manifest.HealthURL()) {
		return fmt.Errorf("AgentDock 核心在计划任务结束后仍响应 %s", manifest.HealthURL())
	}
	return startScheduledCore(ctx, manifest, runtimeRoot)
}

// handOffScheduledCore ends the previous elevated core and starts the same
// scheduled task. The caller does not wait for the port; that task owns the core.
func handOffScheduledCore(ctx context.Context, manifest Manifest, runtimeRoot string) error {
	if err := stopScheduledCore(ctx, manifest, runtimeRoot); err != nil {
		return err
	}
	if testHealth(ctx, manifest.HealthURL()) {
		return fmt.Errorf("AgentDock 核心在计划任务结束后仍响应 %s", manifest.HealthURL())
	}
	recoverAbandonedCoreLocks(manifest, runtimeRoot)
	return StartInteractiveScheduledTask(ctx, runtimeRoot, manifest.AgentDockTaskName)
}

func stopScheduledCore(ctx context.Context, manifest Manifest, runtimeRoot string) error {
	endErr := runScheduledTaskCommand(ctx, "/End", "/TN", scheduledTaskPath(manifest.AgentDockTaskName))
	if endErr != nil && !scheduledTaskNotRunning(endErr) {
		return endErr
	}
	coreBinary := ActiveCoreBinary(runtimeRoot, manifest)
	deadline := time.Now().Add(15 * time.Second)
	for {
		running, err := processRunningAtPath(coreBinary)
		if err != nil {
			return err
		}
		if !running && !testHealth(ctx, manifest.HealthURL()) {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("计划任务结束后 AgentDock 核心仍在运行: %s", manifest.HealthURL())
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func scheduledTaskNotRunning(err error) bool {
	message := err.Error()
	return strings.Contains(strings.ToLower(message), "no running instance") ||
		strings.Contains(message, "没有运行") ||
		strings.Contains(message, "현재 실행") ||
		strings.Contains(message, "실행 중이지")
}

func stopCore(ctx context.Context, manifest Manifest, runtimeRoot string) error {
	coreBinary := ActiveCoreBinary(runtimeRoot, manifest)
	excluded := map[uint32]struct{}{}
	ancestorPIDs, err := ancestorProcessIDsAtPath(coreBinary)
	if err != nil {
		return fmt.Errorf("识别 AgentDock Core 调用链失败: %w", err)
	}
	for processID := range ancestorPIDs {
		excluded[processID] = struct{}{}
	}
	supervisorPID, err := activeTunnelSupervisorPID(runtimeRoot, coreBinary)
	if err != nil {
		return fmt.Errorf("识别 Tunnel supervisor 失败: %w", err)
	}
	if supervisorPID != 0 {
		// Core 与 Tunnel supervisor 共用 agentdock.exe。停止 Core 时必须保留 supervisor，
		// 否则一次普通 Core 重启就会悄悄丢失 Tunnel 的后续自恢复能力。
		excluded[supervisorPID] = struct{}{}
	}

	if manifest.UsesScheduledTask() {
		return stopScheduledCore(ctx, manifest, runtimeRoot)
	}
	if err := stopBinaryProcessesExcept(ctx, coreBinary, excluded, 15*time.Second); err != nil {
		return fmt.Errorf("停止 AgentDock 核心失败: %w", err)
	}
	return nil
}

func startDetachedCore(manifest Manifest, runtimeRoot string) error {
	coreBinary := ActiveCoreBinary(runtimeRoot, manifest)
	if info, err := os.Stat(coreBinary); err != nil || info.IsDir() {
		return fmt.Errorf("找不到 AgentDock 核心程序: %s", coreBinary)
	}
	command := exec.Command(coreBinary, "service", "launch-core", "--runtime-root", runtimeRoot)
	// launch-core 会自行把运行日志写入受限轮转文件；父进程不再持有同一路径的追加句柄。
	command.Dir = manifest.AgentDockDefaultDir
	if command.Dir == "" {
		command.Dir = defaultWindowsWorkDir()
	}
	command.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP | windows.DETACHED_PROCESS,
	}
	if err := command.Start(); err != nil {
		return fmt.Errorf("启动 AgentDock 核心失败: %w", err)
	}
	if err := command.Process.Release(); err != nil {
		return fmt.Errorf("释放 AgentDock 后台进程句柄失败: %w", err)
	}
	return nil
}

func defaultWindowsWorkDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	workDir := filepath.Join(home, "AgentDock")
	if err := os.MkdirAll(workDir, 0o700); err != nil {
		return home
	}
	return workDir
}

func waitForHealth(ctx context.Context, url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if testHealth(ctx, url) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
	return fmt.Errorf("AgentDock 健康检查失败: %s", url)
}

func testHealth(ctx context.Context, url string) bool {
	requestCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode >= 200 && response.StatusCode < 300
}
