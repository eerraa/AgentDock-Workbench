# AgentDock 1.1.6 → 1.1.16101 병목 최소화 구현 계획

문서 ID: AD-116-6100-IMPLEMENTATION-20260924
정본: `D:\Engineering\agentdock-eerraa\docs\eerraa\implementation-plan-1.1.6100.ko.md`
목표: 검증·절차·대기 비용을 최소화하여 구현과 최종 설치파일 제작을 최대한 빠르게 완료한다.
상태: 1.1.6 고정 소스에 적용할 설계·실행 계획이다. 제품 구현·시험·설치 완료를 뜻하지 않는다. §3의 미확인 항목은 해당 단위 착수 직전에 확인하며 독립 단위의 진행을 막지 않는다.

## 1. 기준, 목표, 권한

| 항목 | 실행 기준 |
|---|---|
| 제품 기준 | **AgentDock 1.1.6** |
| BASE | `df7c22f64438ec317e0518eec44d035ad98be2b2` |
| 소스 출처 | `A-m-o-r-F-a-t-i/agentdock`의 `feat/1.1.6-release` |
| 최종 제품 | **1.1.16101 / DownstreamRevision 101 / UpstreamVersion 1.1.6 / Distribution eerraa** |
| 버전 정책 | **main에서만 명시적으로 1.1.16101을 부여한다.** UpstreamVersion=1.1.6과 DownstreamRevision=101은 독립 메타데이터이며 버전 산식으로 추론하지 않는다. 기능 작업 브랜치는 upstream 1.1.6을 유지한다 |
| 작업 저장소 | `D:\Engineering\agentdock-eerraa` |
| 원격 | origin=`eerraa/agentdock`, upstream=`A-m-o-r-F-a-t-i/agentdock` |
| 브랜치 시작점 | 모든 기능 브랜치를 같은 BASE에서 새로 만든다. BASE와 승인된 기능 범위 밖의 코드를 일괄 이식하지 않는다 |
| main 처리 | 검증된 1.1.6 통합 소스를 정상 merge하고 기존 이력·문서 변경·자체 게시 권한 경계를 보존한다. reset/rebase/강제 갱신과 원격 변경은 하지 않는다 |

구현 시작 때 원격 main·v1.1.6 태그·기준 브랜치를 한 번 대조한다. 같은 SHA는 그대로 사용한다. 공식 v1.1.6의 SHA가 다르면 달라진 코드의 영향만 검토하고 이 문서의 BASE를 확정한 뒤 진행한다. BASE 고정 후에는 실제 불일치나 사용자의 기준 변경 요청이 없는 한 반복 fetch·Release 조회를 하지 않는다. 개발 브랜치와 공개 Release는 실제 원격 상태대로 구분하며, 명시적 요청 없이 다른 제품 버전으로 확장하지 않는다.

### 세션 경계

문서 수정 요청에서는 이 정본만 갱신하며 제품 소스·시험 코드·index·main·버전·운영 설정은 변경하지 않는다. 제품 시험·커밋·빌드·설치를 실행하지 않는다.

구현 세션의 허용 범위는 해당 단위의 최소 진입 확인, 기능별 로컬 브랜치 구현, §8의 최소 필수 검증, 로컬 커밋, 지정된 오프라인 Setup 제작이다. main의 마지막 다운스트림 릴리스 커밋에서만 제품 버전을 1.1.16101으로 변경한다. 실기기/VM 설치·운영 교체·push·PR·릴리스 게시를 하지 않는다.

PR은 사용자가 **최종 1.1.16101 Setup**을 직접 설치해 시험한 뒤 별도 세션에서만 연다.

### 신속 개발 원칙 — 검증·병목 최소화

**이 개발의 최우선 실행 목표는 검증과 절차상의 병목을 최소화하여 1.1.16101 구현과 설치파일 제작을 최대한 빠르게 완료하는 것이다.** 검증 범위·보고서·브랜치 수를 늘리는 것이 목표가 아니다. 요구한 기능과 실패 처리 계약은 유지하되, 불필요한 검토와 중복 실행을 없앤다. 각 절의 검토·회귀 목록은 위험과 확인할 동작을 뜻하며, 모든 환경·상태·언어의 조합을 매번 실행하라는 뜻이 아니다. 구체적인 실행량은 §8을 따른다.

- **확인한 것은 재사용한다.** 같은 소스·의존성·실행 조건에서 얻은 실제 읽기/검증 결과는 다시 만들지 않는다. BASE·작업 트리·도구 입력은 시작 때 묶어 확인하고, 이후에는 변경된 부분과 최종 산출물만 확인한다. upstream 문서의 pass를 우리 실행 결과로 복사하지는 않는다.
- **검토는 구현에 필요한 만큼만 한다.** 수정할 파일 전체와 직접 호출자·인접 시험을 한 번 확인한 뒤 구현한다. 저장소 전수 재감사, 무관한 이력 조사, 범위 밖 결함 탐색, 대안 설계의 반복 비교로 구현 착수를 늦추지 않는다. 재구현은 책임 중복을 없애는 최소 범위이며 새 프레임워크를 만드는 기회가 아니다.
- **단위가 막혀도 전체를 멈추지 않는다.** G0를 모든 기능의 선행 관문으로 두지 않는다. 실제 소유권·데이터 손실 위험을 판단할 수 없는 단위만 보류하고 독립 단위는 계속한다. 보안/권한 차단은 우회하지 않으며 같은 실패를 자동 반복 호출하지 않는다.
- **승인된 범위는 연속 실행한다.** 이미 승인된 기능 구현·표적 검사·로컬 커밋·패키징 사이에 재승인을 요청하지 않는다. 한 구현 세션에서 완성하는 흐름을 기본으로 하며, 실제 도구/세션 한계로 멈춘 경우 마지막 확인 지점에서 재개한다. 새로운 파괴 작업·권한 경계·요구 충돌만 사용자 판단을 받는다. 한글화 단일 브랜치와 독립 PR 경계는 유지한다.
- **검증은 변경 위험에 비례한다.** 기능 단위가 완성될 때 관련 compile/정적 검사와 핵심 회귀를 묶어 한 번 실행한다. 파일·커밋마다 전체 suite를 돌리지 않는다. 실패 후에는 원인을 수정하고 실패 slice와 영향 경계만 재실행한다. 변경 없는 영역의 시험을 새로 만들거나 입증된 동작을 반복 측정하지 않는다.
- **검증 인프라가 개발 목표가 되지 않게 한다.** 기존 fixture·cache·runner를 우선 사용한다. 광범위한 OS/테마/DPI/언어 조합, 실시간 장기 대기, 별도 성능 캠페인, 새 VM/브라우저 자동화 환경 구축은 기본 검증에서 제외한다. 필요한 시간 경계는 fake clock과 작은 fixture로 검사한다. 필수 입력이 없으면 한 번 확인 후 무관한 구현을 계속하되, 최종 필수 검사는 미통과인 채 완료로 표시하지 않는다.
- **보고와 문서도 최소화한다.** 기능 단위의 주요 checkpoint, 실행한 검사·실패·미실행, 최종 SHA/설치파일만 간단히 남긴다. 중복 계획·삭제 이력·증거 사본·검증 전용 대형 보고서를 만들지 않는다. 문서만 바꿨으면 문서 diff·참조 확인만 하고 제품 build/test는 하지 않는다.

최소화 대상은 개발 중 검증의 반복·범위와 대기 비용이다. 제품의 권한/소유권 확인, 설치 실패 전파·복원, 번들 무결성, 사용자 데이터·ACL 보존, 표준 설치의 시작→건강 순서를 생략하거나 오류를 경고 성공으로 바꾸지는 않는다. 실제 미실행/실패/차단은 그대로 보고한다. 실기기·VM 설치는 사용자에게 남기고, push·PR은 별도 승인 단계로 유지한다.

## 2. 1.1.6의 재사용 기반과 구현 범위

### 2.1 근거 적용 원칙

아래 경로는 BASE의 재사용 기반과 직접 확인할 소스·시험을 지정한다. 동일한 소스·의존성·실행 조건에서 확인한 결과는 재사용하고, 커밋 메시지·릴리스 설명·시험 소스의 존재를 시험 실행 성공으로 취급하지 않는다.

§3에는 실제 호출 경로와 구현 결정을 위한 미확인 항목을 둔다. 런타임 입구 삭제와 설치 완료 판정은 관련 확인 후에만 확정한다. 이 확인은 해당 기능의 착수 조건이며 독립 기능까지 기다리게 하는 전역 관문이 아니다.

### 2.2 1.1.6에서 이미 제공하는 기반 — 다시 만들지 않음

