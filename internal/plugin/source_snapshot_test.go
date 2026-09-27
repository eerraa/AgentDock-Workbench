package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPluginSourceSnapshotPreservesBytesAndCallerOwnership(t *testing.T) {
	source, target := t.TempDir(), filepath.Join(t.TempDir(), "snapshot")
	if err := os.Mkdir(filepath.Join(source, "nested"), 0700); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"nested/문서.txt": "original 한국어 中文\n", ".agentdock-import.json": "package metadata must not silently disappear"}
	for name, value := range files {
		if err := os.WriteFile(filepath.Join(source, filepath.FromSlash(name)), []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := copyPackageTree(source, target); err != nil {
		t.Fatal(err)
	}
	for name, want := range files {
		path := filepath.FromSlash(name)
		if err := os.WriteFile(filepath.Join(source, path), []byte("subsequent user change"), 0600); err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(target, path))
		if err != nil || string(got) != want {
			t.Fatalf("snapshot %s changed: %q %v", name, got, err)
		}
	}
}

func TestPluginSourceSnapshotBoundsBytesAndAllEntries(t *testing.T) {
	for _, tc := range []struct {
		name      string
		bytes     int64
		entries   int
		wantError bool
	}{
		{"exact", 4, 2, false}, {"bytes", 3, 2, true}, {"directory_entries", 4, 1, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source, target := t.TempDir(), filepath.Join(t.TempDir(), "copy")
			if err := os.Mkdir(filepath.Join(source, "sub"), 0700); err != nil {
				t.Fatal(err)
			}
			original := filepath.Join(source, "sub", "file.txt")
			if err := os.WriteFile(original, []byte("data"), 0600); err != nil {
				t.Fatal(err)
			}
			err := snapshotPluginTree(source, target, tc.bytes, tc.entries)
			if (err != nil) != tc.wantError {
				t.Fatalf("error=%v, wantError=%v", err, tc.wantError)
			}
			if data, readErr := os.ReadFile(original); readErr != nil || string(data) != "data" {
				t.Fatal("source was modified", readErr)
			}
			if info, statErr := os.Stat(filepath.Join(target, "sub", "file.txt")); statErr == nil && info.Size() > tc.bytes+1 {
				t.Fatal("copy exceeded bounded overflow detection")
			}
		})
	}
}

func TestPluginSourceSnapshotDoesNotOverwriteDestination(t *testing.T) {
	source, target := t.TempDir(), t.TempDir()
	for path, value := range map[string]string{filepath.Join(source, "file.txt"): "source", filepath.Join(target, "file.txt"): "retained"} {
		if err := os.WriteFile(path, []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := copyPackageTree(source, target); err == nil {
		t.Fatal("overwrote existing destination")
	}
	if data, err := os.ReadFile(filepath.Join(target, "file.txt")); err != nil || string(data) != "retained" {
		t.Fatal("destination was altered", err)
	}
}

func TestPluginSourceSnapshotRejectsLinkComponents(t *testing.T) {
	for _, directory := range []bool{false, true} {
		name := "file"
		if directory {
			name = "directory"
		}
		t.Run(name, func(t *testing.T) {
			source, target, outside := t.TempDir(), filepath.Join(t.TempDir(), "copy"), t.TempDir()
			private := filepath.Join(outside, "private.txt")
			if err := os.WriteFile(private, []byte("outside sentinel"), 0600); err != nil {
				t.Fatal(err)
			}
			linkTarget := private
			if directory {
				linkTarget = outside
			}
			link := filepath.Join(source, "linked")
			if err := os.Symlink(linkTarget, link); err != nil {
				t.Fatalf("cannot establish required local link fixture: %v", err)
			}
			if err := copyPackageTree(source, target); err == nil {
				t.Fatal("accepted linked package content")
			}
			if _, err := os.Lstat(filepath.Join(target, "linked")); !os.IsNotExist(err) {
				t.Fatal("link content entered snapshot", err)
			}
			if data, err := os.ReadFile(private); err != nil || string(data) != "outside sentinel" {
				t.Fatal("outside source was modified", err)
			}
		})
	}
}

func TestPluginSourceSnapshotRejectsEscapingRelativeFiles(t *testing.T) {
	source := t.TempDir()
	root, err := os.OpenRoot(source)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	for _, name := range []string{"../outside", ".", filepath.Join(t.TempDir(), "absolute.txt")} {
		target := filepath.Join(t.TempDir(), "copy")
		if _, err := snapshotPluginRegularFile(root, name, target, 0600, 8); err == nil {
			t.Fatalf("accepted non-local path %q", name)
		}
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			t.Fatal("invalid relative input created output", err)
		}
	}
}

func TestPluginSourceSnapshotRejectsFileBeyondBudget(t *testing.T) {
	source := t.TempDir()
	file := filepath.Join(source, "data")
	if err := os.WriteFile(file, []byte(strings.Repeat("x", 64)), 0600); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(source)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	target := filepath.Join(t.TempDir(), "bounded")
	copied, err := snapshotPluginRegularFile(root, "data", target, 0600, 4)
	if err == nil || copied != 5 {
		t.Fatalf("copy=%d error=%v, want bounded rejection after five bytes", copied, err)
	}
	if info, err := os.Stat(target); err != nil || info.Size() != 5 {
		t.Fatal("unbounded file copy", err)
	}
	if data, err := os.ReadFile(file); err != nil || len(data) != 64 {
		t.Fatal("source changed", err)
	}
}
