package buildinfo

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestDownstreamIdentityAndCompleteSourceSHA(t *testing.T) {
	old := Commit
	t.Cleanup(func() { Commit = old })
	Commit = "0123456789abcdef0123456789abcdef01234567"
	info := Current()
	if info.SourceCommit != Commit || info.Commit != Commit[:12] {
		t.Fatalf("source identity lost: %+v", info)
	}
	if info.Distribution != "eerraa" || info.UpstreamVersion != "1.1.6" || info.DownstreamRevision != DownstreamRevision {
		t.Fatalf("downstream identity: %+v", info)
	}
	if _, err := json.Marshal(info); err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(Version, ".")
	if len(parts) != 3 {
		t.Fatalf("non-numeric product version: %s", Version)
	}
	for _, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 || n > 65535 {
			t.Fatalf("invalid Windows version component: %s", part)
		}
	}
	// The release owner assigns the main product version explicitly; the upstream
	// baseline and downstream release ordinal remain independent metadata.
	if Version != "1.1.16101" || UpstreamVersion != "1.1.6" || DownstreamRevision != 101 {
		t.Fatalf("main release identity mismatch: version=%s upstream=%s revision=%d", Version, UpstreamVersion, DownstreamRevision)
	}
}
