using System.Text;

namespace AgentDock.ControlPanel;

public static class ExecutionTitleFormatter
{
    public static string Format(string tool, string label, string action = "")
    {
        label = string.Join(" ", label.Split((char[]?)null, StringSplitOptions.RemoveEmptyEntries));
        while (tool.Length > 0 && label.StartsWith(tool, StringComparison.Ordinal) &&
               (label.Length == tool.Length || " ·:-：".Contains(label[tool.Length])))
            label = label[tool.Length..].TrimStart(' ', '·', ':', '-', '：');
        if (tool == "file_edit" && label.StartsWith("EDIT_FILE", StringComparison.Ordinal))
            label = label[9..].TrimStart(' ', '·', ':', '-', '：');
        if (label.Length == 0) label = DefaultAction(tool, action);
        return tool.Length == 0 ? label : tool + " · " + label;
    }

    private static string DefaultAction(string tool, string action) => (tool, action) switch
    {
        ("task_manage", "create") => "创建任务",
        ("task_manage", "resume") => "恢复任务",
        ("task_manage", "checkpoint") => "更新任务进度",
        ("task_manage", "complete") => "完成任务记录",
        ("task_manage", "list" or "get") => "查看任务",
        ("task_manage", "cancel") => UiText.Get("ActivityCancelTask"),
        ("task_manage", _) => UiText.Get("ExecutionManageTask"),
        ("file_edit", "add") => "创建文件",
        ("file_edit", "delete") => "删除文件",
        ("file_edit", "move") => "移动文件",
        ("file_edit", _) => UiText.Get("ExecutionTool_file_edit"),
        ("read_file", _) => UiText.Get("ExecutionTool_read_file"),
        ("list_dir", _) => UiText.Get("ExecutionTool_list_dir"),
        ("search_text", _) => UiText.Get("ExecutionTool_search_text"),
        ("agentdock_context", _) => "加载上下文",
        ("exec_command", _) => "执行命令",
        ("session_observe", _) => UiText.Get("ExecutionTool_session_observe"),
        ("session_act", _) => "操作命令会话",
        ("workspace_manage", _) => "管理工作区",
        ("device_manage", _) => "管理设备",
        ("file_publish", _) => "发布文件",
        ("mcp_search", _) => "搜索扩展工具",
        ("mcp_inspect", _) => "查看扩展工具说明",
        ("mcp_tool_call", _) => "调用扩展工具",
        ("plugin_load", _) => "展开插件",
        ("insertion_ack", _) => "确认收到补充",
        ("plugin_manage", _) => UiText.Get("ExecutionTool_plugin_manage"),
        _ => tool.Contains(':') || tool.Contains('.') ? "调用扩展工具" : "执行工具"
    };
}