| 제공된 기반 | 직접 소스·시험 근거(df7c22f) | 이번 계획에서의 처리 |
|---|---|---|
| 컨텍스트 스냅샷·동일 키 공유 빌드·waiter 취소 분리 | `internal/app/context_snapshots.go`, `internal/snapshot/cache.go`, `context_preparation_test.go`, `context_snapshot_regression_test.go` | 기존 cache를 유지. 별도 컨텍스트 엔진·request identity 캐시를 만들지 않음 |
| 독립 workspace_context·정확한 scoped Skill 참조 | `internal/app/workspace_context.go`, `workspace_context_test.go`, `internal/tool/skill/scoped_reference.go` | 등록 도구·새 화면의 한글화 범위에 포함. 워크스페이스 이동 때 cwd나 인증을 바꾸지 않음 |
| 삽입의 text/structuredContent 이중 전달·표시 모드 정합성 | `internal/mcp/response_additions.go`, `presentation.go` 및 인접 시험 | 응답 조립·TextOnly 동작과 업무 결과 보존 계약을 유지 |
| 전체 MCP 도구 목록·일괄 inspect·호출 후 카탈로그 | `internal/app/specs_mcp.go`, `internal/mcp/client/catalog.go`, `catalog_response_test.go` | 도구와 schema를 보존. discovery는 BASE에서 실제 확인한 결함만 수정 |
| 실행 중 request·실제 response·progress의 저장과 페이지 읽기 | `internal/activity/payload.go`, `internal/app/execution_payload.go`, `execution_payload_test.go`, `ExecutionPayloadView.cs` | **payload 모델을 정본으로 한글화. 제품 summary와 실제 출력을 구분** |
| 활동 중/최근 대화 탐색·역사 정렬 snapshot·접기·도구막대 | `conversation_sidebar_test.go`, `conversation_history.go`, `ExecutionWindow.SidebarStream.cs`, `ExecutionWindow.History.cs` | 새 UX를 유지. 메인페이지 요약 갱신과는 별개로 취급 |
| 편집 트랜잭션의 실제 증감 통계 | `internal/tool/file/edit_outcome.go`, `edit_statistics_test.go`, `git_patch_stage.go`, WSL transaction 시험 | 트랜잭션 통계·부분 복원·unknown과 실제 등록 도구명 표시를 유지 |
| Plugin source adapter·Skill 읽기 lease | `internal/plugin/source_types.go`, `adapters.go`, `internal/skill/state/read_lease.go` | 기존 Store·버전·activate/rollback·Heavy 경계 유지. 별도 registry 만들지 않음 |
| Windows native host·Setup 진단·비동기 터널·구형 layout migration | `cmd/agentdock-shim/task_core_host_windows.go`, `setup_runtime_host_windows.go`, `internal/installer/tunnel_start_windows.go`, `internal/selfupdate/legacy_migration_windows.go` | 유효한 substrate는 재사용하되 단일 작업 소유권과 실제 caller를 재검토. 새 파일 존재를 연결 완료로 보지 않음 |

request의 실행 중 저장·response 재열기 보존, waiter별 취소, 120초 경계와 runtime-confirmed in-flight 동작은 BASE의 회귀 자산이다. 수정이 닿는 경계의 인접 시험만 재사용하며 이미 제공된 기능을 별도 구현하거나 PR로 다시 제출하지 않는다.

### 2.3 구현할 기능과 확인할 경계

| 항목 | 현재 판단 | 다음 조치 |
|---|---|---|
| 한국어 전체 | df7c22f에 ko-KR WPF 리소스가 없고 새 payload UI에도 제품 생성 중문이 있음 | 1.1.6 UI 전체를 한글화 브랜치 하나로 완성 |
| rg 내장 | df7c22f에 `internal/bundledrg`가 없음 | 현재 실행 세대에 귀속되는 검증 번들을 독립 구현 |
| 메인 활동 요약 | `MainWindow.Access.cs`의 요약 조회는 전체 refresh에 묶이고 Number는 누락을 0으로 취급 | 기존 집계 위 갱신 수명·typed parser만 교체 |
| 로컬 건강 표시 | 공개 접근 상태와 로컬 건강의 책임을 분리해야 함 | 기존 DTO 위 독립 표시 수정. native 역할 식별은 런타임 브랜치 책임 |
| 작업 디렉터리 ACL | BASE의 `internal/config/config.go`의 경로 준비에도 EnsurePrivate 호출이 남음 | 홈만 보호하는 경계·회귀를 독립 구현 |
| 단일 상승 소유자 | BASE의 native host 입구와 등록된 tray --run-core-task 연결, detached 터널, 표준 deferred 성공 처리를 별도 확인해야 함 | native host를 단일 소유자로 통합. G0 후 실제 입구 정리 |
| 저장 상태·MCP nil·Setup receipt | BASE가 무효 envelope·optional 인자·파일 공유 경합을 정확히 처리하는지 확인 필요 | 결함이 있는 경우만 최소 patch와 시험. 정상 처리하면 PR 자체 생략 |
| 설치 layout 이전과 rg | 1.1.6의 generation migration에 번들 component 보존을 연결해야 함 | migration의 실제 파일 처리와 복원 경계에 검증 component를 연결 |

**구현 원칙:** BASE의 native host·payload·도구막대·트랜잭션 통계를 유지하며 필요한 책임만 확장·통합한다. 한글화에 필수인 metadata는 같은 한글화 브랜치에 포함한다. 별도 런타임·OAuth hot-reload 엔진·설치 엔진을 추가하지 않는다.

## 3. G0 — 해당 단위 착수 직전에만 닫는 최소 추가 확인

미확인 부분을 정상으로 가정하지 않는다. 다만 G0-1~G0-5 전체를 먼저 끝내야 모든 구현을 시작할 수 있다는 전역 관문은 두지 않는다. R 착수 전에 G0-1~G0-3, C의 잔여 패치 판단 전에 G0-4, Q의 migration 연결 및 최종 패키징 전에 G0-5의 관련 부분만 확인한다. K/S/H/A는 자신의 수정 책임과 직접 의존성이 확인되면 먼저 진행할 수 있다.

동일 BASE에서 이미 읽은 경로는 재조사하지 않는다. 미확인 entry→owner→직접 caller·인접 시험을 한 번 묶어 읽고 구현 결정이 가능해지면 바로 구현한다. 삭제할 입구의 호출자 존재 여부와 설치 소유권 전이는 실제 확인하되, 이를 이유로 unrelated 파일·릴리스 전체를 재감사하지 않는다. 소스 확인은 시험 실행 성공을 뜻하지 않는다.

| ID | 확인할 경로 | 통과 기준 / 실패 시 처리 |
|---|---|---|
| G0-1 | `TaskAdminService.ElevatedCoreArguments/CreateElevatedTask` → stable tray shim의 인자 분기 → `runTaskCoreHost` 또는 WPF `RunElevatedCoreTaskAsync` → Core/supervisor 생성 | 실제 예약 작업이 어느 host를 실행하는지 확인. --task-core-host 파일 존재만으로 사용 중이라고 판정하지 않음. 하나의 실사용 host로 연결할 패치를 정하고 경쟁 host caller를 없앰 |
| G0-2 | `platformServiceAction/startCore/stopCore`, `StartInteractiveScheduledTask`의 service·selfupdate caller, `configureCloudflareTunnel`, `runCloudflaredOnce/applyQuickTunnelURL` | 외부 제어와 내부 Origin 적용 분리, 불필요한 COM Run/직접 Kill/runas/health 대기를 정확히 분류. 자기 task 종료 경로가 없도록 설계 검증 |
| G0-3 | `install.ps1` activation/commit/rollback, `TaskAdminService.RestoreBackup`, `Engine.commit/abandon`, `adapter_recovery.go`, `tunnel_start_windows.go`, `setup_runtime_host_windows.go` | 세대 포인터·task 정의·이전 instance 종료의 실제 commit 순서 확인. 표준 start→Wait 순서 유지, 상승/rollback 건강 대기 제거 경계 확정 |
| G0-4 | `execution_discovery.go`, managed override·insertion·notification store, Setup result reader와 인접 시험 | BASE의 구현과 인접 시험에서 nil/무효 envelope/공유 위반 처리를 확인. 정상 처리하면 제외하고 실제 결함만 독립 PR로 채택 |
| G0-5 | 1.1.6의 local-archive·legacy migration·generation journal과 패키징 script 실제 param | rg가 publish/repair/migration/rollback 때 빠지지 않도록 책임 배치. amd64·unsigned·명시 cloudflared/ISCC 전달을 실제 entry가 지원하는지 확인 |

