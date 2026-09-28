# AgentDock Workbench — 원본 기반 최소 제품 및 Windows 배포 정본

기록일: 2026-09-28 KST
작업: `tsk_c83ac7427fcaa1db`
상태: **Windows 후보 패키지·Setup 설치/제거·실제 다운로드 검증 완료 / 전체 한국어·지원 upgrade 미완료 / 정식 Release 미게시**

이 문서는 현재 AGENTS.md가 가리키는 단일 실행 정본이다. 파일명의 1.1.6100은 과거 경로이며 현재 배포 버전이 아니다. 과거 main의 통합 기록은 기존 Git 이력과 외부 archive에 보존되어 있다. 아래 결과는 이번 원본 기반 후보에 실제 수행한 검사만 나타내며 과거 통합본의 PASS를 합산하지 않는다.

## 1. 기준, 작업 위치, 승인

| 항목 | 실제 기준 |
|---|---|
| 고정 원본 | Workbench v1.1.7 / `b367eaab95202873fb213b8713440bf7822878c4` |
| 최종 정본 경로 | `D:\Engineering\agentdock-workbench` |
| 사용자 원격 | `eerraa/AgentDock-Workbench`; repository ID 1380606195 |
| fork parent | `A-m-o-r-F-a-t-i/AgentDock-Workbench` |
| 현재 제품 작업 | 정본 Git에 연결된 `D:\Engineering\worktrees\agentdock-wb-release` |
| 작업 브랜치 | `work/minimal-wb-release-20260928` |
| 최신 실제 후보 빌드 소스 | `58b6f1496872c37150be925d76458552965f2058`; CI 36338434227, 상세 결과는 4.3절. 이후 후보와 실패는 4.5–4.7절 |
| 기존 local main | `cea147b9a3dbeb6f370c2c5231ca7a50604d7bdd` — 제품 후보로 교체·병합하지 않음 |
| 기존 remote main | `113f709a841ea64f238c528b561ab5f502d95239` — push하지 않음 |
| 현재 후보 소스 버전 | **1.1.17100, 미배포**; 원본 기준 1.1.7과 구분하며 태그는 생성하지 않음 |
| 외부 근거 | `D:\Engineering\archives\agentdock\2026-09-28\minimal-release-c83ac742` |

직전 문서 커밋 8174fec는 제품 변경 없이 검증 기록을 보존했다. 이후 새 배포 설정·후보 버전 변경이 진행 중이며, HEAD와 CI 결과는 외부 실행 기록으로 구분한다. `RESUME_FINAL_STATE.json`은 직전 재개의 종료 관측이지 후속 변경의 완료 증명이 아니다.

사용자는 이번 요청에서 정상 main 통합·사용자 원격 push·새 버전/태그·CI dispatch·Windows x64 정식 Release 게시 및 실제 재다운로드 검증을 명시적으로 승인했다. 별도 게시 승인을 다시 요구할 사안이 아니다. 그러나 불완전한 제품을 게시하거나 검증 실패를 숨길 권한은 아니다.

운영 PC의 Setup·설치·업데이트·제거·복구, 생산 Core 교체·재시작, 운영 Cloudflare/Tailscale·자격 증명 변경은 계속 금지한다. 실제 설치 검증은 진짜 격리 Windows runner 또는 폐기 가능한 VM에서만 수행한다. 사용자 원격 외 push, force push, 기존 공개 태그/바이너리 덮어쓰기, 열린 PR 변경, 다른 작업의 삭제도 금지한다.

## 2. 패치 결정 및 구현 범위

### 2.1 확정 유지: rg와 한국어

rg는 기존 generation·검색·installer·selfupdate 소유자를 사용한다. `internal/bundledrg`, 검색 선택기 및 기존 payload 복사/ZIP/generation 경계에만 연결했다. Windows amd64용 15.2.0 pin, 공식 archive 크기·SHA-256, 실행파일 및 라이선스 pin, 읽기 handle의 변경 방지, 손상/부재 구분을 유지한다. arm64용 rg는 구현하지 않았다.

이전 `33f2aed3c86db00439761e85b6380dd99a33ef40`에서 공급 스크립트와 실제 ZIP 검증을 기존 Windows 빌드·Setup 입력·배포물 검증에 연결했다. 이번에는 두 스크립트의 governance 등록 누락을 보완했다. manifest 읽기 중 취소가 ErrIntegrity로 바뀌는 결정적 반례를 재현하고 취소 원인을 그대로 반환하는 3줄 수정과 회귀를 추가했다. 무결성 검사나 fallback 정책을 약화하지 않았다.

실제 rg 구성요소를 사용한 검색·전달 회귀와 cache-only ZIP 검증은 통과했다. **이후 CI에서 후보 Setup과 ZIP을 빌드하고 실제로 내려받아 rg·라이선스·manifest를 재검증했다(4.3절). 전체 한국어 및 지원 upgrade가 미완료이므로 정식 배포물은 아니다.**

한국어는 후보의 기존 324개 리소스 및 문화권 처리 자산을 보존했다. 이번 재개에서 전체 UI·언어 선택·제품 브라우저/MCP Apps까지 한글화하지 않았다. 사용자 입력·로그·외부 MCP 결과를 번역하거나 Activity 모델을 바꾸는 작업은 하지 않았다. 관련 필수 소스 조회가 실행 전에 차단되어 미검토 대상을 임의로 교체하지 않았다.

### 2.2 MCP discovery — 채택·검증

고정 SDK의 header-annotation 필터는 null Tool을 참조하므로 변환 함수에서만 검사하면 늦다. 원본의 격리 HTTP discovery에서 panic을 재현했다. `internal/mcp/client/protocol.go`의 sending middleware에서 잘못된 결과·배열·null Tool·inputSchema 객체 형태를 먼저 검사하고 MCP_INVALID_RESPONSE를 보존한다.

기존 Manager/Catalog의 last-good 책임, 정상 빈 목록, 선택적 metadata, 원래 transport 오류를 유지한다. 새 SDK fork나 권한 체계는 없다. `discovery_null_test.go` 및 `discovery_schema_test.go`가 null/typed nil, 배열·문자열·숫자·boolean inputSchema, 정상 객체, 다른 method의 원문 보존을 검사한다. SDK/상류가 같은 선행 검증과 회귀를 제공하면 이 작은 보완을 제거할 수 있다.

