using System.Text.Json;

namespace AgentDock.ControlPanel;

public static class InsertionPresentation
{
    public static bool Unconfirmed(string status) => status is "inner_appended" or "outer_forwarded" or "delivery_unknown";
    public static string State(JsonElement item) => item.Text("status") switch
    {
        "pending" => UiText.Get("InsertionStatePending"),
        "reserved" => UiText.Get("InsertionStateReserved"),
        "inner_appended" => UiText.Get("InsertionStateInner"),
        "outer_forwarded" => UiText.Get("InsertionStateForwarded"),
        "acknowledged" => item.Text("acknowledged_by") == "host_context_committed" ? UiText.Get("InsertionStateContextConfirmed") : UiText.Get("InsertionStateReceiverConfirmed"),
        "attached" => UiText.Get("InsertionStateLegacyAttached"),
        "delivery_unknown" => UiText.Get("InsertionStateUnknown"),
        "target_changed" => UiText.Get("InsertionStateTargetChanged"),
        "expired" => UiText.Get("InsertionStateExpired"),
        "cancelled" => UiText.Get("InsertionStateCancelled"),
        _ => UiText.Get("InsertionStateNotRecorded")
    };

    public static string Reason(JsonElement item) => item.Text("delivery_reason") switch
    {
        "legacy_inner_response_without_receipt" => UiText.Get("InsertionReasonLegacy"),
        "receipt_missing_deadline_elapsed" => UiText.Get("InsertionReasonExpired"),
        "host_receipt_not_negotiated" => UiText.Get("InsertionReasonNoHostReceipt"),
        "awaiting_host_receipt" => UiText.Get("InsertionReasonAwaitHost"),
        "awaiting_context_commit" => UiText.Get("InsertionReasonAwaitContext"),
        "process_restarted_before_receipt" => UiText.Get("InsertionReasonRestarted"),
        "inner_response_not_committed" => UiText.Get("InsertionReasonInnerUncommitted"),
        "outer_projection_failed" => UiText.Get("InsertionReasonProjectionFailed"),
        "context_commit_failed" => UiText.Get("InsertionReasonContextFailed"),
        _ => ""
    };

    // Old clients and future fields cannot expose delivery secrets or create a
    // fake executable tool row. The UI keeps only the presentation whitelist.
    public static JsonElement Snapshot(JsonElement item)
    {
        var result = new Dictionary<string, object?> { ["record_kind"] = "insertion" };
        foreach (var key in new[] { "insertion_id", "conversation_id", "task_id", "thread_id", "workspace_id", "text", "status", "created_at", "updated_at", "expires_at", "call_id", "inner_appended_at", "outer_forwarded_at", "acknowledged_at", "acknowledged_by", "outer_call_id", "host_type", "delivery_reason" })
            if (item.Field(key).ValueKind == JsonValueKind.String) result[key == "call_id" ? "inner_call_id" : key] = item.Text(key);
        result["delivery_attempts"] = item.Number("delivery_attempts");
        result["sequence"] = item.Number("sequence");
        result["retry_requested"] = item.Flag("retry_requested");
        return JsonSerializer.SerializeToElement(result);
    }
}

public sealed partial class ExecutionCallRow
{
    private DateTimeOffset? _insertionNow;
    public bool IsInsertion => _value.Text("record_kind") == "insertion";
    public string InsertionText => IsInsertion ? _value.Text("text") : "";
    public string InsertionPreview => InsertionText.Replace('\r', ' ').Replace('\n', ' ');
    public DateTimeOffset TimelineAt => _value.Date("request_received_at") ?? _value.Date("created_at") ?? DateTimeOffset.MinValue;
    public bool CanRedeliverInsertion => IsInsertion && InsertionPresentation.Unconfirmed(Status) && !_value.Flag("retry_requested") && _value.Number("delivery_attempts") is > 0 and < 6 && _insertionNow is { } now && _value.Date("expires_at") > now;
    public bool CanCancelInsertion => IsInsertion && (Status is "pending" or "target_changed" || InsertionPresentation.Unconfirmed(Status));
    public string InsertionHint => State + "\n" + InsertionPresentation.Reason(_value) + (_value.Number("delivery_attempts") > 0 ? UiText.Format("InsertionAttemptCount", _value.Number("delivery_attempts")) : "") + (_value.Flag("retry_requested") ? UiText.Get("InsertionRetryRequested") : "");
    public string InsertionDetails => InsertionText + "\n\n" + InsertionHint + UiText.Get("InsertionMessagePrefix") + Id + UiText.Get("InsertionSentPrefix") + TimelineAt.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss") +
        (_value.Text("inner_call_id").Length > 0 ? UiText.Get("InsertionInnerCallPrefix") + _value.Text("inner_call_id") : "") +
        (_value.Text("outer_call_id").Length > 0 ? UiText.Get("InsertionOuterCallPrefix") + _value.Text("outer_call_id") : "") +
        (_value.Date("acknowledged_at") is { } at ? UiText.Get("InsertionConfirmedPrefix") + at.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss") : "");

    public static ExecutionCallRow FromInsertion(JsonElement item, DateTimeOffset? now = null)
    {
        var row = new ExecutionCallRow(InsertionPresentation.Snapshot(item));
        row._insertionNow = now;
        return row;
    }

    public void ApplyInsertion(JsonElement item, DateTimeOffset? now)
    {
        if (!IsInsertion || Id != item.Text("insertion_id") || ConversationId != item.Text("conversation_id")) throw new InvalidOperationException("Insertion row identity changed.");
        if (_value.Date("updated_at") > item.Date("updated_at")) return;
        _value = InsertionPresentation.Snapshot(item); _insertionNow = now; Notify();
    }
}