확인 결과는 해당 단위의 구현 결정과 필요한 시험에만 반영하고 별도 검토 보고서를 만들지 않는다. 도구 차단이 계속되면 그 경로를 바꾸지 않고 독립 단위로 진행한다. 보류 항목이 최종 필수 계약이면 미완료를 분명히 남기며 전체 완성을 주장하지 않는다.

## 4. 독립 PR 브랜치

각 브랜치는 같은 고정 BASE에서 출발한다. 단위별로 자기 diff만 적용한 트리에서 §8의 해당 변경 최소 검사를 통과해야 한다. 모든 branch에서 동일한 전체 suite나 전체 package를 반복 실행하지 않는다. 구현 순서, Go/WPF/설치기 계층, P 번호 때문에 하나의 기능을 여러 PR로 쪼개지 않는다. 아래 이름은 다음 세션에서 만들 계획이며 현재 존재한다고 가정하지 않는다.

| 브랜치 | 포함하는 완결된 동작 | 분리 경계 |
|---|---|---|
| `feat/korean-localization` | 제어판·트레이·활동센터·신규 payload/역사/편집통계·브라우저 인증/status·MCP App·설치/제거 UI, 필요한 표시 metadata·저장·호환·시험 **전부** | 한글화 필수 backend를 별도 PR로 분리하지 않음. eerraa 브랜드·업데이트 정책 제외 |
| `feat/windows-bundled-rg` | pin/라이선스·검증·검색 선택·package·generation/migration/rollback·시험 | 운영/사용자/WSL PATH 전체 변경 금지 |
| `fix/windows-elevated-runtime` | task controller·단일 host/Job·Core/supervisor·Quick Origin·상승 설치/복원·관련 표준 실패 전파·native 역할 식별 | 같은 수명 계약이므로 한 PR. 번역·별도 UI 표시 개선 제외 |
| `fix/workspace-acl-boundary` | 홈 보호와 작업 디렉터리 ACL 불변 | 런타임 PR 없이도 적용 가능 |
| `fix/runtime-health-display` | 기존 모델 위 로컬 건강과 공개 터널 표시 분리 | Core 역할 신원 체계는 런타임 PR. 미수용 새 DTO를 숨은 전제로 쓰지 않음 |
| `fix/activity-summary-refresh` | 기존 overview API의 정확한 parsing·visible 자동 갱신·stale/최근 호출 표시·시험 | 집계 엔진/미열람 상태/한국어 번역 제외 |
| `fix/state-store-preservation` | 유효성 미확인 envelope 비파괴 거절·last-good 보존 | G0-4에서 남은 결함이 있는 store만. Setup receipt와 분리 |
| `fix/mcp-discovery-null-safety` | optional/typed-nil params·nil result의 실제 남은 결함 | G0-4 결과 이미 해결됐으면 생성하지 않음 |
| `fix/setup-receipt-read-race` | 공유/잠금 위반의 bounded pending 읽기·nonce/deadline 보존 | G0-4 결과 이미 해결됐으면 생성하지 않음. 설치 소유권과 분리 |
| `integration/release-1.1.6100` | 검증한 branch 조합. 제품 버전 변경과 최종 패키징은 main에서만 수행 | upstream PR 아님 |

모든 branch를 반드시 만들어 개수를 맞추는 목표가 아니다. 이미 upstream에서 해결한 항목은 생략한다. 독립 기능의 진짜 선행 코드가 필요하면 몰래 다른 PR을 merge하지 말고 수용 후 rebase하거나 같은 기능으로 묶는다.

**한글화는 끝까지 한 브랜치다.** 다른 기능 branch는 자기 의미에 필요한 기본 리소스 키만 추가한다. 한국어 resource·문구·generated presentation 처리는 모두 `feat/korean-localization`에서 한다. 통합 브랜치에 흩어진 한국어 hotfix를 남기지 않는다. 별도 기능의 새 문구는 같은 한글화 브랜치에서 담당하되, upstream 제출 diff는 제출 기반에 실제 존재하는 기능을 대상으로 다시 산출한다. 미수용 DTO·동작에 대한 숨은 의존성이 생기지 않게 resource 추가와 기능 연결을 구분한다. 거절된 기능 전용 잔여 키는 같은 한글화 브랜치에서 정리한다. 단독 한글화에서는 키/placeholder와 실제 표시·원문 보존의 최소 회귀를, 최종 조합에서는 충돌·새 연결 부분만 검증한다. 같은 전체 UI 행렬을 두 번 돌리지 않는다.

## 5. R — 상승 런타임: 1.1.6 substrate 위 단일 책임으로 재구현

### R1. 외부 task lifecycle

위치: `internal/desktopruntime/service_windows.go`, `service_startup_windows.go`, `task_scheduler_windows.go`, `task_com_windows.go`, service command entry, `RuntimeService.cs`.

`scheduledTaskController`(신규 책임)를 Go desktopruntime 내부에 둔다. task 정의/보안 조회는 기존 native COM 기반을 재사용하며 1.1.6의 VARIANT ABI 수정은 보존한다. lifecycle 명령 정책은 한 곳에만 둔다.

- `InspectOwnedTask`(신규): 정확한 AgentDock 이름/root/SID, InteractiveToken, Highest, 단일 안정 tray action·인자, Enabled, 기존 instance identity를 읽는다. 다른 소유자의 동명이인 task는 손대지 않는다.
- `EndOwnedTask`(신규): 같은 작업에 `schtasks /End`를 보내고 캡처한 이전 instance 종료를 bounded wait한다. 포트 소멸이나 이미지 이름을 종료 증거로 대신하지 않는다.
- `RunOwnedTask`(신규): 같은 작업에 `/Run` 요청을 보내고 수락 오류를 반환한다. 상승 시작은 새 포트를 기다리지 않는다.
- `ApplyScheduledAction`(신규): start/stop/restart와 configure handoff가 공유한다. 이미 중지인 정상 task만 멱등 stop 성공. 부재·권한 거부·다른 action은 오류다.

호출은 `패널/트레이 → native service action → platformServiceAction → controller`로 통일한다. C#에 별도의 /End→/Run 규칙을 다시 두지 않는다. 외부 runtime 실패의 agentdock.exe runas 재실행과 언어별 not-running 부분 문자열 판정은 제거한다. 명시적 설치/권한 전환의 TaskAdmin UAC까지 없애는 것은 아니다.

root별 기존 operation lock을 재사용한다. configure가 잡은 lock을 controller가 다시 획득하지 않게 public/locked 내부 연산을 나눈다. C# 클릭 방지 semaphore는 사용자 경험용이고 프로세스 간 정본 lock을 대신하지 않는다.

### R2. native 작업 host와 Job

**작업 host는 1.1.6의 `cmd/agentdock-shim/task_core_host_windows.go`를 우선 통합 대상으로 삼는다.** G0-1에서 실제 entry와 호출자를 확인하고 하나의 host로 연결한다.

작업 정의의 외부 계약은 **안정 tray `--run-core-task --runtime-root <root>`**를 유지한다. 설계 방향은 stable tray shim의 이 입구를 하나의 native host 책임으로 연결하고, 별도로 추가된 --task-core-host 및 WPF host와의 경쟁을 없애는 것이다. 어떤 이름을 alias로 남길지는 실제 호출자 검토 뒤 결정한다. 단순 인자 감지만으로 installer trial/실행 호환 검사를 우회하지 않는다.

host는 고정 세대, 한 Job, 현재 Core 자식, 최대 한 named/quick supervisor, 취소/종료 수명을 소유한다. cloudflared는 supervisor의 자식으로 같은 Job 안에 남는다. 새로운 예약 작업·Windows 서비스·공개 관리 서버는 없다.

현재 Start→Attach 간격을 남기지 않는다. `StartChildInJob`(신규 Windows 시작 helper)는 `STARTUPINFOEX`의 JOB_LIST와 필요한 HANDLE_LIST를 통해 자식 생성 때 Job과 필요한 pipe handle만 지정한다. Job handle은 불필요하게 상속하지 않는다. 귀속 실패/미지원이면 unowned 실행으로 폴백하지 않고 생성 실패로 반환한다. 무관한 command 실행 엔진 전체를 교체하지 않는다.

Core와 supervisor의 예상/미예상 종료를 구분하고 자식·pipe·Job handle을 회수한다. 필요한 supervisor 생성 실패를 mode=none처럼 성공 처리하지 않는다. 건강 미준비는 살아 있는 host 내부 상태로 유지하며 프로세스 crash나 Task Scheduler 재시작 백오프로 처리하지 않는다.

