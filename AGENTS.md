# AgentDock development rules

This repository is the user's extension fork. Keep changes on the current feature
branch. Never publish to the upstream author's repository.

## Before editing

1. Load `agentdock_context` for this working tree. Check Git status, the current
   task, and local rules. Read the complete files being changed.
2. Find the existing responsibility and tests. Make a bounded implementation plan
   before changing code. Do not create parallel runtimes, task systems or themes.
3. Preserve unrelated changes. Do not reset, stash or overwrite another agent's work.

## Implementation boundaries

- `Runtime.Call` / `callObserved` owns external tool admission and root call identity.
  Children inherit an immutable binding. Never infer a conversation from the latest
  task, selected window or account. Unattributed calls stay unattributed.
- Tool RPCs, internal spans and command processes have separate lifetimes. Count
  only roots in totals. Unknown timings and outcomes are nullable, never fake zero
  or success. Use monotonic durations inside a process.
- Persist redacted facts before presenting them. Keep parameters, diff previews,
  output, caches, queues and subscriptions bounded. Never recursively observe
  journal writes, UI queries, heartbeats or background probes as AI requests.
- Views render state; services own I/O and cancellation. Bind selected, checked and
  expanded states explicitly. Hover or keyboard focus is not business state.
- Theme colors, borders and sizing belong in shared resources. Do not set hardcoded
  brushes in event handlers. Keep default and focused control geometry identical.
- All I/O must have cancellation, errors and resource cleanup. Do not synchronously
  block the UI with `.Wait()` / `.Result`; use `async void` only for event handlers.
- Keep one canonical implementation. Name files and types by responsibility, not
  New, Final or V2. Put necessary old-data conversion at a documented boundary.
- Keep authentication, permission checks, installation rollback and configuration
  ownership intact. A display setting must never restart the running core.

## Verification and delivery

Run Go formatting/tests and Windows desktop tests. Add regression tests for each
changed behavior; do not remove tests or weaken assertions to obtain a green build.
Use injected time for activity boundaries. Verify concurrent conversations, calls
without tasks, old records, rejected writes and asynchronous command completion.

Canonical implementation and delivery record:
`docs/eerraa/implementation-plan-1.1.6100.ko.md`. Its historical filename is not a
release version. The fixed source baseline is the Workbench v1.1.8 prerelease tag,
`4bd778d4077bbe58cfe19e4abb777f660694377b`, merged onto the earlier v1.1.7
(`b367eaab95202873fb213b8713440bf7822878c4`) fork line. Keep only pinned rg,
Korean presentation and reproducer-proven minimal fixes at the existing owners. Do
not reintroduce composite hosts, Origin pipes, epochs or parallel state systems.

This fork is updated only by running a new Setup. The tray/window update entry
opens the fork Releases page, and `agentdock update` refuses online checks and
downloads; only the maintenance `--local-archive` path remains. Do not re-enable
the online updater or publish updater archives as a user update path.

The 2026-09-28 user request explicitly authorizes normal main integration and
push, a new consistent version/tag, CI dispatch and a formal Windows x64 Release
only in `eerraa/AgentDock-Workbench`. Prefer the existing Windows package workflow.
This supersedes the old source-only restriction, but not any integrity check.
Never push upstream, rewrite shared history, move public tags, overwrite released
bytes, change open PR heads or publish incomplete product work as finished.

Never run Setup or install/upgrade/uninstall/recover the operational PC. Never
replace/restart production Core or change production Cloudflare/Tailscale or
credentials. Real installer acceptance requires a genuine isolated Windows runner
or disposable VM; do not impersonate one through environment variables. Local
explicit disposable unit/native fixtures are allowed. Preserve runner-only test
guards and distinguish a local build from an unexecuted runner test.

A new distribution must compare above the published 1.1.16102 and 1.1.17100
releases; never reuse a version for new bytes or confuse it with the upstream
baseline.
Source, tests, packaging, installation, publication and actual asset redownload
verification are separate delivery states. A formal release requires all of them,
including Korean and rg content in the actual payload, immutable source identity,
checksums, honest signing status, draft=false and prerelease=false. The release
approval is already explicit; missing tests are not missing user approval.

Keep other worktrees, uncommitted evidence and open PR branches intact. Preserve
failed tests and blocked-request messages. Do not retry an explicit safety or
permission denial through another tool, encoding or argument decomposition; only
independent permitted work may continue. Checkpoint the exact remaining boundary.
Keep the 120s activity, 180s request eligibility and 300s insertion expiry separate.
