package insertion

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInsertionTextIdentityAndDeadline(t *testing.T) {
	store, now, target := fixture(t)
	text := "  원문 보존\n" + strings.Repeat("추가 지시 🙂 ", 80) + "\n  "
	item, err := store.Add(context.Background(), target, "summary_once", text)
	if err != nil {
		t.Fatal(err)
	}
	if item.Text != text {
		t.Fatal("summary replaced original")
	}
	same, err := store.Add(context.Background(), target, "summary_once", text)
	if err != nil || same.ID != item.ID || !same.ExpiresAt.Equal(item.ExpiresAt) {
		t.Fatalf("summary changed idempotency: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(store.root, "queue.json"))
	if err != nil {
		t.Fatal(err)
	}
	var disk diskState
	if err := json.Unmarshal(data, &disk); err != nil {
		t.Fatal(err)
	}
	if len(disk.Items) != 1 || disk.Items[0].Text != text {
		t.Fatal("original duplicated or changed in storage")
	}
	reloaded, err := New(store.root, "summary_restart", func() time.Time { return *now })
	if err != nil {
		t.Fatal(err)
	}
	list, err := reloaded.List(context.Background(), target.Owner, target.Conversation)
	if err != nil || len(list) != 1 || list[0].Text != text || !list[0].ExpiresAt.Equal(item.ExpiresAt) {
		t.Fatalf("reload changed message: %v", err)
	}
	reserved, err := reloaded.Reserve(context.Background(), target, "call_summary", *now)
	if err != nil || len(reserved) != 1 || reserved[0].Text != text {
		t.Fatal("preview was delivered instead of full text")
	}
	finished, err := reloaded.Finish(context.Background(), "call_summary", true)
	if err != nil || len(finished) != 1 || finished[0].Text != text {
		t.Fatal("attachment lost original")
	}
}