### R3. Quick Origin과 local health

대상: `configureCloudflareTunnel`, `startCloudflareTunnel`, `runCloudflaredOnce`, `applyQuickTunnelURL`, `invalidateQuickTunnelAfterExit`, `regenerateQuickTunnel`, effective runtime 설정과 task host.

configured mode와 ready를 분리한다. Quick 선택·URL 미확정·재연결 중에도 mode=quick이고 ready=false다. 이를 manifest none으로 바꿔 supervisor 생성이 빠지게 하지 않는다. host·UI·CLI가 modefile/manifest 우선순위를 각각 구현하지 않게 기존 정규화 경계 하나를 사용한다.

cloudflared 생성은 **해당 작업이 소유한 현재 Core가 `GET http://127.0.0.1:{실제 port}/healthz`에 성공한 뒤**에만 허용한다. 임의의 200 응답을 자기 Core라고 추정하지 말고 owner의 Core identity/수명과 함께 확인한다. 0.0.0.0을 probe 목적지로 쓰지 않는다. 초기/회복 health-pending은 stop 가능한 고정 재검사이며 cloudflared crashRetry를 올리지 않는다.

Quick URL 채택 때문에 자기 작업에 /End를 보내지 않는다. `applyQuickTunnelURL → host-private Origin request → replaceCoreOrigin`(신규)으로 **같은 host가 Core 자식만 교체**한다. Job/task/supervisor/cloudflared 수명은 유지한다. OAuth 서버 전체를 hot reload하는 별도 엔진은 만들지 않는다.

통신은 host와 자신이 생성한 supervisor 사이의 익명 pipe로 한정한다. 작은 UTF-8 JSON frame, 상한 16 KiB, 한 in-flight 요청, schema/request ID/owner epoch/고정 generation, `replace_core_origin` 한 연산만 허용한다. 입력에 executable·명령·임의 파일 경로를 받지 않는다. stdout 프로토콜과 회전 로그/stderr는 구분하고 pipe EOF·기한·취소·중복 ID를 명시 처리한다. 이미 존재하는 동등한 private 제어 경계가 G0에서 확인되면 그것을 확장하고 두 채널을 만들지 않는다.

Origin 파일의 writer는 host 하나다. 새 Core가 같은 세대·기대한 Origin으로 시작하고 소유 Core의 health가 확인된 뒤 ack한다. 그 뒤 supervisor가 public/ready를 게시한다. 실패는 ready=false이며 일반 성공으로 바꾸지 않는다. 새 Origin 반영 중의 오류, 설정 오류, local health 실패, 실제 cloudflared child_exit, 사용자 stop을 서로 다른 결과로 분류한다. crash backoff는 실제 child_exit에만 적용한다.

### R4. 상승 configure·자동 시작

하나의 operation gate 아래 `전체 입력 검증 → 설정 기록 → 기존 task /End → 기존 instance 종료 확인 → 같은 task /Run → 반환`이다. 새 포트·Quick URL·공개 터널을 기다리지 않는다. 실패한 validation 때문에 token·설정을 부분 변경하지 않으며, handoff 실패 때 이전 설정 복원 여부를 명확히 보고한다.

상승 경로의 독립 cloudflared HKCU Run, 별도 detached start proxy, configure 전후 직접 cloudflared Kill/launch를 없앤다. 외부 `tunnel start/launch`가 상승 소유자 밖의 supervisor를 만들지 못하게 한다. 필요한 표준 모드 caller까지 제거하지 않는다. Tailscale 공식 서비스는 이 Cloudflare 소유 모델의 대상으로 확대하지 않는다.

### R5. 설치·rollback 완료 판정

대상: `install.ps1`의 task 준비/활성화/commit/rollback, `TaskAdminService`의 저장·정의·복원, 기존 `internal/installer` transaction/active pointer/adapter verification.

기존 generation·파일 journal은 유지한다. **파일과 task 소유권 전이만 재구현**한다. TaskAdmin helper가 복원 안에서 몰래 Run하고 바깥 스크립트가 다시 Stop/Run하는 중복 정책은 없앤다.

상승 성공 경로:
`Preflight → SnapshotFilesAndTask → EndOldInstance → StageNewGeneration → BindSameTask → VerifyFilesAndTaskOwner → Commit/Finalize → InstalledNotStarted`.

완료 증거는 transaction ID와 target generation 파일·active pointer, 같은 AgentDock task의 SID/InteractiveToken/Highest/안정 tray action/root/Enabled, 이전 instance 종료다. source shim과 generation의 관계를 검증하지 않은 owner_verified 플래그만으로 성공하지 않는다. 실패 가능한 외부 작업이 끝나기 전에 journal을 버리지 않는다. logon trigger가 미커밋 generation을 실행할 창을 만들지 않도록 Disabled 준비·commit·Enabled 전환 순서와 복원 조건을 G0-3에서 확정한다.

상승 설치는 임시 추출본 service start나 120초 health 대기에 의존하지 않는다. no-start/skip-health는 **owner 검증을 통과한 이 계약의 결과**이며 소유자 없는 성공 탈출구가 아니다. 요청된 상승 helper 실패를 자동 standard 성공으로 바꾸지 않는다.

표준 설치는 **`service start → Wait-AgentDockHealth`의 상대 순서**를 유지한다. service start가 실패하면 후속 Wait를 호출하지 않고 rollback한다. 건강 실패를 runtime-launch-deferred 경고로 삼키지 않는다. 공개 터널 준비 비동기화라는 1.1.6 개선과 로컬 Core 건강 실패를 혼동하지 않는다.

실패 경로:
`Failure → EndNewOwnerIfAny → RestoreGeneration/Files → RestoreTaskDefinition/Security → VerifyRestoredOwner → OptionalSingleResume → SealRollback`.

이전 실행 재개는 한 곳에서 최대 한 번이다. 상승이면 같은 task /Run이며 새 건강 대기는 없다. 표준 복원도 두 번째 장시간 start/health 대기를 만들지 않는 비대기 시작 책임을 사용한다. 상승/rollback의 별도 tunnel start와 abandon --require-health를 없앤다. 일반 Engine API의 명시적 health 검증이나 표준 성공의 필요한 검증 자체를 일괄 삭제하지 않는다. 파일 복원·task 복원·실행 요청 수락·health는 별도 사실이며 확인하지 않은 health를 true로 기록하지 않는다.

복원 실패는 rollback-failed다. 실제 복원 journal은 보존하고 설치를 성공으로 표시하지 않는다.

### R6. 역할 식별·사용하지 않는 입구

상태 조회는 Go native 역할 판정으로 모은다. supervisor나 동일 이미지의 command/updater를 CoreRunning으로 세지 않는다. PID 하나나 경로 이름만 믿고 종료하지 않으며, host가 소유한 handle/생성시간/고정 세대와 연결한다. 제한 권한 QueryFullProcessImageName을 유지하고 query 실패는 stopped와 구별한다.

G0 후 C# task lifecycle·중복 process P/Invoke, 경쟁 native/WPF host, 상승 detached wrapper의 실제 caller가 없어지면 제거한다. `installerTrialHostEntry/liveInstallerTrial`, `StartInteractiveScheduledTask` 같은 이름은 **후보 검토 대상**이며 df7의 존재·caller를 다시 확인한다. 표준 install trial·local-archive migration·복구에 필요한 좁은 입구는 보존한다. 상승 설치에서만 필요 없어졌다는 이유로 모든 trial 지원이나 stable shim parent lifetime을 삭제하지 않는다.

## 6. K — 한국어 전체: 하나의 branch에서 새 1.1.6 UI에 구현

BASE의 리소스/UiText/ActivityText/ResourceManager를 확장하고 별도 테마나 렌더러를 만들지 않는다. 번역 사전은 현재 키·placeholder에 맞는 것만 재사용하며 XAML·ExecutionModels 구조를 유지한다.

### K1. 정확한 화면 인벤토리

제어판·트레이·권한·설정·확인/오류 대화상자, `workspace_context`/`mcp_tool_list` 등 새 도구의 제품 라벨, 신규 호출/출력 탭과 payload state/byte/line/page 안내, 편집 통계·부분/unknown 표시, 활동/역사/접기/도구막대, 브라우저 OAuth/status, MCP App, Inno 설치/복구/제거를 모두 포함한다.