### 2.3 Windows 프로젝트 ACL — 채택·검증

원본 Config.Normalize가 사용자가 선택한 프로젝트와 하위 파일의 DACL을 실제로 바꾸는 반례를 임시 디렉터리에서 재현했다. `internal/config/config.go`의 호출 경계만 변경하여 Windows에서는 AgentDockHome만 EnsurePrivate 대상이다. 프로젝트 경로 생성·검증 및 Unix 동작은 유지한다.

`workspace_acl_windows_test.go`는 반복 Normalize 후 프로젝트·하위 폴더·파일의 DACL과 데이터 불변, 자체 Home 보호, 파일/상대 경로 거부를 검증한다. 실제 프로젝트의 ACL을 시험용으로 변경하지 않았다.

### 2.4 Core 역할·root 및 종료 handle — 채택·검증

원본이 동일 실행파일의 형제 controller·다른 root의 Core까지 종료하고, controller만 있어도 Running으로 표시하는 반례를 inert Windows 프로세스로 재현했다.

`service_process_role_windows.go`는 기존 Core 명령과 절대 runtime root를 현재 프로세스 handle에서 확인한다. command line 조회는 크기를 제한하며 상태 파일·cache·epoch·pipe를 만들지 않는다. 실제 종료 handle에서 이미지 경로·역할을 다시 검사하고 취소·조회 실패를 성공으로 바꾸지 않는다. `service_windows.go`의 Core 상태·정지 경계에만 연결했다. installer의 기존 경로 전체 종료 API와 Tunnel supervisor 구현은 그대로다.

실제 disposable 프로세스 시험은 controller/supervisor/다른 root 생존, 해당 Core만 종료, idempotent 정지, unhealthy Core 표시, 잘못된 경로·역할·handle과 취소 시 비파괴 거부를 확인했다. 접근 불가 경계의 오류 처리를 검증했으나 다른 계정의 실제 서비스에 종료를 시도한 시험은 아니다.

후속 CI 36364901575에서 이 선택이 열 수 없는 같은 이름의 프로세스(`ERROR_ACCESS_DENIED`)를 전체 정지 실패로 바꾸는 결함이 드러났다(4.7절). 열 수 없는 후보는 이 root의 Core임을 증명할 수 없으므로 종료된 PID와 같이 선택·종료하지 않는다. 이는 원본의 경로 기반 열거와 기존 "unknown, never stopped" 계약과 같다. 열린 handle의 경로·역할 조회 실패와 선택된 Core의 종료 handle 거부는 계속 오류다.

### 2.5 TaskAdmin — 목적별 축소·검증

현재 root·SID·task name·정의·실행파일·인자 및 복구자료 대상 검증을 유지했다. 복원 입력 검증은 기존 Task의 정지·삭제보다 먼저 실행한다. 새 Task 정의는 현행 stable tray action만 허용하고, 구형 stable Core action과 정확한 기존 PowerShell -File launcher는 기존 task/명시적 복구 입력에서만 allowLegacyAction으로 허용한다. 이 제품이 실행하지 않는 `--task-core-host` action은 제외했다.

보존된 1.1.16102의 실제 build-report와 Setup/ZIP/install.ps1 3개 hash가 일치함을 확인했다. 그 정확한 source `4d719ce75ac748cb24e8ad4b7d92ff7c56e24225`는 root/task/SID가 결합된 schema 1 JSON과 Unicode task XML을 기록했다. 따라서 schema 1을 무조건 삭제하지 않고 읽기 전용 호환 adapter로 한정했다. 새 기록은 원본의 schema 2를 사용한다. 원본의 schema 0 absent-task 읽기 호환은 유지하되 소유권 없는 기록은 변조 작업을 승인하지 못한다.

추가로 schema 1의 WasRunning을 원본 schema 2와 동일하게 처리해 복원 도중 Task를 즉시 시작하는 결함을 재현했다. 1.1.16102는 외부 제어자가 source 파일·manifest 복원 후 재개하도록 했으므로 schema 1은 그 순서를 보존하고 schema 2만 원본 즉시 재개 동작을 유지한다. 새 복구 coordinator를 추가하지 않았다.

실제 COM/NTFS fixture는 schema 1 정의와 명시적 deny DACL의 정확한 복원, disabled 상태, 원문 JSON/XML 바이트 불변 및 다른 root/task/SID 거부를 확인했다. 이 결과는 실제 제품 업데이트·Setup 검증을 대신하지 않는다. 원본 PrivilegeTransition, TaskSecurityDescriptor 및 그 원본 보안 시험 파일은 변경하지 않았다.

### 2.6 보류·제외

composite runtime host, Origin 교체 pipe, host/Tunnel epoch·state, 광범위한 health/version cache, Activity observer/insertion preview, 옛 downstream revision·온라인 업데이트 일괄 금지 framework는 채택하지 않았다. activity/insertion/permission/process/shim 및 MCP Manager/Catalog, 의존성 go.mod/go.sum의 원본 대비 diff가 0임을 확인했다.

Setup receipt 재시도, Named/Quick/Tailscale 표시 묶음, 생성 시점 Job 재작성, timeout 일괄 변경, cold-start/성능 보완, 보조 upstream 묶음은 자동 이식하지 않는다. 현재 원본의 구체적 결함이 입증될 때만 작은 변경으로 재판단한다. 열린 upstream PR 여부는 제품 채택 기준이 아니다.

## 3. 이번 소스 이력

