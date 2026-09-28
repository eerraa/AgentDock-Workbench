using System.Text.Json;
using System.Windows;
using System.Windows.Automation;
using System.Windows.Controls;
using CheckBox = System.Windows.Controls.CheckBox;
using ComboBox = System.Windows.Controls.ComboBox;
using Orientation = System.Windows.Controls.Orientation;

namespace AgentDock.ControlPanel;

internal sealed class PermissionSettingsEditor
{
    private readonly JsonElement _policy;
    private readonly CheckBox _edit = new() { Content = UiText.Get("PermissionProfileEdit"), Margin = new Thickness(0, 14, 0, 6) };
    private readonly CheckBox _inherit = new() { Content = UiText.Get("PermissionProfileInherit"), Margin = new Thickness(0, 6, 0, 6) };
    private readonly StackPanel _fields = new();
    private readonly StackPanel _categories = new();
    private readonly ComboBox _filesystem;
    private readonly ComboBox _network;
    private readonly ComboBox _boundary;
    private readonly ComboBox _approval;
    private readonly ComboBox _reviewer;
    private readonly Dictionary<string, CheckBox> _gates = new();
    private string _scope = "global";
    internal StackPanel View { get; } = new();

    internal PermissionSettingsEditor(JsonElement policy)
    {
        _policy = policy;
        View.Children.Add(_edit); View.Children.Add(_inherit); View.Children.Add(_fields);
        AutomationProperties.SetAutomationId(_edit, "PermissionSettingsEdit");
        AutomationProperties.SetAutomationId(_inherit, "PermissionSettingsInherit");
        _filesystem = Choice(UiText.Get("PermissionFilesystemLabel"), "PermissionFilesystem", [new("deny", UiText.Get("PermissionFilesystemDeny")), new("read", UiText.Get("PermissionFilesystemRead")), new("write", UiText.Get("PermissionFilesystemWrite"))]);
        _network = Choice(UiText.Get("PermissionNetworkLabel"), "PermissionNetwork", [new("allow", UiText.Get("PermissionNetworkAllow")), new("deny", UiText.Get("PermissionNetworkDeny"))]);
        _boundary = Choice(UiText.Get("PermissionBoundaryLabel"), "PermissionSandboxBoundary", [new("none", UiText.Get("PermissionBoundaryNone")), new("workspace", UiText.Get("PermissionBoundaryWorkspace"))]);
        _fields.Children.Add(Note(UiText.Get("PermissionBoundaryNotice")));
        _approval = Choice(UiText.Get("PermissionApprovalLabel"), "ApprovalPolicy", [new("on-request", UiText.Get("PermissionApprovalRequest")), new("never", UiText.Get("PermissionApprovalNever")), new("granular", UiText.Get("PermissionApprovalGranular"))]);
        _fields.Children.Add(_categories);
        foreach (var (key, label) in new[] { ("file_writes", UiText.Get("PermissionCategoryFiles")), ("commands", UiText.Get("PermissionCategoryCommands")), ("network", UiText.Get("PermissionCategoryNetwork")), ("mcp", UiText.Get("PermissionCategoryMcp")), ("management", UiText.Get("PermissionCategoryManagement")), ("other", UiText.Get("PermissionCategoryOther")) })
        {
            var box = new CheckBox { Content = UiText.Get("PermissionApprovalPrefix") + label, Margin = new Thickness(8, 3, 0, 3) };
            AutomationProperties.SetAutomationId(box, "ApprovalCategory_" + key);
            _gates.Add(key, box); _categories.Children.Add(box);
        }
        _reviewer = Choice(UiText.Get("PermissionReviewerLabel"), "ApprovalReviewer", [new("user", UiText.Get("PermissionReviewerUser")), new("auto_review", UiText.Get("PermissionReviewerAuto"))]);
        _fields.Children.Add(Note(UiText.Get("PermissionReviewerNotice")));
        _edit.Checked += (_, _) => UpdateEnabled(); _edit.Unchecked += (_, _) => UpdateEnabled();
        _inherit.Checked += (_, _) => UpdateEnabled(); _inherit.Unchecked += (_, _) => UpdateEnabled();
        _approval.SelectionChanged += (_, _) => UpdateEnabled();
        SelectScope("global", "");
    }

    private static TextBlock Note(string text) => new() { Text = text, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(0, 7, 0, 6) };
    private ComboBox Choice(string title, string id, ExecutionChoice[] values)
    {
        _fields.Children.Add(Note(title));
        var box = new ComboBox { ItemsSource = values, DisplayMemberPath = "Title", SelectedValuePath = "Id", MinHeight = 34 };
        AutomationProperties.SetAutomationId(box, id); _fields.Children.Add(box); return box;
    }
    private static string Value(ComboBox box) => (box.SelectedItem as ExecutionChoice)?.Id ?? "";
    private static void Select(ComboBox box, string value, string fallback) => box.SelectedValue = value.Length == 0 ? fallback : value;

    internal void SelectScope(string scope, string id)
    {
        _scope = scope; _edit.IsChecked = false; _edit.IsEnabled = scope != "conversation";
        var settings = _policy.Field("settings"); var overridden = false;
        foreach (var existing in _policy.Array("scopes"))
            if (existing.Text("kind") == scope && existing.Text("id") == id && existing.Field("settings").ValueKind == JsonValueKind.Object)
            { settings = existing.Field("settings"); overridden = true; }
        _inherit.IsChecked = scope == "workspace" && !overridden;
        var profile = settings.Field("permission_profile");
        Select(_filesystem, profile.Text("filesystem"), "write");
        Select(_network, profile.Text("network"), "allow");
        Select(_boundary, profile.Text("sandbox_boundary"), "none");
        Select(_approval, settings.Field("approval_policy").Text("mode"), "on-request");
        Select(_reviewer, settings.Text("approval_reviewer"), "user");
        var categories = settings.Field("approval_policy").Field("granular");
        foreach (var (key, box) in _gates) box.IsChecked = categories.Flag(key);
        UpdateEnabled();
    }

    private void UpdateEnabled()
    {
        _inherit.Visibility = _scope == "workspace" ? Visibility.Visible : Visibility.Collapsed;
        _inherit.IsEnabled = _edit.IsChecked == true;
        _fields.IsEnabled = _edit.IsChecked == true && !(_scope == "workspace" && _inherit.IsChecked == true);
        if (_approval is not null) _categories.Visibility = Value(_approval) == "granular" ? Visibility.Visible : Visibility.Collapsed;
    }

    internal void AddChange(Dictionary<string, object> change)
    {
        if (_edit.IsChecked != true || _scope == "conversation") return;
        if (_scope == "workspace" && _inherit.IsChecked == true) { change["inherit_settings"] = true; return; }
        var approval = new Dictionary<string, object> { ["mode"] = Value(_approval) };
        if (Value(_approval) == "granular") approval["granular"] = _gates.ToDictionary(item => item.Key, item => item.Value.IsChecked == true);
        change["settings"] = new Dictionary<string, object> {
            ["permission_profile"] = new { filesystem = Value(_filesystem), network = Value(_network), sandbox_boundary = Value(_boundary) },
            ["approval_policy"] = approval, ["approval_reviewer"] = Value(_reviewer)
        };
    }
}