BASE의 `ExecutionPayloadView.Position/StateLabel/Describe/ApplyPage/ReadFailed`에 남은 제품 문구를 리소스로 옮긴다. 새 내비게이션과 display 경로 호환도 그대로 둔다. locale 변경은 Core 재시작이나 업무 상태 변경이 아니다. 리소스 키/placeholder는 자동 검사로 전체를 확인하고, 레이아웃·focus·screen-reader 이름은 실제 변경 화면과 긴 문구의 대표 사례에 한정한다. 모든 화면×테마×DPI×언어의 조합을 만들지 않는다.

### K2. 제품 생성 설명과 실제 출력 분리

`describeExecution`, trusted root ingress, `localManagementStart/Finish`, permission.update·승인 거절/만료·복구 producer에서 필요한 code+인자/출처를 만들어 기존 Event→projection에 함께 저장한다. 사용자 인자가 임의 presentation/label_source를 주입하지 못하게 하고 redaction·크기 제한을 지킨다. 성공 설명이 뒤의 오류를 덮지 않도록 summary가 교체되면 descriptor도 교체/제거한다.

`加载上下文`, `修改执行权限`, `操作系统权限未改变。`는 각각 컨텍스트 불러오기·실행 권한 변경·운영체제 권한은 변경되지 않았습니다로 표시한다. 실제 scope/mode/revision 값, 경로·명령·질의·사용자 제목과 stdout/stderr는 데이터로 보존한다.

**1.1.6에서는 제품 summary를 출력 탭의 대체 결과로 쓰지 않는다.** CallTitle/CallDescription과 RequestPayload/ResponsePayload 표시를 분리한다. payload preview/페이지 원문은 재번역하지 않고 제품 state·빈 기록·오류 안내만 번역한다. 기존 기록도 출력 미저장/일부/불명 상태를 유지한다. 저장된 journal을 일괄 치환하지 않고 신뢰 가능한 producer 형식만 좁게 호환하며, 근거가 부족하면 원문을 남긴다.

### K3. 브라우저·MCP App·설치기

OAuth 동의·비밀번호·등록 redirect·CSP·HTML escaping과 TextOnly/App·insertion 응답 구조는 바꾸지 않는다. MCP App의 현재 pinned renderer를 확인한 뒤 정식 locale 입력이 있으면 그 경계를 사용한다. 없으면 기존 compatibility adapter를 한 경계로 제한하고 exact-anchor/resource 회귀를 재사용한다. renderer나 메시지 경계를 바꾼 경우만 기존 브라우저 fixture의 해당 smoke를 추가하며 단순 사전 추가에 별도 브라우저 환경을 구축하지 않는다. HTML adapter는 현재 pinned renderer의 실제 계약에 맞게 적용하며 치환 실패를 침묵시키지 않는다. 별도 protocol fork/두 번째 renderer는 만들지 않는다.

원본 Korean.isl·번역자 크레딧·Inno 라이선스·pin을 함께 보존한다. AgentDock 제품 문구는 설치기 리소스에서 처리하며 `scripts/install/install.ps1`는 ASCII다. 번역 문구를 비교해 권한·상승·복원 동작을 결정하지 않는다.

회귀: 실제 root/local-management 생성→저장→재열기→loaded WPF 라벨/설명/payload state; 사용자가 같은 중문을 쓴 반례; 실제 출력 미변형; 새 페이지·오류·history·edit stats; en/zh-CN/ko-KR 키·placeholder parity; 언어 전환이 삽입·권한·Core 수명을 바꾸지 않음. 라벨 한 개 번역이나 문자열 존재 검사만으로 완료하지 않는다.

## 7. 독립 수정과 다운스트림 정책

### S. 메인 활동 요약

`RuntimeExecutionOverview → CallStatistics`를 집계 정본으로 유지한다. 새 사이드바 streaming이 메인 요약 갱신을 해결했다고 가정하지 않는다.

`ActivityClient` 경계에 typed `ExecutionSummarySnapshot`/parser(신규)를 둔다. 필수 running/pending/unknown은 비음수 정수이며 missing/null/string/음수/지원하지 않는 schema는 invalid이지 0이 아니다. 기존 범용 Number helper를 전역 변경하지 않는다.

visible 메인창은 즉시 조회하고 3초 주기로 한 요청만 실행한다. 숨김·교체·닫힘에는 요청/구독을 취소하고 재표시 때 즉시 읽는다. 느린 Core/Nexus 전체 refresh나 건강 false→true를 전제로 하지 않는다. 활동창 overview도 같은 DTO/parsing 책임을 사용하되 창이 닫혀도 쓸모없는 polling을 계속하지 않는다. port/인증/root 세대가 바뀐 뒤 구 응답이 덮지 못하게 한다.

정상 숫자에는 sample/last-success 시간을 연결한다. 오류에서는 unavailable 또는 마지막 값+stale를 표시한다. last_tool_call_at은 최근 도구 요청의 수신 시각이며 완료 시각이나 summary 조회 시각이 아니다. 현재 실행이 모두 끝나 카운터가 0이어도 최근 요청이 있었다는 사실을 별도로 표시한다. unknown은 결과 상태 불명이지 미열람 개수가 아니다. 120/180/300초 계약과 root/child/diagnostic 필터는 바꾸지 않는다. 한국어 문구는 K branch가 담당한다.

회귀: healthy=true 유지 상태의 자동 0→1→0, 승인 0→1→0, unknown 증가, 빠른 완료 후 정상 0과 최근 호출 갱신, invalid 응답·오프라인·숨김·재표시·취소·응답 역전. 단발 수동 refresh 시험으로 대체하지 않는다.

### H. 로컬 건강 표시

`RuntimeService.ReadHealthAsync/GetSnapshotAsync`, `MainWindow.ApplySnapshot/ApplyLiveRuntimeStatus`, tray status가 대상이다. 로컬 health는 `/healthz`로만 결정하고 공개 터널 down은 별도 공개 상태다. 두 렌더링 경로를 같은 snapshot projection으로 묶으며 mode/URL/port의 시점을 일치시키고 편집 중 사용자 입력을 덮지 않는다. R의 새 역할 DTO 없이도 단독 branch에서 기존 모델로 동작해야 한다.

회귀: local healthy/public down, Core running/health failing, query unavailable, refresh 역전, 편집 입력 유지. native supervisor 식별은 R 시험에서 별도로 검증한다.

### A. 작업 디렉터리 ACL

`Config.Normalize`/경로 준비는 AgentDockHome에만 보호 처리를 적용한다. 사용자 작업 디렉터리와 기존 자식의 보호 상속 DACL을 변경하지 않는다. 경로 존재·디렉터리·정규화 검증은 유지한다. 기존 작업 디렉터리 ACL의 자동 복구 기능은 추가하지 않는다. Windows만의 요구를 근거로 다른 OS 보안을 일괄 완화하지 않는다.

회귀: 작업 디렉터리·기존 자식 DACL 전후 동일, 의도한 홈 보호 ACL 직접 확인, 사용자 경로·여러 SID·실패 시 보존. BASE의 경로 준비 함수에 작은 패치와 필요한 인접 시험을 적용한다.

### Q. rg 내장

`internal/bundledrg`와 `search_rg_selection`을 BASE의 검색·패키징 경계에 추가한다. 검증 가능한 함수·공식 pin·라이선스는 현재 계약에 맞는 범위만 재사용하며 이 기능 브랜치에 downstream release/version assertion은 넣지 않는다. 실제 실행 Core의 옆 `tools/rg`를 선택한다. mutable active pointer·cwd·개발 cache를 번들 선택 근거로 쓰지 않는다.

공식 고정 artifact·크기/SHA·아키텍처·라이선스를 검증하고 실행 종료까지 Windows 교체 방지 읽기 handle을 유지한다. bundle 완전 부재만 허용 PATH→기존 Go fallback이다. present-but-incomplete/corrupt/redirected bundle은 fail-closed다. `--no-config`, `-e <query> -- <path>`, no-match exit 1 정상 처리, regex/실행 실패를 다른 엔진 성공으로 바꾸지 않는 계약을 유지한다.

package·generation publish·repair·rollback과 **1.1.6 legacy migration**의 경로에 같은 검증 component를 연결한다. migration이 Core와 필요한 component를 함께 처리하는지 확인한다. 지원할 수 없는 payload는 변경 전에 명확히 거절하고 Setup을 안내하되, 지원 가능한 migration 경로는 유지한다. 다른 OS·Windows ARM64의 기존 검색은 보존한다. 운영/user/session/WSL PATH는 변경하지 않는다.

회귀: bundle 유효/부재/부분/변조/취소·PATH 오염·Unicode/대시 query·실제 rg binary 검색·migration/repair/rollback의 component와 hash 일치. 최종 필수 실제 binary 시험은 fixture 부재를 skip하지 않는다.

### C. 저장 상태·discovery·receipt