| commit | 범위 |
|---|---|
| `33f2aed3c86db00439761e85b6380dd99a33ef40` | 이전 재개에서 rg 공급·ZIP·Setup 입력 검증 연결 |
| `0c4433e24cdd117ab1e3fa6d8475bb1a76ae7d67` | MCP discovery 및 Windows 프로젝트 ACL |
| `3d262e80975e26b53d242cb855b3394e645bc9b1` | Core 역할/root/동일 handle 확인 |
| `2bee6fdcef5ec44d2190a0dd70c4051cda4681c2` | 구형 복구 입력 한정과 schema별 재개 순서 |
| `b1d8589dce326763bf6f8447163b8ba21e60f53e` | rg manifest 취소, MCP 형태 회귀, TaskAdmin 정적 계약 및 스크립트 inventory |
| `bda69caaf04ad861f73c29449bbf440837a29d93` | 사용자 fork 업데이트·배포 주소, Task rollback runtime root와 모든 버전 formal source 검사 |
| `381b932e5729c2b03d3a87ffe1131cc2b059810e` | 1.1.17100 후보 및 원격/서명/필수 시험/불변 게시 검증 경로 |
| `4c4eed6e50363b14bdc09b4fb33a898f3d34d3aa` | 제품 tree를 바꾸지 않는 기존 main 이력 보존 merge |
| `58b6f1496872c37150be925d76458552965f2058` | 실제 Setup repair 실패의 구형 PowerShell action 호환과 실패 증거 보존 |
| `900723d6cf9096244e9b8c9b0eed633caf8abd3b` | 남은 한국어 표시와 Setup result.json 공유 위반 대기 |
| `f253e31` | 도구 생성 제목의 descriptor 결합, Setup E2E 세 언어 확인(이전 미커밋 3개 파일) |
| `361d426` | 접근 불가 같은 이름 프로세스로 Core 상태·정지가 실패하는 결함 |

직전 b1d8589 제품 소스의 WB_BASE 대비 차이는 43개 파일, 3,785줄 추가·53줄 삭제다. 여기에는 이전 후보의 한국어·rg 및 시험 코드가 포함된다. 이 숫자를 이번에 새로 완성한 기능 수로 해석하지 않는다. 이번 재개의 소스 커밋 4개만 외부 resume-source-patches에 format-patch 및 SHA-256으로 추가 보존했다.

## 4. 실제 검증과 실패 보존

이 절의 표는 직전 b1d8589 재개의 검증 기록이다. 후속 최신 후보/CI 결과는 4.3절에 별도로 기록한다. 수치를 서로 더하지 않는다. 전체 Go, 부분 회귀, native assertion, 시나리오 및 패키지는 다른 집계다.

| 검사 | 관측 결과와 근거 |
|---|---|
| MCP/ACL 원본 반례 | 2 tests FAIL, 2 packages FAIL. `minimal-baseline-01` |
| MCP/ACL 최초 수정 후 | 148 PASS / 2 SKIP / 0 FAIL, 3 packages PASS. `minimal-fixed-01` |
| Core 원본 반례 | 2 tests FAIL. `core-role-baseline-01` |
| Core 수정 후 전체 관련 패키지 | 200 PASS / 2 SKIP / 0 FAIL. `core-role-fixed-01` |
| 구형 복구 조기 재개 반례 | schema 1의 Run 호출을 기록하는 test double에서 FAIL; 실제 Task 시작 없음. `native-legacy-baseline` |
| manifest 취소 반례 | 1 FAIL; 수정 후 bundle 관련 suite 통과. `rg-cancel-baseline-01`, `rg-script-fixed-01` |
| 직전 b1d8589 전체 Go | `go test -json -p 1 -count=1 -timeout=180s ./...`: 2,135 PASS / 85 SKIP / 0 FAIL, 시험 패키지 56 PASS, 시험 없는 패키지 8. `resume-full-go-02` |
| go vet / go build | 각각 exit 0. `resume-vet-final`, `resume-build-final` |
| 핵심 race + 실제 rg | 명시한 12 packages, `-race -tags=bundled_rg_integration -p 1 -count=1 -timeout=600s`: 912 PASS / 48 SKIP / 0 FAIL. `resume-critical-race-final` |
| rg 일반 연결 | 실제 번들로 4 packages 361 PASS / 42 SKIP / 0 FAIL. `rg-connected-resume-01`; 이후 manifest 수정은 최신 tagged race에도 포함 |
| cache-only rg 패키징 | directory 8 / ZIP 9 / 연결 3 PASS, SKIP 0. `rg-package-resume-final.json`; 새 제품 Setup 생성 검사는 아님 |
| Windows 제품/정책/native/layout build | Release win-x64, 각 최종 build 경고 0·오류 0. `native-final-build`, `policy-final-build`, `layout-final-build` |
| Windows 순수 정책 | 890 assertions PASS. `policy-final-test` |
| Task 소유권 / 보안·복구 pure contract | 40 / 99 assertions PASS. `native-final-contract` |
| 실제 native COM/NTFS | 22 scenarios / 272 assertions / 0 failures. `native-final-recovery/native-privilege-validation.json` |
| source identity | native 검증 전후 1,337 파일 불변. Go/race/Windows 마지막 검증도 snapshot 해시와 status 불변. 후속 제품 C# blob은 native 검증한 2bee6fd와 동일 |
| WPF 실제 렌더링 | NOT RUN. 원본 GitHub Actions 전용 guard 유지; 환경 위조·guard 삭제 없음 |
| 최종 버전 package / 실제 Setup / 업데이트 / 게시·재다운로드 | NOT RUN / NOT RELEASED |

과거 두 180초 timeout의 원시 stack은 Config.Normalize → EnsurePrivate → SetNamedSecurityInfo에서 대기했다. 이번에는 HOME/USERPROFILE/APPDATA/LOCALAPPDATA/TEMP/AgentDock Home·workspace를 시험 자식 프로세스에만 분리했고 제품의 프로젝트 ACL 수정과 함께 동일 180초 전체 관문을 통과했다. 당시 지연의 모든 OS 내부 원인을 확정했다는 의미는 아니다.

이번 첫 전체 실행 `resume-full-go-01`의 2,120 PASS / 4 FAIL / 85 SKIP도 보존한다. Plugin 2건은 외부 archive 아래 지나치게 깊게 만든 TEMP에서 Git의 Filename too long이 발생했다. 제품 Plugin 코드·시험·hook 정책·전역 Git 설정을 그대로 두고, 임시 경로만 짧은 고유 `D:\Engineering\tmp\wb-*`로 바꿔 동일 2개 시험과 전체 suite를 통과했다. 나머지는 새 schema별 재개 동작을 반영하지 못한 정적 source fragment와 rg 공급 스크립트 inventory 누락으로, 실제 동작 시험을 유지하면서 수정했다.

첫 native contract의 schema 0 오류 메시지 회귀는 원본 오류 계약을 복원해 해결했다. 실패한 로그를 삭제하거나 assertion을 약화하지 않았다. 인코딩·줄바꿈·PowerShell의 patch 출력 인자 오류는 일반 실행/검증 오류로 처리했으며 보안 차단이라고 바꾸어 기록하지 않는다.

