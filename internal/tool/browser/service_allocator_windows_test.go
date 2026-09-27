package browser

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// Exercise the same upstream allocator options with a disposable native child,
// not an installed browser or any user's profile.
func TestWindowsLocalExecAllocatorOptionsUseRequestTimeout(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	request := StartRequest{Headless: true, Viewport: Viewport{Width: 800, Height: 600}, Timeout: 200 * time.Millisecond}
	options := localExecAllocatorOptions(executable, filepath.Join(t.TempDir(), "profile"), request)
	var child *exec.Cmd
	options = append(options, chromedp.ModifyCmdFunc(func(cmd *exec.Cmd) {
		cmd.Args = []string{executable, "-test.run=^TestWindowsLocalAllocatorHelper$"}
		cmd.Env = append(cmd.Environ(), "AGENTDOCK_BROWSER_ALLOCATOR_HELPER=1")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
		child = cmd
	}))
	allocator, stopAllocator := chromedp.NewExecAllocator(context.Background(), options...)
	browser, stopBrowser := chromedp.NewContext(allocator)
	defer func() {
		stopBrowser()
		stopAllocator()
		if child == nil || child.ProcessState == nil {
			t.Error("allocator did not reap its native helper")
		}
	}()
	done := make(chan error, 1)
	started := time.Now()
	go func() { done <- chromedp.Run(browser) }()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "websocket url timeout reached") {
			t.Fatalf("allocator error=%v, want request-derived websocket URL timeout", err)
		}
		if elapsed := time.Since(started); elapsed < request.Timeout || elapsed > 2*time.Second {
			t.Fatalf("timeout elapsed=%s, expected %s within native launch allowance", elapsed, request.Timeout)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("local allocator ignored the request timeout")
	}
}

func TestWindowsLocalAllocatorHelper(t *testing.T) {
	if os.Getenv("AGENTDOCK_BROWSER_ALLOCATOR_HELPER") != "1" {
		return
	}
	_, _ = os.Stdout.WriteString("allocator fixture started\n")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