G0-4에서 이미 해결됐으면 구현/PR을 생략한다. 남았다면 기존 소유자 함수에 다음 최소 계약을 적용한다.

state 보존: 파일 부재에만 초기값을 만든다. 기존 null·필드 누락·미래 schema·trailing data·과대 envelope는 원본 불변 거절. rejected update는 마지막 정상 연결/리비전·notification claims·insertion 소비 상태를 바꾸지 않는다. 하나의 범용 새 저장 엔진을 만들지 않는다.

discovery: optional SDK params의 interface 안 typed-nil과 nil ListToolsResult를 정확히 처리한다. HTTP/SDK/직접 호출 및 실제 optional metadata를 시험한다. 범용 recover로 panic을 삼키지 않으며 BASE의 catalog/response_additions/presentation 계약을 유지한다.

receipt: `ERROR_SHARING_VIOLATION`/`ERROR_LOCK_VIOLATION` 및 미완료 파일 읽기는 기존 deadline 안의 pending이다. launch nonce 불일치·잘못된 JSON·기타 I/O는 성공으로 바꾸지 않는다. 일시적 읽기 경합 해결을 설치 실패 warning 성공과 혼동하지 않는다.

### D. downstream 전용

온라인 업데이트 정책·자체 URL·Distribution·현재 계획은 integration 전용이다. online check는 최신 버전을 안다고 주장하지 않고 offline-manual/online-updates-disabled를 반환한다. online apply는 네트워크/파일 변경 전에 실패한다. upstream 패키지를 fallback으로 내려받지 않는다.

1.1.6의 local-archive/generation/legacy migration API 자체를 온라인 정책 때문에 삭제하지 않는다. remote 다운로드를 동반하는 entry와 이미 제공된 local payload 처리의 책임을 구분해 offline 정책을 한 곳에 둔다. CI 자동 게시/태그/원격 push는 이번 구현에서 실행하지 않는다.

## 8. 신속 구현 순서와 최소 검증

### 8.1 한 번의 구현 흐름

1. 시작 때 BASE·dirty 상태·필수 도구 입력을 묶어 확인한다. 이미 확인한 사실을 매 branch마다 다시 수집하지 않는다. G0는 해당 단위 직전에만 닫고 upstream이 이미 해결한 항목은 즉시 제외한다.
2. BASE에서 독립 기능 branch를 만든다. 원본 main을 reset하거나 일괄 cherry-pick하지 않는다. `D:\Engineering\agentdock-eerraa\dist\worktrees\implementation-1.1.6100`의 단일 구현 worktree를 재사용할 수 있다. branch/worktree/작업 관리 자료를 필요 이상으로 늘리지 않는다.
3. R은 controller/host→Quick→installer/rollback→역할/미사용 입구 순으로 연속 구현한다. K/S/H/A/Q/C는 실제 의존성만 고려해 진행한다. 하나가 막히면 독립 단위로 전환하며 승인된 단계 사이에 재승인을 기다리지 않는다. 한글화의 모든 코드·리소스·시험은 K branch에 둔다.
4. 각 기능이 의미 있는 단위로 완성되면 변경 경계의 compile/정적 검사와 최소 회귀를 묶어 한 번 수행하고 로컬 commit한다. 중간 compile은 실제 컴파일 위험을 줄이는 데 필요할 때만 실행한다. 실패 시 그 실패와 영향 범위만 재검증한다.
5. 완료한 branch를 integration에 정상 merge한다. 단독 branch 검증 결과는 재사용하고, 통합 시에는 merge 충돌·공유 호출 경로·기능 조합을 한 번 확인한다. **전체 suite는 기본 필수 관문이 아니다.** 공통 admission·직렬화·저장 형식·응답 조립처럼 실제 영향이 넓거나 표적 실패가 파급을 보여줄 때만 관련 suite로 넓힌다. 그때도 가능한 최소 범위로 한 번 실행하며 전체 OS 행렬로 확대하지 않는다.
6. main의 마지막 downstream release commit에서만 1.1.16101/Revision101/UpstreamVersion1.1.6을 적용한다. 버전만 바뀌면 이전 기능 회귀를 반복하지 않고 version/resource/source 정합성을 확인한다. clean commit에서 amd64 최종 package를 한 번 만들고 산출물을 검사한다. 코드/패키징 결함으로 결과가 무효가 된 경우에만 해당 검증과 package를 다시 실행한다.

### 8.2 실제 실행할 최소 검사

§5~§7의 회귀 목록은 지켜야 할 동작·위험 목록이다. 아래처럼 기존 시험을 묶고 대표 경계로 줄여 검증한다. 이미 같은 계약을 검사하는 fixture가 있으면 새 fixture를 만들지 않는다. 가벼운 테이블 기반 경우들은 한 실행으로 처리한다.

| 변경 단위 | 필수 최소 확인 | 기본적으로 제외할 반복/확장 |
|---|---|---|
| task controller·host/Job | 정상 End→종료→Run, 다른 owner 거절, 자식 실행 전 Job 귀속과 host 종료 정리, 부분 생성 실패의 대표 사례를 기존 mock/helper로 검사 | 실제 운영 task 제어, Windows 버전별 반복, 같은 helper를 branch/commit마다 재구축 |
| Quick·상승 configure | 건강 미준비에서 spawn/crashRetry=0, 회복 뒤 생성, Origin 적용 때 Core만 교체·자기 task End=0·ready 순서, configure 포트 비대기를 fake clock/child로 검사 | 실제 공개 터널/네트워크 대기, 120초 실시간 대기, 같은 수명 경계의 중복 시험 |
| 설치·복원/표준 실패 | 변경된 정상 경로, owner 없는 성공 거절, 표준 start→health와 실패 전파, 파일/포인터·작업 소유권의 서로 다른 실패 책임을 대표하는 rollback 사례·resume≤1·재대기=0 | 신규/upgrade/repair × 모든 실패 지점의 전수 조합. 경로가 실제로 다를 때만 해당 사례 추가 |
| K 한글화 | 전체 리소스 키/placeholder 자동 검사, 변경된 producer→저장/재열기→새 UI 표시와 사용자/실제 출력 불변의 대표 loaded fixture | 화면×테마×DPI×언어 전체 조합, 번역만을 위한 새 브라우저/설치 시험 환경, 변경 없는 화면의 반복 렌더링 |
| S/H/A | 자동 카운터 0→1→0·invalid/stale, local healthy/public down, 작업/자식 DACL 불변과 홈 보호를 각 기존 표적 fixture에서 검사 | 전체 활동 엔진·모든 창·다중 PC 검증, 관련 없는 수명 회귀 반복 |
| Q rg·패키징 | 실제 고정 rg 기본 검색, bundle 누락/변조 구분, 바뀐 generation/migration/rollback의 component 보존, 최종 package의 pin/라이선스/구조 확인 | 전역 PATH 조작, 모든 OS/architecture 실기기 시험, 같은 bundle 재다운로드·반복 hash 캠페인 |
| C 잔여 소규모 패치 | 1.1.6에도 남은 각 결함을 직접 잡는 기존/추가 회귀와 영향 package 검사 | 이미 해결된 결함의 새 PR, 대형 공통 저장 엔진/시험 harness, 무관한 suite |
| 1.1.6 기존 기능 | 수정이 닿은 공통 호출 경계의 기존 인접 회귀만 선택 | 변경 없는 cache·Skill·catalog·payload·편집 엔진 전체를 새 기준이라는 이유로 재검증 |
| 최종 조합·산출물 | 공유 경로/충돌 부분 smoke, 전체 source와 버전 일치, amd64/unsigned·cloudflared 유효 서명·rg/core skills/체크섬 확인 | branch별 전체 Setup 생성, 별도 후보 패키지 누적, 실기기·VM 설치 및 자동 운영 교체 |

표의 핵심 계약을 확인하는 데 실제 OS 동작이 필요하면 작은 격리 helper 시험은 남긴다. 코드 문자열의 존재만으로 Job 소유권·복원 성공을 판정하지 않는다. 안전/권한/데이터 무결성 결함은 속도를 이유로 무시하지 않는다. timeout/cancel/입력 반례는 해당 로직을 수정한 경우 같은 fixture의 작은 사례로 추가한다.

### 8.3 반복 실행·대기 제한

검사 결과는 관련 소스·의존성·입력·설정이 그대로면 재사용한다. 전체 HEAD나 문서/버전만 바뀌었다는 이유로 무관한 결과를 폐기하지 않는다. 컴파일·WPF 검사는 기존 증분 build와 fixture를 사용하고 같은 소스를 연속 build/publish하지 않는다. 필요한 빌드/테스트 실행 하나가 여러 계약을 입증하면 별도 wrapper로 다시 실행하지 않는다.

