using System.Collections;
using System.Globalization;
using System.Resources;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using System.Text.RegularExpressions;
using AgentDock.ControlPanel;

internal static class LocalizationIntegrationTests
{
    internal static void Run(Action<bool, string> check)
    {
        var resources = new ResourceManager("AgentDock.ControlPanel.Resources.UiStrings", typeof(UiText).Assembly);
        Dictionary<string, string> Inventory(CultureInfo culture) =>
            resources.GetResourceSet(culture, true, false)!.Cast<DictionaryEntry>()
                .ToDictionary(item => (string)item.Key, item => (string)item.Value!);
        var neutral = Inventory(CultureInfo.InvariantCulture);
        string Slots(string text) => string.Join(",", Regex.Matches(text, @"(?<!\{)\{([0-9]+)(?:[^{}]*)\}(?!\})")
            .Select(match => match.Groups[1].Value).Distinct().Order());
        var now = DateTimeOffset.Parse("2026-09-27T00:00:00Z");
        const string original = "  원문 中文 <&>\r\n🙂 keep spaces  ";
        foreach (var (locale, pending, confirmed) in new[] {
            ("en", "Waiting for the next tool call", "Confirmed in the model context"),
            ("zh-CN", "等待下一次工具调用", "模型上下文已确认接收"),
            ("ko-KR", "다음 도구 호출 대기", "모델 컨텍스트 수신 확인됨") })
        {
            UiText.ApplyPreference(locale);
            var entries = locale == "en" ? neutral : Inventory(CultureInfo.GetCultureInfo(locale));
            check(neutral.Keys.ToHashSet().SetEquals(entries.Keys), "resource key parity: " + locale);
            foreach (var (key, value) in neutral)
            {
                check(!string.IsNullOrEmpty(entries[key]), "nonempty resource " + locale + "/" + key);
                check(Slots(value) == Slots(entries[key]), "format argument parity " + locale + "/" + key);
            }
            var item = JsonSerializer.SerializeToElement(new { insertion_id = "locale-1", conversation_id = "conv-locale",
                text = original, status = "pending", created_at = now, updated_at = now, expires_at = now.AddMinutes(5),
                receipt_token = "secret-must-not-render", call_id = "call-original", outer_call_id = "call-outer" });
            var message = ExecutionCallRow.FromInsertion(item, now);
            check(message.State == pending && message.Title == UiText.Get("InsertionUserSupplement"), "localized pending label: " + locale);
            check(message.InsertionText == original && message.InsertionDetails.StartsWith(original, StringComparison.Ordinal), "exact supplement original: " + locale);
            check(!message.Technical.Contains("secret-must-not-render") && message.Tool == "" && !message.NeedsApproval, "supplement provenance and whitelist: " + locale);
            var receipt = JsonSerializer.SerializeToElement(new { insertion_id = "locale-1", conversation_id = "conv-locale", text = original,
                status = "acknowledged", acknowledged_by = "host_context_committed", updated_at = now.AddSeconds(1), created_at = now });
            message.ApplyInsertion(receipt, now.AddSeconds(1));
            check(message.State == confirmed && !message.CanRedeliverInsertion && message.InsertionText == original, "receipt evidence and original: " + locale);
            var user = new ExecutionCallRow(JsonSerializer.SerializeToElement(new { tool_name = "third:read", activity_label = original, activity_label_source = "user" }));
            check(user.OriginalLabel == original && user.Title.Contains("원문 中文") && user.Tool == "third:read", "user label provenance: " + locale);
            var label = "Read file fixture.toml";
            var hash = Convert.ToHexString(SHA256.HashData(Encoding.UTF8.GetBytes(label)));
            var owned = new ExecutionCallRow(JsonSerializer.SerializeToElement(new { tool_name = "read_file", activity_label = label, activity_label_source = "tool",
                title_text = new { schema_version = 1, code = "tool.read_file", args = new[] { "fixture.toml" }, text_hash = hash } }));
            check(owned.OriginalLabel == label && owned.Title.Contains(UiText.Format("OwnedTool_read_file", "fixture.toml")), "trusted generated text uses selected resource: " + locale);
        }
        UiText.ApplyPreference(UiText.SimplifiedChinesePreference);
    }
}