race compiler는 기존 검증 cache의 x64 GCC를 자식 CC/CGO/PATH에만 지정했다. 시스템 PATH나 다른 프로젝트 compiler를 변경하지 않았다. 600초는 race 계측의 외부 시험 예산이며 제품 timeout이나 내부 assertion은 변경하지 않았다.

## 4.1 두 번째 재개 — 배포 경계와 후보 버전

사용자의 계속 요청으로 bda69ca에 배포 주소 및 복구 호출 수정을 커밋했다. 수정 전 fork URL/복구 root/formal source 검사 6개 실패를 기록했고, 수정 후 scripts/test와 selfupdate 2개 패키지 130 PASS / 0 SKIP / 0 FAIL을 확인했다. 실제 installer에서 추출한 복구 호출과 순수 함수는 가짜 Start-Process로 실행해 원래 runtime root/SID/backup 인자와 URL 선택을 검사했다. 운영 설치나 UAC를 실행한 시험이 아니다.

이후 기존 windows-package.yml 및 workbench-acceptance.yml에 사용자 fork/x64, Linux와 Windows 동일 SHA, 실제 native 복구·WPF·Setup·역사적 1.1.16102 업그레이드/rollback 관문을 연결했다. 한국어 전체 검증이 없는 상태로 publish할 수 없으며, 후보 빌드는 그 미완료 상태를 not_run으로 기록한다. 과거 버전 baseline은 변경하지 않은 정확한 역사적 commit을 CI에서 다시 빌드하는 방식이고, 원래 게시 bytes와 동일하다고 주장하지 않는다.

정식 게시 전 모든 필수 시험과 checksummed metadata를 다시 검사한다. 이미 공개된 Release는 변경하지 않고, 기존 draft에서도 일치하는 bytes만 재사용한다. 실패한 조회를 부재로 해석하지 않으며 업로드 중 파일을 덮어쓰지 않는다. draft 및 공개 직후 각각 새 디렉터리에 실제 재다운로드하여 크기·SHA·구성·버전·source commit을 확인하는 경로를 기존 workflow에 구현했다. 현재 이 경로의 GitHub 동작은 시험 더블로만 검사했다.

publication policy를 포함한 scripts/test 108 PASS / 0 SKIP / 0 FAIL 및 브랜드·버전 일치 검사를 통과했다. 정상 신규 게시, 부분 draft 재개와 11개 실패 조건을 실제 workflow 본문에 대해 시험 더블로 실행했다. 시험 더블의 PowerShell scope 및 배열 인자 전달 오류도 실패 로그를 보존하고 수정했다. 필수 시험 누락/다른 SHA/다른 원격/변조 bytes는 승인하지 않는다. 이 기록은 실제 Release 게시 성공을 뜻하지 않는다.

새 미배포 버전은 1.1.17100이다. GUI informational version에 source SHA를 넣고 Core 외 shim/arbiter도 같은 SHA인지 검증하도록 했다. 최종 패키지·실제 CI 결과는 아직 이 절의 로컬 회귀 수치에 포함하지 않는다.

두 번째 재개의 한국어 3-way 후보는 제품 밖 korean-port에 준비했으나 충돌 부분 읽기 요청이 동일한 명시적 보안 문구로 실행 전에 차단됐다. resume2-localization-conflict-block.json에 기록했으며 우회 재조회·제품 적용하지 않았다. 따라서 전체 한국어 관문은 계속 미완료다.

## 4.2 실제 후보 CI와 구형 Task repair 보완

후보 source 381b932의 tree f1c8f66a를 그대로 유지하는 이력 보존 merge 4c4eed6e50363b14bdc09b4fb33a898f3d34d3aa를 만들고 사용자 원격의 작업 브랜치에만 push했다. 기존 local main cea147b9와 remote main113f709a는 변경하지 않았다. 기존 main과 1.1.16102 source가 조상으로 보존되지만 그 옛 제품 tree를 다시 이식하지 않았다.

실제 GitHub 실행36337080381은 정확한4c SHA에 대해 amd64, installation_tests=true, publish=false로 실행했다. Linux backend/race 및 Windows native 검증이 통과했다. artifact를 실제 내려받아 같은 SHA를 확인했다. native Go221 PASS/36 SKIP/0 FAIL, COM/NTFS22시나리오272assertions/0실패, WPF60개100/125/150/200% offscreen표본27538assertions가 기록됐다. 물리 키보드·모니터 검증 및 한국어 전체 검증은 아니다.

같은4c SHA의 로컬 전체Go는2157 PASS/85 SKIP/0 FAIL, vet/build는exit0, 핵심12패키지race+실제rg는912 PASS/48 SKIP/0 FAIL이다. 실제GUI빌드의ProductVersion은1.1.17100+전체4cSHA이며NotSigned이다. 운영UI/Setup은 실행하지 않았다.

이 CI는 Windows ZIP/오프라인Setup 빌드와 실제패키지metadata/checksum 검증까지 통과했지만 **Setup repair에서 실패**했다. 처음 설치·반복 설치가 완료된 뒤, 원본 E2E가 만드는 기존 PowerShell 예약작업의 소유권 검사에서 거부됐다. 역사적1.1.16102 upgrade/rollback은 이 실패 때문에 아직 도달하지 않았다. 이 실패를 통과로 바꾸거나 테스트를 제거하지 않는다. 첫 실행에는 실패시패키지보존단계가 없어 최종Setup/ZIP artifact를 회수하지 못했으며, CI로그와별도native/Linux artifact만 보존했다.

동일한 XML/인자로 로컬 순수정책 반례를 재현했다. 기존 root/SID/단일action/InteractiveToken 검사를 유지하면서 정확한 powershell.exe 또는 시스템WindowsPowerShell 경로와 기존 start-agentdock.ps1 -File 인자만 기존task/복구 경계에서 식별하도록 보완했다. 신규task는 여전히stable native tray만 사용한다. 임의상대실행파일·다른root·다른SID·추가명령·Command/EncodedCommand·다중action은 거부하며 PowerShell을 새로 시작하는제품경로는 추가하지 않았다. 수정후구형action32assertions와기존소유권40/복구99assertions가 통과했다. 실제Setup재검사는 후속CI에서 별도로 확인해야 한다.

