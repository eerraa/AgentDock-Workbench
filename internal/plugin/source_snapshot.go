package plugin

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Source: uvwt/agentdock 120527c7a385e24fcbe98ef1786c824b086e7aae.
// Reuse the complete bounded snapshot helper set without changing Workbench's
// package layout or metadata/digest rules. The later metadata-ignore policy
// belongs to uvwt's different package contract and is not enabled here.
func snapshotPluginTree(source, destination string, maxBytes int64, maxFiles int) error {
	source = filepath.Clean(source)
	root, err := os.OpenRoot(source)
	if err != nil {
		return err
	}
	defer root.Close()
	if err := os.MkdirAll(destination, 0o700); err != nil {
		return err
	}
	var total int64
	entries := 0
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == source {
			return nil
		}
		relative, err := filepath.Rel(source, path)
		if err != nil || !filepath.IsLocal(relative) || relative == "." {
			return fmt.Errorf("invalid Plugin source path %q", path)
		}
		entries++
		if entries > maxFiles {
			return fmt.Errorf("Plugin source exceeds %d entries", maxFiles)
		}
		info, err := validatePluginSnapshotPath(root, relative, entry.IsDir())
		if err != nil {
			return fmt.Errorf("unsafe Plugin source path %q: %w", path, err)
		}
		target := filepath.Join(destination, relative)
		if info.IsDir() {
			return os.MkdirAll(target, 0o700)
		}
		remaining := maxBytes - total
		if remaining < 0 {
			return fmt.Errorf("Plugin source exceeds %d bytes", maxBytes)
		}
		mode := info.Mode().Perm() & 0o755
		if mode == 0 {
			mode = 0o600
		}
		copied, err := snapshotPluginRegularFile(root, relative, target, mode, remaining)
		if err != nil {
			return err
		}
		total += copied
		if total > maxBytes {
			return fmt.Errorf("Plugin source exceeds %d bytes", maxBytes)
		}
		return nil
	})
}

func validatePluginSnapshotPath(root *os.Root, relative string, wantDirectory bool) (os.FileInfo, error) {
	clean := filepath.Clean(relative)
	if !filepath.IsLocal(clean) || clean == "." {
		return nil, errors.New("path is not local to the Plugin source")
	}
	parts := strings.Split(clean, string(filepath.Separator))
	current := ""
	var info os.FileInfo
	for index, part := range parts {
		if current == "" {
			current = part
		} else {
			current = filepath.Join(current, part)
		}
		var err error
		info, err = root.Lstat(current)
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("path contains symlink component %q", filepath.ToSlash(current))
		}
		if index < len(parts)-1 && !info.IsDir() {
			return nil, fmt.Errorf("path component %q is not a directory", filepath.ToSlash(current))
		}
	}
	if wantDirectory {
		if info == nil || !info.IsDir() {
			return nil, errors.New("path is not a directory")
		}
	} else if info == nil || !info.Mode().IsRegular() {
		return nil, errors.New("path is not a regular file")
	}
	return info, nil
}

func snapshotPluginRegularFile(root *os.Root, relative, destination string, mode os.FileMode, remaining int64) (int64, error) {
	before, err := validatePluginSnapshotPath(root, relative, false)
	if err != nil {
		return 0, err
	}
	in, err := root.Open(relative)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	opened, err := in.Stat()
	if err != nil {
		return 0, err
	}
	if !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		return 0, errors.New("Plugin source file changed while opening snapshot")
	}
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return 0, err
	}
	copied, copyErr := io.Copy(out, io.LimitReader(in, remaining+1))
	closeErr := out.Close()
	if copyErr != nil {
		return copied, copyErr
	}
	if closeErr != nil {
		return copied, closeErr
	}
	if copied > remaining {
		return copied, fmt.Errorf("Plugin source exceeds byte limit")
	}
	afterOpen, err := in.Stat()
	if err != nil {
		return copied, err
	}
	afterPath, err := validatePluginSnapshotPath(root, relative, false)
	if err != nil {
		return copied, err
	}
	if !os.SameFile(opened, afterOpen) || !os.SameFile(opened, afterPath) ||
		opened.Size() != afterOpen.Size() || opened.ModTime() != afterOpen.ModTime() || opened.Mode() != afterOpen.Mode() {
		return copied, errors.New("Plugin source file changed while creating snapshot")
	}
	return copied, nil
}
