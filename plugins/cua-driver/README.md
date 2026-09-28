# CUA 데스크톱 제어 플러그인

AgentDock에 연결된 에이전트(예: ChatGPT)가 이 Windows PC의 화면을 보고 앱을 조작하게 한다. 화면 캡처와 입력은 로컬 [CUA](https://github.com/trycua/cua) `cua-driver`가 하고, 이 플러그인은 그것을 AgentDock에 연결한다.

| 구성 | 내용 |
|---|---|
| `mcp.json` | MCP 서버 `cua-driver`. `cua-driver mcp --socket \\.\pipe\cua-driver`로 실행 중인 cua-driver 데몬에 붙는다. 도구 이름은 `cua-driver:<도구>` |
| `skills/cua-desktop/SKILL.md` | 에이전트용 사용법: 관찰 → 동작 → 검증, 승인 처리, 안전 규칙 |

## 1. 전제: cua-driver

cua-driver는 AgentDock에 포함되지 않는다. [CUA 공식 문서](https://cua.ai/docs)의 Windows 설치 방법으로 설치한 뒤 아래를 확인한다.

```powershell
cua-driver --version
cua-driver autostart status  # 자동 실행 등록과 데몬 실행 여부
cua-driver autostart enable  # 로그온 시 데몬 자동 실행 등록(한 번만)
cua-driver autostart kick    # 지금 바로 데몬 시작
```

데몬이 없으면 플러그인은 설치돼도 도구 목록을 가져오지 못한다. 데몬이 관리자 권한으로 돌면 관리자 창도 조작할 수 있고, 일반 권한이면 관리자 창에는 입력이 전달되지 않는다.

## 2. 자동 설치 (기본)

AgentDock Setup이 설치를 마친 뒤 이 플러그인을 `%USERPROFILE%\.agentdock\plugins\cua-driver`에 넣는다. 다음 Setup은 새 버전으로 갱신한다.

- 사용자가 끈 플러그인은 꺼진 채로 갱신한다.
- 사용자가 삭제하면 이후 Setup이 다시 넣지 않는다.
- 같은 이름을 사용자가 직접 설치했으면 Setup이 건드리지 않는다.

기록은 `%USERPROFILE%\.agentdock\plugins\.bundled.json`에 남는다. 실패해도 AgentDock 설치는 유지되고 Setup 경고(`bundled-plugins-deferred`)만 남는다.

## 3. 수동 설치

플러그인 원본은 이 저장소의 `plugins\cua-driver`, 또는 Release ZIP(`agentdock_windows_amd64.zip`)의 `share\agentdock\plugins\cua-driver`다.

| 방법 | 절차 | 이후 Setup 갱신 |
|---|---|---|
| Setup과 같은 방식 | `agentdock plugin bootstrap --bundle <plugins 폴더의 절대 경로> --home "$env:USERPROFILE\.agentdock"` | 된다 |
| 폴더 복사 | `cua-driver` 폴더를 `%USERPROFILE%\.agentdock\plugins\cua-driver`로 복사 → 제어판 "기능과 플러그인"에서 새로 고침 | 안 된다(직접 설치로 본다) |
| 채팅 에이전트 | 에이전트에게 `plugin_manage`를 `{"action":"install","source":"<cua-driver 폴더의 절대 경로>","confirmed":true}`로 호출하게 한다 | 안 된다(직접 설치로 본다) |

직접 설치한 플러그인을 Setup 관리로 바꾸려면 "기능과 플러그인"에서 삭제한 뒤 첫 번째 방법을 쓴다. `.bundled.json`에 이미 기록돼 있으면 삭제로 간주되므로 그 항목도 지운다.

## 4. 확인과 사용

1. 제어판 "기능과 플러그인"에 `cua-driver`가 켜져 있는지 본다.
2. 에이전트에게 `cua-driver:list_windows`와 `cua-driver:get_window_state`로 창을 보게 한다. 스크린샷이 보이지 않는다고 하면 스킬의 `view_image` 대안을 쓰게 한다.
3. 권한 모드가 `rules`면 동작마다 승인을 묻는다. 승인할 때 "이 작업공간에서 허용"을 고르면 그 도구 하나(`cua-driver:click` 등)만 허용된다. 승인된 호출의 결과는 에이전트에게 돌아가지 않으므로, 에이전트는 다시 관찰한 뒤 진행한다.

## 5. 주의

- 원격 에이전트가 이 PC의 앱을 조작한다. 화면의 웹페이지나 문서 내용이 에이전트를 속이려 할 수 있으므로, 결제·전송·삭제 같은 작업은 승인 단계에서 확인한다.
- 끄기는 "기능과 플러그인"의 스위치, 제거는 같은 화면의 삭제를 쓴다.