후속실패원인을보존하도록 원래E2E에 known4개설치로그의최대2000줄씩선택적보존을연결했고, 실패후에도fixture정리를계속한다. workflow의실패산출물은unverified-windows-failure로명시하여정식게시관문과구분한다. 재검증이성공하기전설치완료·지원upgrade완료로보고하지않는다.

## 4.3 최신 실제 검증 상태 — source 58b6f149 / CI 36338434227

수정 커밋 `58b6f1496872c37150be925d76458552965f2058`을 사용자 원격 작업 브랜치에 정상 push했다. 같은 SHA의 실제 GitHub runner에서 `publish=false`, `installation_tests=true`, `architectures=amd64`로 다시 실행했다. 이전 실행 36337080381의 실패는 그대로 보존한다.

| 완료 상태 | 실제 결과 |
|---|---|
| Linux backend / race | 동일 SHA에서 모두 PASS |
| Windows 전체 Go | 2,172 PASS / 76 SKIP / 0 FAIL; 시험 패키지 56 PASS, 시험 없는 패키지 8 |
| 정적 분석 / 순수 정책 | 두 단계 모두 PASS |
| Windows native COM/NTFS | 22개 시나리오, 272 assertions, 실패 0 |
| WPF offscreen | 60개 렌더링 표본, 27,538 assertions; 100/125/150/200% 배율 |
| 실제 ZIP·오프라인 Setup | 빌드 및 패키지 metadata/checksum 검사 PASS |
| 실제 Setup E2E | 설치·반복 설치·구형 PowerShell 작업 repair·제거 PASS |
| 역사적 1.1.16102 upgrade/repair/rollback | **FAIL — 과거 버전 baseline 설치 준비에서 실패. 새 버전으로의 upgrade에는 도달하지 않음** |
| 전체 한국어 | 미완료; 정식 게시를 요구하지 않은 후보 실행이므로 필수 한국어 게시 단계는 SKIP |
| 전체 workflow | **failure**; 검증 완료 패키지 업로드/정식 게시 단계는 실행되지 않음 |

Windows Go 수치는 실제 다운로드한 `backend-validation.jsonl`에서 집계했다. 로컬 4c4eed6e의 2,157 PASS / 85 SKIP와 다른 환경의 결과이므로 합산하지 않는다. 물리 키보드/모니터 검증은 하지 않았고 offscreen 결과로 대신하지 않는다.

과거 버전 baseline은 원래 1.1.16102 source `4d719ce75ac748cb24e8ad4b7d92ff7c56e24225`를 수정하지 않고 다시 빌드했다. 그 빌드는 성공했지만 빈 경로 설치에서 stable shim이 `active-version.json`을 찾지 못했고, 이후 rollback도 실패했다. 로그의 관측 오류는 `read AgentDock active version ... active-version.json: The system cannot find the file specified`이다. 새 후보로의 upgrade가 실패했다고 바꾸어 기록하지 않는다. 정상적인 구형 설치 fixture를 마련하는 경로와 실제 지원 upgrade/rollback은 아직 검증되지 않았다.

실패 패키지는 별도 artifact `unverified-windows-failure-36338434227`로 보존했다. 이를 호스트의 외부 archive에 실제로 내려받아 기존 패키지 검증기를 다시 실행했고 다음을 확인했다.

- 버전 1.1.17100 / source 58b6f149... / 채널 candidate-not-released / Windows amd64가 일치한다.
- Core·arbiter·shim의 Go source 정보와 UI informational version의 source SHA가 일치한다. Setup 제품 버전도 일치한다.
- rg 15.2.0의 manifest, 실행파일 및 라이선스 5개 파일과 원본 크기/SHA 검증을 통과했다. Core Skill 3개는 별도 임시 Home에서 fresh/repeat bootstrap을 통과했다.
- Setup의 실제 Authenticode 상태는 NotSigned이다. 빌드에서 official cloudflared의 유효 서명을 검사했다. 운영 Setup이나 생산 Core를 시작하지 않았다.

| 후보 파일 | bytes | SHA-256 |
|---|---:|---|
| AgentDockSetup-amd64.exe | 110,869,060 | `c251e134a65f912362d27573bcc582237416071f47fd5f4438402101ee61632d` |
| agentdock_windows_amd64.zip | 96,600,971 | `a8ac85ed8005fe075ba321e048488b03217e5db5a7a4bfbf985ce346151fc3ca` |
| install.ps1 | 112,799 | `582d05b610b98e66fdbad43bb65c50048c34ac3d33a671603fc8e93ba44c3af4` |

이는 GitHub Actions 후보 artifact의 실제 재다운로드 검증이다. **정식 GitHub Release에서 재다운로드한 결과가 아니며 설치 권장/게시 승인이 아니다.** 위치는 외부 evidence의 `ci-36338434227-candidate/release`다. 실행 결과·실패 로그·Setup 4개 로그·원본 native/Linux 증거를 함께 보존했다. `upgrade-validation.json`의 passed=false를 변경하지 않았다.

## 4.4 세 번째 재개 — 업그레이드 실패 증거와 시험 오판 방지

이번 계속 요청은 실제 HEAD `5bfa5e7d45dc307426ac6beaba9f02a0d2115d50`와 CI 36338434227을 대조한 뒤 진행했다. 4.3절의 Setup 성공과 역사적 baseline 실패는 기존 실제 결과이며, 이번에 새로 통과한 upgrade 결과가 아니다. 새 제품 버전이나 바이너리를 만들지 않았다.

한글화 외부 후보의 충돌을 읽어 현재 출력 제한·승인·삽입 메시지 로직을 보존할 이식 범위를 검토했으나, 외부 `resume3-resolve-korean.py` 작성 요청이 명시적으로 실행 전에 차단됐다. 해당 파일은 생성되지 않았고 한국어 제품 코드를 적용하지 않았다. 같은 작업을 인자 분해나 다른 도구로 재전송하지 않았다. 원문과 실제 파일 부재는 `resume3-localization-write-block.json`에 보존했다.

독립적인 검토에서 기존 upgrade 시험은 실행 단계와 입력 bytes의 SHA를 보고서에 남기지 않았고, 기대한 실패의 비정상 종료만 확인하여 의도된 trial 오류 지점에 실제 도달했는지를 입증하지 않았다. 이 시험 경계만 보완했다. `scripts/test/test-windows-upgrade-isolated.ps1`은 이제 phase·실제 프로세스 종료 코드·스크립트/ZIP SHA와 실행 전후 입력 불변을 기록한다. rollback 주입 시험은 정확한 오류 표식에도 도달해야 한다. 기존 baseline 설치/health, target upgrade, committed pointer, 자격 증명·사용자 데이터 보존과 rolled_back 검사는 그대로 남는다.

