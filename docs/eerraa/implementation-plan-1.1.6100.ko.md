# AgentDock Workbench — 원본 기반 최소 제품 및 Windows 배포 정본

기록일: 2026-09-28 KST
작업: `tsk_c83ac7427fcaa1db`
상태: **부분 구현 검증 완료 / 전체 제품 및 Windows Release 미완료**

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
| 직전 전체 검증 소스 | `b1d8589dce326763bf6f8447163b8ba21e60f53e`; 후속 변경의 검증은 4.1절에 별도 기록 |
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

실제 rg 구성요소를 사용한 검색·전달 회귀와 cache-only ZIP 검증은 통과했다. **새 제품 Setup/Release ZIP은 아직 만들지 않았으므로 최종 배포물까지 완료라고 하지 않는다.**

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

### 2.5 TaskAdmin — 목적별 축소·검증

현재 root·SID·task name·정의·실행파일·인자 및 복구자료 대상 검증을 유지했다. 복원 입력 검증은 기존 Task의 정지·삭제보다 먼저 실행한다. 새 Task 정의는 현행 stable tray action만 허용하고, 구형 stable Core action은 기존 task/명시적 복구 입력에서만 allowLegacyCore로 허용한다. 이 제품이 실행하지 않는 `--task-core-host` action은 제외했다.

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

직전 b1d8589 제품 소스의 WB_BASE 대비 차이는 43개 파일, 3,785줄 추가·53줄 삭제다. 여기에는 이전 후보의 한국어·rg 및 시험 코드가 포함된다. 이 숫자를 이번에 새로 완성한 기능 수로 해석하지 않는다. 이번 재개의 소스 커밋 4개만 외부 resume-source-patches에 format-patch 및 SHA-256으로 추가 보존했다.

## 4. 실제 검증과 실패 보존

아래 수치는 서로 더하지 않는다. 전체 Go, 부분 회귀, native assertion, 시나리오 및 패키지는 다른 집계다.

| 검사 | 관측 결과와 근거 |
|---|---|
| MCP/ACL 원본 반례 | 2 tests FAIL, 2 packages FAIL. `minimal-baseline-01` |
| MCP/ACL 최초 수정 후 | 148 PASS / 2 SKIP / 0 FAIL, 3 packages PASS. `minimal-fixed-01` |
| Core 원본 반례 | 2 tests FAIL. `core-role-baseline-01` |
| Core 수정 후 전체 관련 패키지 | 200 PASS / 2 SKIP / 0 FAIL. `core-role-fixed-01` |
| 구형 복구 조기 재개 반례 | schema 1의 Run 호출을 기록하는 test double에서 FAIL; 실제 Task 시작 없음. `native-legacy-baseline` |
| manifest 취소 반례 | 1 FAIL; 수정 후 bundle 관련 suite 통과. `rg-cancel-baseline-01`, `rg-script-fixed-01` |
| 전체 Go 현재 소스 | `go test -json -p 1 -count=1 -timeout=180s ./...`: 2,135 PASS / 85 SKIP / 0 FAIL, 시험 패키지 56 PASS, 시험 없는 패키지 8. `resume-full-go-02` |
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

## 5. 미완료 및 명시적 차단

직전 재개에서 두 개의 읽기 요청이 실행 전에 차단됐다. 원문은 다음과 같다.

> 요청의 보안 상태를 결정하지 못해 이 도구 요청은 OpenAI에 의해 차단되었습니다.

Windows 한글화 대상 전체 파일/OAuth template/관련 frontend diff 조회와, 별도의 Core CLI/MCP Apps 시험·한글화 source 등의 조회다. 세부 운영 정책 사유는 응답에 제공되지 않았다. `resume-localization-read-block.json`, `resume-mcp-cli-read-block.json`에 실행되지 않았음을 기록했다. 다른 도구·문자 인코딩·인자 분해로 같은 조회를 재전송하지 않았다. 독립적인 최소 결함·시험·문서 작업만 계속했다. 이 제품 패치가 OpenAI 사전 검사를 고쳤다고 주장하지 않는다.

남은 제품 관문은 전체 한글화 검토·구현과 실제 표시/배포물 검증이다. 두 번째 재개에서 원격 다운로드/업데이트 대상과 Windows workflow를 사용자 fork에 맞게 수정했으나, 실제 CI·설치·게시 완료 여부는 개별 실행 결과로 판정한다. 기존의 all-platform 게시 경로로 Windows 외 배포를 확장하지 않는다.

그다음 새 버전 일관성, 원본 대비 최종 tree 재검토, 기존 main 이력 보존 통합, 실제 격리 runner의 WPF/신규 설치/지원 이전 버전 업데이트·복구, 동일 commit Windows Setup·ZIP, 정상 main/tag push, 정식 Release 및 실제 새 디렉터리 재다운로드 검증이 필요하다. 이미 승인된 게시 권한을 다시 묻는 단계가 아니라 아직 구현·검증되지 않은 완료 조건이다.

현재 후보는 새 버전 1.1.17100으로 지정했으나 Release로 게시하지 않았다. 이 버전은 보존된 1.1.16102보다 크며 Go와 Windows UI 버전을 일치시켰다. Setup/ZIP/metadata/checksum/서명을 최종 소스에 묶기 전 태그를 생성하지 않는다. 필수 표시 및 설치 검사가 빠졌으므로 전체 완료 판정을 내리지 않는다.

## 6. 보존·운영 보호와 다음 복구 지점

정본 main과 독립 후보 `agentdock-workbench-product`는 이번 소스 작업의 변경 대상이 아니다. rebuild worktree의 미커밋 13개 항목은 읽기 참조만 하고 수정·삭제하지 않았다. 열린 upstream PR의 제품 채택 여부와 무관하게 기존 11개 head를 보존한다. 아직 통합하지 않은 작업 브랜치와 후보를 삭제하지 않는다.

native fixture root 및 Task는 최종 시험에서 제거되고 잔류 0임을 확인했다. 완료한 Go 시험의 고유 root 잔류 및 추가 종료 점검은 외부 `RESUME_FINAL_STATE.json`과 `resume-fixture-cleanup.json`에 관측값을 기록한다. 원시 실패·PASS·patch·빌드 결과는 제품 checkout 밖에 보존한다. 이는 같은 D: 볼륨의 로컬 기록이지 별도 재해복구 저장소가 아니다.

운영 Setup·업데이트·제거·복구 및 생산 Core 재시작 명령을 실행하지 않았다. 실제 설치 완료나 지원 upgrade 경로 완료를 native fixture의 성공으로 대체하지 않는다. Windows Release URL·신규 설치파일 다운로드 링크·게시 hash는 아직 존재하지 않으므로 작성하지 않는다.

복구 시 이 문서와 작업 task, 실제 Git/외부 종료 기록을 먼저 대조한다. 검증된 제품 변경을 다시 만들거나 옛 bloated main을 그대로 push하지 않는다. 차단 경계를 우회하지 않고, 아직 남은 한글화와 Windows 배포 완료 조건을 분리해 유지한다.
