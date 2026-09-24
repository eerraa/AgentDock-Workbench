package app

import (
	"github.com/uvwt/agentdock/internal/activity"
	"testing"
)

func TestOwnedExecutionDescriptionsPreserveUnrecognizedAndUserText(t *testing.T) {
	for name, prefix := range executionDescriptions {
		got := executionDescriptionText(name, prefix+" · path 中文 한국어")
		if name == "mcp_tool_call" {
			if got != nil {
				t.Fatal("third-party invocation description localized")
			}
			continue
		}
		if got == nil || got.Code != "tool."+name || len(got.Args) != 1 || got.Args[0] != " · path 中文 한국어" {
			t.Fatalf("owned descriptor %s: %+v", name, got)
		}
		if executionDescriptionText(name, "user requested text") != nil {
			t.Fatal("user text gained descriptor")
		}
	}
	raw := activity.Event{ToolName: "permission.update", Title: "修改执行权限", Status: "succeeded", Summary: "scope=conversation scope_id=conv_1 mode=full revision=2；操作系统权限未改变。"}
	got := describeOwnedManagement(raw)
	if got.TitleText == nil || got.SummaryText == nil || len(got.SummaryText.Args) != 4 || got.Summary != raw.Summary {
		t.Fatal("permission presentation lost original")
	}
	raw.Status = "failed"
	raw.Summary = "permission failure 原文"
	got = describeOwnedManagement(raw)
	if got.SummaryText != nil || got.Summary != raw.Summary {
		t.Fatal("permission failure was presented as success")
	}
	raw.ToolName = "third:permission.update"
	got = describeOwnedManagement(raw)
	if got.TitleText != nil || got.SummaryText != nil {
		t.Fatal("third-party data gained trusted presentation")
	}
}
