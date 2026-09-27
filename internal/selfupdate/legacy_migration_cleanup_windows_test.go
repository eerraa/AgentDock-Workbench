//go:build windows

package selfupdate

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

func migrationCleanupFixture(t *testing.T) (string, windowsLegacyMigrationPlan) {
	t.Helper()
	root, err := os.MkdirTemp("", "agentdock-legacy-migration-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	runtimeRoot := t.TempDir()
	return root, windowsLegacyMigrationPlan{
		ParentPID: 1234, RuntimeRoot: runtimeRoot,
		CorePath:   filepath.Join(runtimeRoot, "bin", "agentdock.exe"),
		TrayPath:   filepath.Join(runtimeRoot, "bin", "agentdock-tray.exe"),
		PayloadDir: filepath.Join(root, "payload"), Version: "v0.9.1", CleanupRoot: root,
	}
}

func holdMigrationCleanupMutex(t *testing.T) func() {
	t.Helper()
	runtime.LockOSThread()
	name, err := windows.UTF16PtrFromString(windowsLegacyMigrationMutexName)
	if err != nil {
		runtime.UnlockOSThread()
		t.Fatal(err)
	}
	handle, err := windows.CreateMutex(nil, true, name)
	if err != nil {
		if handle != 0 {
			_ = windows.CloseHandle(handle)
		}
		runtime.UnlockOSThread()
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			t.Skip("another migration owns the native mutex")
		}
		t.Fatal(err)
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			if err := windows.ReleaseMutex(handle); err != nil {
				t.Error(err)
			}
			_ = windows.CloseHandle(handle)
			runtime.UnlockOSThread()
		})
	}
}

func TestLegacyMigrationCleanupPreservesPriorRecovery(t *testing.T) {
	for _, kind := range []string{"malformed_plan", "invalid_plan", "valid_plan"} {
		t.Run(kind, func(t *testing.T) {
			root, plan := migrationCleanupFixture(t)
			backup := filepath.Join(root, "stable-backup", "original.exe")
			if err := os.MkdirAll(filepath.Dir(backup), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(backup, []byte("original-recovery-material"), 0600); err != nil {
				t.Fatal(err)
			}
			var err error
			if kind == "malformed_plan" {
				planPath := filepath.Join(root, "migration-plan.json")
				if err := os.WriteFile(planPath, []byte("{"), 0600); err != nil {
					t.Fatal(err)
				}
				err = runWindowsLegacyMigrationHelper(t.Context(), filepath.Join(root, "agentdock-legacy-migration-helper.exe"), planPath)
			} else {
				if kind == "invalid_plan" {
					plan.ParentPID = 0
				}
				err = finalizeWindowsLegacyMigration(t.Context(), plan)
			}
			if err == nil {
				t.Error("a prior unverified recovery directory must not be reused")
			}
			data, readErr := os.ReadFile(backup)
			if readErr != nil || string(data) != "original-recovery-material" {
				t.Fatalf("prior recovery was deleted or overwritten: %v", readErr)
			}
		})
	}
}

func TestLegacyMigrationMalformedPlanCannotCleanActiveOwner(t *testing.T) {
	release := holdMigrationCleanupMutex(t)
	defer release()
	root, _ := migrationCleanupFixture(t)
	planPath := filepath.Join(root, "migration-plan.json")
	if err := os.WriteFile(planPath, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		done <- runWindowsLegacyMigrationHelper(context.Background(), filepath.Join(root, "agentdock-legacy-migration-helper.exe"), planPath)
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Error("malformed plan accepted")
		}
	case <-time.After(2 * time.Second):
		release()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("helper remained blocked")
		}
		t.Fatal("malformed-plan cleanup must not wait on or interrupt a live owner")
	}
	if data, err := os.ReadFile(planPath); err != nil || string(data) != "{" {
		t.Fatalf("a failed second helper removed the active owner's files: %v", err)
	}
}

func TestLegacyMigrationMutexWaitHonorsCancellation(t *testing.T) {
	release := holdMigrationCleanupMutex(t)
	defer release()
	root, plan := migrationCleanupFixture(t)
	marker := filepath.Join(root, "retained")
	if err := os.WriteFile(marker, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- finalizeWindowsLegacyMigration(ctx, plan) }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	var result error
	timely := true
	select {
	case result = <-done:
	case <-time.After(500 * time.Millisecond):
		timely = false
		release()
		select {
		case result = <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("duplicate helper never left the mutex")
		}
	}
	if !timely || !errors.Is(result, context.Canceled) {
		t.Errorf("mutex wait ignored cancellation: timely=%v err=%v", timely, result)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "unchanged" {
		t.Errorf("unowned cleanup directory changed: %v", err)
	}
}