전체 테스트를 먼저 돌려 실패 목록을 모으는 대신 수정 함수·직접 caller의 표적 검사부터 실행한다. 결함 재현은 가능하면 작은 실패 assertion 하나로 잡고 패치 뒤 같은 검사를 통과시킨다. 실제 설치 재현이나 전체 baseline suite 실행을 구현의 선행 조건으로 삼지 않는다. 새 회귀 harness의 비용이 패치보다 커지면 기존 fixture에 좁은 사례를 추가한다.

도구/의존성 실패는 실행 결과와 exit code를 먼저 확인하고, 기존 설치·cache를 한 번 확인한다. 해결 불가한 항목은 짧게 보류 사유를 남기고 독립 구현을 계속한다. 필요 없는 VM·SDK·도구 버전의 반복 설치/다운로드로 우회하지 않는다. 범위 밖 기존 실패나 비필수 검사 부재는 승인된 구현을 전부 멈추는 사유로 삼지 않되, 최종 산출물에 미치는 영향을 명시한다. 최종 필수 compile·소유권/복원·무결성 검사가 실패하면 완성/통과로 표시하지 않는다.

상세 결과 문서를 새로 만들지 않고 단위별로 실행 명령, pass/fail/not-run, 다음 행동만 남긴다. skip·assertion 약화·오류를 경고 성공으로 바꾸는 방식은 사용하지 않는다. 실기기/VM 설치는 사용자가 최종 파일로 진행하므로 미실행으로 남긴다.

## 9. 1.1.16101 릴리스 계약과 실행 명령

- Core·WPF·Setup 버전은 1.1.16101, DownstreamRevision=101, UpstreamVersion=1.1.6, Distribution=eerraa다. buildinfo와 release tool의 명시적 릴리스 식별자·버전 정렬·전체 source_commit을 확인한다.
- PowerShell 7과 clean 최종 main commit이 필수다. 원본 main의 문서를 삭제/stash해서 억지로 clean하게 만들지 않는다.
- 출력: `D:\Engineering\agentdock-eerraa\dist\windows-release-1.1.16101`.
- Windows payload/Setup은 amd64만. 기존 공통 WSL helper의 Linux payload 계약은 Windows 아키텍처 선택과 별개이며 임의 제거하지 않는다.
- cloudflared: `D:\Engineering\agentdock-eerraa\dist\windows-release\cloudflared.exe` 그대로 재사용. 빌드 직전 Authenticode Valid·예상 Cloudflare 서명자를 확인하며 자동 latest 다운로드로 대체하지 않는다.
- Inno: `D:\Engineering\.agentdock-build-tools\inno\ISCC.exe`.
- AgentDock payload는 unsigned. `-SignedBuild`/서명 키·신뢰 설정 변경을 사용하지 않는다. cloudflared의 원래 유효 서명은 유지한다.
- `scripts/install/install.ps1`는 ASCII다. UTF-8 오류 전달 개선과 script source 인코딩을 혼동하지 않는다.

G0-5에서 실제 script param을 확인하고, 아래 지원 계약을 만족하는 package entry를 구현·검증한 **뒤에만** 실행한다.

```powershell
$ErrorActionPreference = 'Stop'
$repo = 'D:\Engineering\agentdock-eerraa'
$out = 'D:\Engineering\agentdock-eerraa\dist\windows-release-1.1.16101'
$cf = 'D:\Engineering\agentdock-eerraa\dist\windows-release\cloudflared.exe'
$inno = 'D:\Engineering\.agentdock-build-tools\inno\ISCC.exe'
if ($PSVersionTable.PSVersion.Major -lt 7) { throw 'PowerShell 7 required.' }
$changes = @(git -C $repo status --porcelain)
if ($LASTEXITCODE -ne 0) { throw 'Git status failed.' }
if ($changes.Count -ne 0) { throw 'Release tree is not clean.' }
if (-not (Test-Path -LiteralPath $cf -PathType Leaf)) { throw 'Pinned cloudflared input is missing.' }
$sig = Get-AuthenticodeSignature -LiteralPath $cf
if ($sig.Status -ne [Management.Automation.SignatureStatus]::Valid) { throw 'Invalid cloudflared signature.' }
if (-not (Test-Path -LiteralPath $inno -PathType Leaf)) { throw 'ISCC is missing.' }
& "$repo\packaging\windows\build-windows-release.ps1" `
    -Architectures amd64 -OutputDirectory $out `
    -CloudflaredBinary $cf -InnoCompiler $inno
if ($LASTEXITCODE -ne 0) { throw 'Release build failed.' }
```

script 안의 각 native tool exit code도 검사해야 한다. 기존 script가 release 하위에 Setup을 쓰면 실제 경로를 보고한다. 파일명을 추정해 존재하지 않는 Setup을 전달하지 않는다. 최종 산출물은 사용할 Setup.exe 하나이며 package 구조·rg·core skills·checksums·Core/WPF/Setup version·source SHA·unsigned/서명 경계를 검증한다.

최종 보고는 BASE SHA, 각 branch HEAD와 단독 시험, integration source SHA, 실제 Setup 경로/크기/SHA-256, 사용자 설치 시험 미실행을 포함한다. 최종 전달에 필요하지 않은 중간 산출물은 나열하지 않는다. 패키지 생성이 설치/rollback 검증을 통과했다는 뜻은 아니다. 최종 source 변경 뒤 dirty 재빌드는 하지 않는다.

## 10. 단일 정본과 문서 소유권

**개발 실행 계획의 정본은 `docs\eerraa\implementation-plan-1.1.6100.ko.md` 하나다.** 이 파일이 1.1.6 기반·1.1.16101 목표·병목 최소화 원칙·기능별 구현·최소 검증·패키징·다음 세션 지시를 모두 소유한다. 변경은 같은 파일에 반영하며 별도 실행안이나 중복 정본을 만들지 않는다. 첨부용 사본도 같은 파일명과 내용으로 제공한다.

본문에는 현재 적용할 기준·결정·필요한 착수 조건만 작성한다. 변경 이력·폐기안 비교·과거 세션 지시·삭제 기록은 넣지 않는다. 실제 제품의 데이터 호환·복원 요구는 현재 기능 계약으로 기술한다.

자체 문서는 `docs/eerraa`에 두고 upstream 문서는 원래 위치에 유지한다. `KOREAN-NOTICE.md`는 설치기 번역·라이선스의 출처 설명이며 개발 계획서가 아니다. 라이선스 원문·번역자 크레딧·pin·패키징 리소스와 제품이 소비하는 core-skills·schema·fixture를 보존한다.

구현에 필요한 개발 문서 포인터는 downstream 문서 커밋에서 이 정본 경로를 가리키도록 한다. 실행 권한과 세션 경계는 §1에 모으고 다른 개발 문서에 중복 작성하지 않는다. 문서 정리 때문에 main·index를 reset하거나 제품 코드를 바꾸지 않는다.

## 11. 남길 동작·없앨 완화책·합칠 중복·제거할 입구

| 종류 | 최종 자기검토 체크 |
|---|---|
| 남김 | 1.1.6 snapshot/workspace-context/response additions/catalog/payload/편집 통계/Skill lease, 기존 generation journal·DPAPI, 표준 start→health, 단일 task/Job, 한국어·rg·홈 보호 |
| 없앰 | 상승 소유자 없는 설치 성공, 표준 deferred 건강 성공, runtime runas 재시도, 상승 detached 터널·HKCU 중복, Quick 자기 task 종료, local-health 실패의 crash backoff, rollback 두 번째 장기 대기, summary를 실제 출력처럼 표시 |
| 합침 | task action 정책·operation gate, native host, effective tunnel mode, Origin writer, 역할 조회, 제품 생성 문구 renderer, summary parser·갱신 수명 |
| caller 증명 후 제거 | 경쟁 WPF/native task host, 사용하지 않는 별도 --task-core-host alias와 중복 WPF 연결(안정 tray의 --run-core-task 계약은 유지), C# 중복 task/process helper, 상승 전용 detached wrapper, 사용되지 않는 trial/COM wrapper. 유효 표준/local-archive caller는 남김 |
| 확인 전 금지 | 미확인 G0 항목을 해결됐다고 기록하거나 그 입구부터 삭제. 실행하지 않은 시험을 pass로 표시. 확인하지 않은 공개 Release 상태를 확정적으로 표시 |