증거 보존은 알려진 네 단계의 log/ini, 최대 8개 파일만 대상으로 한다. 파일 1MiB와 본문 524,288 문자 상한을 넘는 입력은 생략 이유를 기록한다. 크기 상한을 넘는 파일은 해시를 위해서도 전체 읽지 않고 source_sha256을 null로 둔다. BOM을 판별하여 UTF-16/UTF-8 원문을 읽고 생성된 fixture의 Bearer/OAuth/Tunnel 자격 증명과 Authorization 값을 제거한다. 실제 원본/제거 후 bytes의 hash를 구분하며 원본을 바꾸거나 이전 export를 덮어쓰지 않는다. 정상·실패 CI artifact 경로 모두 `upgrade-evidence/*`를 보존하도록 기존 workflow에 연결했다. 진단 export 실패를 성공으로 숨기지 않는다.

새 `upgrade_evidence_windows_test.go`는 Windows PowerShell 5.1에서 실제 script의 네 함수만 AST로 추출해 검증한다. 설치 프로세스는 inert 함수로 대체하며 원본 script의 최상위 설치·레지스트리·서비스 동작을 실행하지 않는다. 정상 프로세스 반환을 설치 건강 상태로 오인하지 않기, 엉뚱한 지점의 실패 거부, 주입 지점 도달, 예상 외 성공 거부, 입력 변조 거부 및 민감정보 제거·크기 제한·원문 보존을 43 assertions로 확인했다. 실행 결과는 mock 5개이며 실제 Setup 5회가 아니다.

| 이번 변경 검증 | 실제 결과 |
|---|---|
| Windows script suite | `resume3-scripts-verified`: 110 PASS / 0 SKIP / 0 FAIL |
| PowerShell 5.1 실제 함수 회귀 | 위 suite에 포함: 43 assertions, 5개 inert 결과 |
| `go vet ./scripts/test` | `resume3-vet-verified`: exit 0 |
| PowerShell 구문 / gofmt / diff whitespace | 통과 |
| 새 CI·새 바이너리·실제 역사적 upgrade | 이번 보완 후 미실행; 4.3절의 실패를 변경하지 않음 |

첫 함수 시험의 Get-FileHash 미해결은 PS7 부모의 모듈 경로를 PS5.1 자식이 상속한 시험 환경 오류였다. 자식 PSModulePath를 해당 Windows PowerShell 내장 Modules로만 지정하여 해결했으며 시스템 모듈·PATH·권한을 변경하지 않았다. 첫 실패 기록 `resume3-upgrade-evidence-01`도 보존한다.

직전 종료 기록 `RESUME2_FINAL_STATE.json`이 작성자를 unknown으로 기록한 workflow/upgrade script/새 Go test 세 파일은 이번 계속 요청에서 작성한 파일들이다. 그 당시 관측 기록을 소급 수정하지 않고 이번 검증·커밋으로 귀속을 명시한다. 다른 worktree의 변경이나 열린 PR은 수정하지 않는다. 남은 한글화와 실제 역사적 upgrade/rollback을 마치기 전 main/tag/정식 Release를 갱신하지 않는다.

## 4.5 네 번째 재개 — 기존 한국어 이식과 중복 검증 제거

2026-09-28 사용자의 진행 및 과설계·과검증 방지 요청에 따라 새 시험 프레임워크나 실행 소유자는 추가하지 않았다. 기존 한국어 PR `f60a58bcfffd13940ef168d010ae39b2e8f7772b`의 리소스·표시 구현을 현재 후보에 이식했다. 각 언어의 실제 리소스는 324개에서 1,028개로 늘었으며 Windows 제어판·실행 창·브라우저 상태/OAuth 페이지·MCP Apps·Setup에 연결했다. 출력의 Unicode scalar 예산, 승인 검토 주체, 삽입 수신 상태, 호출 전이의 현재 책임은 유지한다. 영어·한국어 설명을 중국어 기본 문구로 덮던 제목 판별은 제거하고 원문을 유지한다. 선택적 표시 메타데이터는 기존 기록·redactor·projection 경계에서만 처리하며 사용자 입력·로그·외부 도구 응답을 번역하지 않는다.

제품 빌드는 경고·오류 0이었다. 기존 localization-only 회귀는 3개 언어의 실제 리소스와 원문 보존을 검사하여 6,246 assertions를 통과했다. 집중 Go 회귀는 최신 `resume4-focused-final`에서 51 PASS / 0 SKIP / 0 FAIL이다. 최초 pin 검사의 실패는 이식 도구가 공식 언어 파일의 줄바꿈을 정규화한 탓이었다. 기대 해시를 바꾸지 않고 원래 Git blob의 정확한 바이트를 복원하여 해결했다. 기존 순수 정책 시험은 새 표시 계약에 맞춰 전체 원문·상태·출력 크기를 계속 확인한다.

과거 기준판은 기본 `script` 채널 대신 기존 설치기가 제공하는 `setup` 채널로 준비하도록 workflow의 인자 한 줄을 수정했다. 구형 소스나 버전 포인터를 조작하지 않았다. 아직 실제 CI의 설치·업그레이드 성공을 가정하지 않는다. 같은 SHA의 WPF 검사를 native acceptance와 build에서 두 번 실행하던 중복도 제거했다. 필수 native acceptance의 원래 검사·스크린샷과 build 의존성은 그대로 유지하며 검증 범위를 줄이지 않는다.

다만 v1.1.7에서 추가된 권한 설정·새 삽입 상태 등의 잔여 화면 문구를 일괄 보완하는 외부 `resume4-complete-ui-strings.py` 작성 요청이 실행 전에 차단됐다. 원문: “요청의 보안 상태를 결정하지 못해 이 도구 요청은 OpenAI에 의해 차단되었습니다.” 파일 부재를 확인했고 같은 작업을 우회 재전송하지 않았다. 따라서 1,028개 리소스 검사 성공을 전체 화면 한글화 완료로 부풀리지 않는다. 현재 후보는 잔여 문구 미완료이며 정식 게시 대상이 아니다.

