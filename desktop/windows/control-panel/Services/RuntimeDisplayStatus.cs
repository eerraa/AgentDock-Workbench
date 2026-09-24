namespace AgentDock.ControlPanel;

// Presentation only: public transport readiness must never rewrite local
// Core liveness or /healthz. Both the full panel and live tray use this policy.
internal sealed record RuntimeDisplayStatus(string HeaderKey, string ServiceKey, string HealthKey, string BrushKey)
{
    internal static RuntimeDisplayStatus From(RuntimeSnapshot? value)
    {
        if (value is null) return new("StatusUnavailable", "Unknown", "Unknown", "WarningBrush");
        var publicDown = value.TunnelMode.ToLowerInvariant() switch {
            "quick" or "named" => !value.CloudflaredRunning || string.IsNullOrWhiteSpace(value.PublicOrigin),
            "funnel" => value.Tailscale is not { Ready: true },
            _ => false
        };
        return new(value.Healthy ? (publicDown ? "LocalHealthyPublicUnavailable" : "RunningNormally")
                : value.CoreRunning ? "RunningHealthFailed" : "Stopped",
            value.CoreRunning ? "Running" : "Stopped", value.Healthy ? "Healthy" : "Unavailable",
            value.Healthy ? (publicDown ? "WarningBrush" : "SuccessBrush") : value.CoreRunning ? "WarningBrush" : "SecondaryText");
    }
}

internal sealed class RuntimeObservationOrder
{
    private DateTimeOffset _latest = DateTimeOffset.MinValue;
    internal bool Accept(DateTimeOffset startedAt)
    {
        if (startedAt < _latest) return false;
        _latest = startedAt;
        return true;
    }
}