OS 설계 참고: Microsoft Learn `UpdateProcThreadAttribute`의 JOB_LIST/HANDLE_LIST 및 redirected child I/O 계약을 따른다. 새 helper가 이미 이 저장소에 있다는 뜻은 아니다. 생성 시 Job 귀속과 허용 handle만 상속하는 동작은 별도 OS helper fixture로 입증해야 한다.
`https://learn.microsoft.com/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribute`

## 12. 후속 작업 지시

정본 경로는 `docs\eerraa\implementation-plan-1.1.6100.ko.md`로 유지한다. 고정 BASE와 독립 기능 브랜치의 기존 검증을 재사용한다. 기능 변경은 기존 소유 브랜치에서 수행하며 그 제품 버전은 upstream 1.1.6이다. 검증된 결과를 main에 정상 merge한 뒤 main에서만 명시적 제품 버전 1.1.16101, DownstreamRevision 101, UpstreamVersion 1.1.6을 적용한다. 버전 산식으로 다른 값으로 고치지 않는다. clean main의 지정 amd64 오프라인 Setup을 생성하고 파일·버전·source SHA·무결성을 확인한다. 설치 실행·운영 교체·push·PR·릴리스 게시는 하지 않는다.

## 현재 구현 상태 및 재개 경계

제품 기준과 BASE는 `1.1.6` / `df7c22f64438ec317e0518eec44d035ad98be2b2`로 유지한다. 시작 시 수행한 원격 ref 대조를 재사용한다.

| 단위 | 독립 branch / SHA | 구현과 실제 검증 |
|---|---|---|
| R | `fix/windows-elevated-runtime-lifetime` / `e0ab961dc6a634ecbe3e7bca7397646e7361aec0` | create-time Job/host 강제 종료·자식 정리, private channel, Task 소유권, 같은 사용자 kernel ACL, 설치 파일우선 복원 표적 통과. Task XML 12개·설치 AST mock 20개 및 WPF compile 통과. |
| K | `feat/korean-localization` / `2cabf2f21021f9b3c0729d403a2be6a591a8fb01` | 실제 제품 어셈블리 3언어 리소스·표시 6,096개, 기존 pure model 674개, producer 원문 보존·backend·공식 Inno 한국어/라이선스 핀 검사 통과. |
| G | `feat/windows-bundled-ripgrep` / `9f3fba36d08f1c5465171f09498fd1c6c2d55335` | 실제 rg 고정 번들·검색·generation 검사와 캐시-only 8개 결과 재사용. legacy source 자신의 sidecar 보존/변조 거절 추가 회귀 통과. 최초 외부 명령 세션 timeout과 네 패키지의 ok 출력은 구분한다. |
| A | `fix/workspace-acl-boundary` / `adf8acceaeedadd1b183844528e29458bd5f6aa4` | 실제 Windows 임시 workspace/root/child/file DACL 불변, AgentDockHome의 별도 private ACL 및 경로 검증 통과. |
| O | `feat/windows-loopback-health` / `d46a4fcc905bbd2a407cef3ecdbecd2cfcf4c7c1` | 기존 loopback 건강/cache 89개 및 WPF 검증 결과 재사용. 설치 버전 읽기에서 실행파일을 호출하지 않는다. |
| S | `fix/activity-summary-refresh` / `1fc380e3e1f330e1a359424e433f2712d339cbe8` | 메인 표시 중 자동 0→1→0, pending/unknown, strict schema/null 거절·last-good stale·숨김/재표시·late result 34개 통과. |
| H | `fix/runtime-health-display` / `2620fdb8eb642e14c8f9b4bc34270012c285a6b3` | full/live/tray 공통 로컬 건강·공개 터널 분리, 이전 snapshot 거절·편집 초안 보존 19개 통과. |
| C-discovery | `fix/mcp-discovery-null-boundary` / `02d60bf69ee2c0a9106f7dbadb9b50b65835345a` | 실제 tools:[null] SDK panic 재현 후 pre-SDK 검증과 last-good catalog 보존 회귀 통과. 정상 optional metadata 및 원본 오류 보존. |
| C-receipt | `fix/setup-receipt-retry` / `4c87488034d794ba4c2d125acd41a19e81b2663e` | 실제 Windows sharing lock·부분 JSON·nonce·원본 exit/error 보존 통과. 원래 deadline 안에서 일시적인 공유/잠금 오류만 대기. |
| D | `feat/eerraa-offline-distribution` / `b9a645c766b6d01d842854eeda62f1d2675e3449` | 온라인 check/update 0 HTTP 요청·upstream fallback 차단, own distribution/source SHA metadata 검사 및 compile 통과. |
| 추가 지시 미리보기 | `feat/activity-insert-summary` / `d04f1187b7d088927760baf3bf8479d41dedd6e5` | 160 rune 파생 미리보기와 원문 단일 저장·응답·중복 방지·만료 기존 결과 재사용. 메인 요약 S와 별개 기능이다. |

최종 기능 통합 소스는 `4d1e699982cd2948a94395afc0afe016c51f4865`이며 main은 이 제품 소스와 기존 자체 문서·게시 권한 경계를 함께 보존한다. 충돌이 있었던 R/O health 응답은 service·PID·origin hash를 함께 유지하고 해시값·원문 비노출·GET/HEAD 계약을 검사했다. K 리소스 병합은 기존 키와 자리표시자를 보존했고 S의 독립 observer를 유지했다. 통합 WPF는 경고 0/오류 0이며, 통합 S 34개·H 19개·설치 AST 20개 및 한국어 설치 계약이 통과했다. 동일 기능 전체 suite는 반복하지 않았다.

`fix/windows-named-state-security`는 BASE 그대로이며 별도 변경이나 PR 단위로 취급하지 않는다. 필요한 same-user runtime coordination 보안은 R이 소유한다. 과거 Q의 지정 5개 참조 객체는 로컬에 없어서 그 객체 자체를 검증했다고 주장하지 않는다. 실제 현재 소스에서 재현된 discovery/receipt와 메인 요약·상태 표시 문제는 위 기능 단위로 처리했다.

추가 state-reader 소스 대조와 묶음 조회 요청은 도구의 보안 상태 판정 단계에서 차단되었다. 해당 차단은 우회하지 않았으며 전수 state-reader audit는 미검증이다. 미확인 state 소유권/복원 입구는 수정하지 않고 BASE를 유지한다. 이는 확인된 결함이나 앞서 통과한 표적 시험의 실패를 뜻하지 않는다.

이번 검증은 임시 fixture/프로세스와 모의 OS adapter에 한정한다. 실제 설치된 elevated Task의 UAC·session 전환, 실기기/VM Setup 설치·업그레이드·복구, 실제 Cloudflare/Tailscale 인터넷 준비, 전수 UI/DPI/OS 행렬은 실행하지 않았다. 운영 교체·push·PR·릴리스 게시는 하지 않는다.

main의 최종 downstream release commit에서만 `Version=1.1.16101`, `DownstreamRevision=101`, `UpstreamVersion=1.1.6`을 적용한다. 버전·문서만 바뀌면 위 기능 결과를 재사용한다. 최종 clean 소스는 PowerShell 7에서 지정 cloudflared/ISCC를 사용하여 `D:\Engineering\agentdock-eerraa\dist\windows-release-1.1.16101`에 amd64 오프라인 Setup 하나를 만든다. 실제 생성 성공·파일 크기·SHA-256·서명은 생성 뒤 산출물과 기존 build-report로 확인하며, 이 소스 checkpoint 자체를 설치/실행 시험 완료 증거로 쓰지 않는다.

상태 후속 검증은 native 상태 파서/조회 수명 39개, 표시 판정 27개, 기존 활동 요약 34개가 통과했고 WPF compile은 경고 0/오류 0이다. 통합 어셈블리의 현재 설치본 읽기 전용 상태 조회에서 Core/health/cloudflared/ready=true, named, RunningNormally를 확인했다. 일반 권한 운영 UI의 실제 화면 교체 검증은 미실행이다.

`fix/core-skill-bundle-manifest`의 `0eda649fd8a51bdf3bb23aa6731afad5cd10770d`는 bundle manifest의 필수 Skill version 누락을 수정한다. 실제 bundle 생성→격리 bootstrap과 LF/CRLF 불변·잘못된 version 필드 거절 회귀가 통과했다. R의 후속 COM null BSTR 수정은 빈 작업 디렉터리를 숫자 0으로 오판하지 않으며, 소유 root 허용과 다른 root·정수·missing 거절 6사례 및 기존 순서/소유권 회귀를 통과했다. 원래 설치 후속 worktree의 미커밋 파일은 그대로 보존한다. 패키징은 이 두 설치 수정이 포함된 clean main만 사용한다.