제품 커밋 `e9b2c6d23068f41ce0efd64230d8f7d6ceb110cc`은 작업 브랜치에 push했다. 이 정확한 SHA의 실제 CI `36362689792`는 Linux 전체/race, Windows 전체 Go·정적 분석·순수 정책, native 복구 22개 시나리오/272 assertions, WPF 60개 표본/27,538 assertions, ZIP·Setup 빌드 및 패키지 검증을 통과했다. 전체 CI 결과는 **failure**다. 첫 Setup 설치 중 native launcher가 `result.json`을 열 때 `The process cannot access the file because it is being used by another process`라는 파일 공유 위반이 발생했다. 역사적 기준판 단계는 실행되지 않았으므로 `-BaselineChannel setup` 수정의 실제 upgrade 효과도 미검증이다. 앞선 58b6f149의 Setup 성공을 이 후보의 성공으로 합산하지 않는다.

공유 위반이 발생한 기존 결과파일 소유자와 보존 PR `782bb9196b5aee5ee42dcd63eb15ee1b65e9a7e2`를 검토하는 조회도 같은 사전 차단 문구로 실행되지 않았다. 점유 프로세스와 구체적인 코드 원인은 확정하지 않았고 미검토 PR을 이식하지 않았다. 실패를 고치지 않은 채 같은 전체 CI를 다시 실행하지도 않았다. 별도의 로컬 Go fixture 정리 요청과 좁은 툴바의 리소스 조회 역시 실행 전에 차단됐으며, 같은 요청을 다른 경로로 재전송하지 않았다. 따라서 로컬 임시 폴더 제거 및 추가 화면 문구 완료를 주장하지 않는다.

실패 후보 artifact `unverified-windows-failure-36362689792`를 외부 `ci-36362689792-candidate`에 내려받아 기존 검증기로 버전·source SHA·실제 bytes·rg 15.2.0·라이선스·Core Skill 3개를 확인했다. Setup SHA-256은 `8e13a4f1188d35d89b79423d49eeda88e57b511591c5c77cfd0dc1b1c19af4a9`, ZIP은 `85e2918be6100e48aa2d639b425160983336e324a7d80887cd7f1afd66d2fc38`이다. 후보 채널은 `candidate-not-released`, AgentDock 서명 상태는 unsigned다. 이 재다운로드는 배포물 바이트 검증이지 설치 성공이나 정식 Release 게시가 아니다. 운영 PC의 Setup·설치·재시작·네트워크 변경은 수행하지 않았다. 최신 제품 소스는 e9b2c6d이며 뒤따르는 기록 커밋은 제품 바이너리를 새로 빌드한 것으로 해석하지 않는다.

## 4.6 완성 요청 — 남은 표시와 Setup 수신 경계

남은 권한 설정·추가 메시지 상태·검색·기본 동작 이름 등 90개 문구를 기존 UiText에 추가했다. 각 언어 리소스는 1,118개이며, 출력·입력·진단 원문을 다시 쓰지 않는다. 원본 PrivilegeTransition와 TaskAdmin의 기술 진단은 원문으로 보존했다. 실제 WPF 회귀의 언어를 한국어로 선택하여 기존 60개 표본에서 한국어 화면을 확인하며 새 배율 행렬이나 시험 실행기를 추가하지 않았다. 순수 정책 회귀는 기존 문구의 의미를 검증하도록 중국어를 명시한다.

Setup 공유 위반은 기존 broker가 result.json 읽기의 모든 비-부재 오류를 즉시 실패로 처리하는 경계에서 발생했다. ERROR_SHARING_VIOLATION과 ERROR_LOCK_VIOLATION만 기존 50ms polling과 기존 전체 기한 안에서 대기하도록 보완했다. 새 timeout·재시작·성공 추정은 없다. 실제 임시 파일의 독점 handle로 읽기 실패를 재현하고 handle 해제 후 원래 결과가 읽히는지 확인했다. 접근 거부·다른 오류는 즉시 반환하며 기존 nonce·JSON·자식 종료 검사는 유지한다. 실제 점유 프로세스가 무엇이었는지는 확정하지 않는다.

집중 Setup 회귀 12 PASS / 1 helper SKIP / 0 FAIL, 기존 순수 정책 890 assertions와 실제 3개 언어 리소스·원문 회귀 6,786 assertions가 통과했다. 한국어 관문은 후보 CI에서도 기존 검사 결과를 확인하도록 하여 게시 직전에만 누락을 발견하는 불필요한 재빌드를 방지한다. 변경 소스의 최종 설치·업그레이드·게시 성공은 CI 결과로 별도 확정한다. 기존 failed CI와 차단 기록은 보존한다. 이번 native 진단 표시/리소스 조회 요청도 실행 전에 차단되어 같은 조회를 재전송하지 않았으며, 이미 읽은 화면 소스의 독립적인 표시 수정과 Setup 결함 수정만 수행했다.

## 4.7 이관 후 재개 — Setup 설치 후 Core 정지 실패

source `900723d` / CI 36364901575는 Linux·Windows Go·정적 분석·순수 정책·한국어·native 22개 시나리오/272 assertions·WPF 60개 표본/27,538 assertions·ZIP/Setup 빌드와 패키지 검증을 통과했지만 전체 결과는 **failure**다. 실패 단계 로그를 실제로 조회했다. 첫 Setup(당시 중국어) 설치 자체는 성공했고, 직후 원래 E2E의 `service stop`이 `agentdock: inspect Core candidate 896: Access is denied.`로 실패했다(`test-windows-setup-e2e.ps1:274`). 오류에 wrapper가 없으므로 예약작업 `/End` 뒤의 정지 대기 선택에서 발생했다. 반복 설치·repair·제거와 역사적 upgrade는 미도달이다. 앞선 result.json 공유 위반은 이 실행에서 재발하지 않았다.

원인은 2.4절 선택기가 snapshot의 모든 `agentdock.exe`를 열고 접근 거부를 즉시 오류로 반환한 경계다. 다른 사용자의 프로세스나 Core 종료 직후 재사용된 PID가 있으면 상태·정지가 모두 실패한다. PID 896의 실제 소유자는 로그에 없어서 확정하지 않는다. 58b6f149의 같은 선택기가 한 번 통과한 것은 이 경합이 비결정적임을 뜻한다. timeout 추가·재시도·assertion 변경 없이 `ERROR_ACCESS_DENIED`만 종료된 PID와 같이 미선택으로 처리했다.

회귀 `TestCoreSelectionSkipsInaccessibleSameNameProcess`는 빈 DACL로 만든 같은 이름 fixture에 동일 사용자가 실제 `ERROR_ACCESS_DENIED`를 받는지 먼저 확인한 뒤, 선택·정지가 성공하고 이 root의 Core만 종료하며 다른 프로세스는 생존함을 검사한다. 수정 전 코드에서는 CI와 같은 `inspect Core candidate …: Access is denied.`로 실패했고, 수정 후 기존 Core 역할 시험 5개와 함께 6 PASS / 0 SKIP / 0 FAIL이다. 로컬 worktree의 `.git` 소유자가 Administrators여서 Go VCS stamping이 거부되는 환경 문제는 전역 Git 설정을 바꾸지 않고 자식 프로세스의 `safe.directory` 환경값으로만 해결했다.

이전 미커밋 3개 파일은 `FINISH_ATTEMPT_STATE.json`의 SHA-256과 일치함을 확인한 뒤 그대로 `f253e31`로 커밋했다. 해당 로컬 검사(표시 6,789 assertions, scripts/test 112 PASS)는 이관 전 결과이며 재실행하지 않았다. 실제 Setup의 세 언어, 설치 후 Core 정지·반복 설치·repair·제거와 역사적 1.1.16102 upgrade/rollback은 이 소스의 후속 CI 결과로만 판정한다.

## 5. 미완료 및 명시적 차단

직전 재개에서 두 개의 읽기 요청이 실행 전에 차단됐다. 원문은 다음과 같다.

> 요청의 보안 상태를 결정하지 못해 이 도구 요청은 OpenAI에 의해 차단되었습니다.

Windows 한글화 대상 전체 파일/OAuth template/관련 frontend diff 조회와, 별도의 Core CLI/MCP Apps 시험·한글화 source 등의 조회다. 세부 운영 정책 사유는 응답에 제공되지 않았다. `resume-localization-read-block.json`, `resume-mcp-cli-read-block.json`에 실행되지 않았음을 기록했다. 다른 도구·문자 인코딩·인자 분해로 같은 조회를 재전송하지 않았다. 독립적인 최소 결함·시험·문서 작업만 계속했다. 이 제품 패치가 OpenAI 사전 검사를 고쳤다고 주장하지 않는다.

남은 제품 관문은 전체 한글화 검토·구현과 실제 표시/배포물 검증이다. 두 번째 재개에서 원격 다운로드/업데이트 대상과 Windows workflow를 사용자 fork에 맞게 수정했으나, 실제 CI·설치·게시 완료 여부는 개별 실행 결과로 판정한다. 기존의 all-platform 게시 경로로 Windows 외 배포를 확장하지 않는다.

이번 재개에서는 한국어 병합 충돌 조회와, 실제 CI 실패 후 역사적 설치 준비/기존 upgrade harness 조회의 두 요청이 같은 문구로 실행 전에 차단됐다. `resume2-localization-conflict-block.json`과 `resume2-historical-inspection-block.json`에 기록했다. 해당 조회를 우회·분해·재시도하지 않았으며, 한국어 후보와 과거 바이너리를 임의로 적용·수정하지 않았다.

남은 필수 조건은 전체 한국어의 구현/실제 표시/배포물 검증, 동작 가능한 역사적 baseline 준비와 지원 upgrade·실패 rollback 검증, 그 최종 source의 전체 재검증, canonical local/remote main 통합, 동일 SHA 태그·최종 패키지, 정식 Release 게시와 게시 후 실제 재다운로드다. 후보 버전 지정·main 이력 보존 merge·CI/WPF·Setup·후보 bytes 검증은 이미 수행했으므로 미실행이라고 반복하지 않는다. 게시 권한의 재승인이 아니라 실제 남은 구현·검증 조건이다.

현재 후보는 새 버전 1.1.17100으로 지정했으나 Release로 게시하지 않았다. 이 버전은 보존된 1.1.16102보다 크며 Go와 Windows UI 버전을 일치시켰다. Setup/ZIP/metadata/checksum/서명을 최종 소스에 묶기 전 태그를 생성하지 않는다. 필수 표시 및 설치 검사가 빠졌으므로 전체 완료 판정을 내리지 않는다.

## 6. 보존·운영 보호와 다음 복구 지점

정본 main과 독립 후보 `agentdock-workbench-product`는 이번 소스 작업의 변경 대상이 아니다. rebuild worktree의 미커밋 13개 항목은 읽기 참조만 하고 수정·삭제하지 않았다. 열린 upstream PR의 제품 채택 여부와 무관하게 기존 11개 head를 보존한다. 아직 통합하지 않은 작업 브랜치와 후보를 삭제하지 않는다.

직전 로컬 native fixture root 및 Task는 제거와 잔류 0을 확인했다. 이번 실제 CI의 historical baseline 실패 fixture는 로그에서 보존된 것으로 관측됐으며 모든 CI fixture의 정리 성공을 주장하지 않는다. 로컬 Go 시험 root의 소유권·참조 프로세스 확인과 제거 결과는 `resume2-fixture-cleanup.json` 및 추가 종료 기록에 남긴다. 최신 Git/PR/Release 종료 상태는 `RESUME2_FINAL_STATE.json`에 기록한다. 원시 실패·PASS·patch·빌드 결과는 제품 checkout 밖에 보존한다. 이는 같은 D: 볼륨의 로컬 기록이지 별도 재해복구 저장소가 아니다.

운영 Setup·업데이트·제거·복구 및 생산 Core 재시작 명령을 실행하지 않았다. 실제 설치 완료나 지원 upgrade 경로 완료를 native fixture의 성공으로 대체하지 않는다. 정식 Windows Release URL은 아직 없다. 후보 Setup/ZIP와 hash는 4.3절에 실제 상태로 기록했으며 공개 Release 게시와 구분한다.

복구 시 이 문서와 작업 task, 실제 Git/외부 종료 기록을 먼저 대조한다. 검증된 제품 변경을 다시 만들거나 옛 bloated main을 그대로 push하지 않는다. 차단 경계를 우회하지 않고, 아직 남은 한글화와 Windows 배포 완료 조건을 분리해 유지한다.
