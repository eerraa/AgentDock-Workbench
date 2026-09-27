# AgentDock — 업스트림 우선 수렴과 로컬 통합

정본: `docs/eerraa/implementation-plan-1.1.6100.ko.md`  
실행일: 2026-09-27 KST  
상태: **실험 분리 및 선택한 upstream 소스 통합·필수 로컬 검증·로컬 main 정상 반영 완료**. 최초 27개 충돌을 모두 해결했다. 전체 Go suite와 Windows desktop·다국어·실제 격리 task/복구 검증을 거쳐 normal merge commit을 만들고, 변경되지 않은 clean main에 fast-forward했다. 새 패키징·운영 설치·원격 게시 및 아래 명시한 추가 upstream 변경은 완료 범위와 구분한다.

## 1. 현재 기준과 권한

| 항목 | 고정 값 |
|---|---|
| LOCAL_BASE | `4d719ce75ac748cb24e8ad4b7d92ff7c56e24225` |
| WB_BASE | `b367eaab95202873fb213b8713440bf7822878c4` — Workbench v1.1.7 |
| UVWT_REVIEW_HEAD | `0111fc053acbf1ccdcdb0d0f29d2b6ec0f006269` |
| 검토한 Workbench 최신 main | `a2b4b515e0a09d38ce359ef33b2502b5319115f5` — 기준선으로 채택하지 않음 |
| uvwt 최신 공개 릴리스 | v0.9.1; tag object `74b341c5b38f3d3b4f98697e5555debd7a5e4575`, peeled commit `637ea64bcf3c4177158a0b1806bf9bde545ee4d4` |
| 통합 브랜치 | `integration/upstream-convergence-2026-09-27` |
| 제품 worktree | `D:\Engineering\worktrees\agentdock-maintenance` |
| 기존 배포 식별 | 1.1.16102; 이번 변경의 새 배포 버전은 미정 |

주 추종은 `A-m-o-r-F-a-t-i/AgentDock-Workbench`, 보조 검토는 `uvwt/agentdock`이다. 기존 `A-m-o-r-F-a-t-i/agentdock` URL이 Workbench로 redirect됨을 공개 API로 확인했다. `workbench`와 `uvwt` remote를 분리하고 태그는 `refs/upstream-tags/<remote>/`에 fetch했다. 기존 tag/remote ref의 무심한 덮어쓰기·prune·원격 변경은 하지 않았다. 작업 중 기준을 계속 최신 HEAD로 움직이지 않는다.

기존 df7c22f 기반 1.1.6 고정 계획, 자체 구현 형태의 우선 보존, 기능별 다수 브랜치 및 기존 배포 버전으로 Setup을 다시 만드는 지시는 이번 사용자 요청으로 대체한다. df7c22f는 현재 두 이력의 merge-base일 뿐 LOCAL_BASE가 아니다. AGENTS.md의 해당 낡은 실행 지시도 이 정책에 맞춘다. 인증·명시적 deny·권한·사용자 데이터·ACL·원문·복구·게시 승인 경계는 완화하지 않는다.

허용: 실제 상태 확인, 공개 upstream 조회/fetch, 독립 로컬 백업과 검증, 검증된 실험 전용 대상의 정리, 코드/회귀/정본 수정, 빌드/검사, 로컬 commit과 검증 후 main 정상 merge.

금지: push/PR/CI dispatch/tag/Release 게시, upstream 저장소 변경, 공유 이력 reset/rebase/force push, Setup 실행, 생산 Core 교체·재시작·운영 설정 변경, 다른 프로젝트·공유 journal 삭제. 새 배포 버전이 정해지지 않았으므로 설치파일은 만들지 않는다. 소스 compile은 패키징·설치·게시 완료가 아니다.

## 2. 기준선 검증 근거

정확한 SHA의 GitHub Actions workflow 결과를 조회하여 보관했다. upstream 시험과 우리 통합본의 시험을 구분한다.

| 기준 | 실제 확인 | 판단 |
|---|---|---|
| WB b367eaab / v1.1.7 | CI 36129347704, CodeQL 36129347772, Windows Installer 36129347567, 수동 package 36129349697, 수동 publish 36132246715 모두 success | 최근의 일관된 릴리스 기준으로 채택 |
| 같은 b367eaab의 tag-trigger 36132172894 | failure | 위 성공과 별도 보존. 상세 실패 원인은 확인하지 못함 |
| WB a2b4b515 / 최신 main | CI 36250662527·CodeQL 36250662581 success; Windows Installer 36250662532·package 36253143050 failure | Windows 제품 기준으로 곧바로 채택하지 않음 |
| uvwt 0111fc05 | CodeQL 36244370445 success | 보조 소스 검토 기준이며 전체 제품 검증 성공을 뜻하지 않음 |

CI 상세 jobs 조회 한 요청은 OpenAI 안전 검사에 의해 pre-dispatch 차단됐다. 같은 요청을 재시도하거나 다른 App/경로/인자로 우회하지 않았다. workflow 수준 결과만 근거로 사용하며 job-level 원인은 unknown이다. 이 통합이 OpenAI pre-dispatch 차단을 해결한다고 주장하지 않는다.

두 upstream의 patch-id 차이만으로 미반영 기능 수를 판단하지 않는다. 동등 구현·현재 호출자·실제 부족한 요구와 의존성을 확인해 Workbench에 없는 필요한 uvwt 수정만 출처를 남기고 적용한다.

## 3. 실험 보관 및 정리 결과

독립 실험 Git: `D:\Engineering\agentdock-experiments`  
접근 제한 원본 보관: `D:\Engineering\archives\agentdock\2026-09-27`

독립 저장소는 자체 .git/objects와 ref를 가진다. product의 linked worktree/submodule이 아니며 alternates 의존이 없다. 원래 모든 79개 ref/SHA를 대조하고 `git fsck --full`을 통과했다. 원본 diag ref는 `63ad4aa574593971c5aec2e8577fe126cb786bfe`로 남겼다. 별도 `archive/studies-2026-09-27`의 `19bedb9677cdd613c9ab9f260e40376c77dc19bf`는 LONG240와 COMMON240-v2를 서로 다른 디렉터리에 원문 보관한 commit이다. 시험을 재실행하거나 run/study 결과를 합산하지 않았다.

원본 파일·ignored 산출물·세 worktree·Git 메타데이터·staged/unstaged binary diff·untracked/ignored 목록·등록 정보·선별 실험 task 원본을 보관했다. 총 8,338파일 / 1,283,125,207bytes를 SHA-256 및 원본/보관본 전체 경로 집합으로 대조했다. 실험 ZIP 5개는 CRC 검사와 실제 해제 후 모든 항목 읽기·해시를 검증했다. 독립 diag worktree 1,858개 파일도 원본과 일치했다. 복원 중 세 파일의 줄바꿈/stat 상태를 내용 diff와 index blob 표 불변을 확인한 뒤 refresh했다. 원문 바이트는 변경하지 않았다.

bundle SHA-256: `85d4d3373343604eee23e810f325db91f3d21ed1efe5baccc2d53ec3b333480b`. 전체 file-manifest.jsonl에는 원본/보관 경로, 크기, 해시, 출처, 확인 가능한 study/run을 기록했다. source_commit은 worktree 귀속이며 개별 산출물의 build origin이라고 주장하지 않는다. 모르는 버전은 null이다. 같은 D: 볼륨의 로컬 보관이며 별도 디스크 재해복구 사본이 아니다.

삭제 전 등록 원문을 보존하고 oai_bench만 disable하여 TestKit와 자식 writer 종료를 확인했다. 검증 후 다음만 제거했다: dist/oai-long240, dist/oai-unified-v2, dist/oai-safety-diag-1.1.16102, dist/AgentDock_OAI_Long240_Kit_KO.zip, diag worktree와 해당 local ref, oai_bench 등록. worktree는 Git 관리 명령을 사용했고 ref는 예상 SHA 대조 후 삭제했다. 원격 ref를 삭제하지 않았다.

보존: 정상 1.1.16102 설치/rollback 자료, 제품 PR worktree·브랜치, 제품 소유 dist/maintenance, 회귀 fixture/testdata/라이선스, 다른 프로젝트와 생산 데이터. 완료된 실험 task 2개만 지원 API로 archive했다. 미완료 실험 task 6개는 API가 completed 상태만 archive하므로 그대로 남겼다. 결과를 성공/취소로 조작하지 않았다.

미정리: 현재 연결을 유지하는 생산 Core가 소유한 oai_fixture, 공유 journal/diagnostics, 위 미완료 task. 해당 메모리/공유 기록의 선택적 export·정리 완료를 주장하지 않는다. 확보하지 못한 과거 첨부는 재구성하지 않았으며 로컬 보관본 밖의 존재 여부는 unknown이다. 생산 Core와 운영 설치는 변경하지 않았다.

## 4. 채택한 구현과 잔여 제품 차이

공통 구현은 고정한 Workbench에 수렴한다. 최초 27개 충돌을 모두 해결했으며 `git ls-files -u`가 비어 있다. 파일 전체를 한쪽 것으로 일괄 대체하지 않고 실제 책임과 사용자 계약에 따라 병합했다. 새 실행기, task 시스템, native typed extension, 모델 API controller 또는 A–G 시험은 추가하지 않았다.

### 4.1. 중복 제거와 공통 책임

- activity binding·실행 사실·수명 전이는 Workbench `call_transitions`로 모았고 journal은 batch admission을 사용한다.
- insertion의 중복 `Summary`/`Summarize` 및 footer를 제거하고 Workbench timeline/presentation을 사용한다. 추가 메시지 원문·호출/대화 ID·만료·수신 확인을 보존하며 재조회·재전달이 원 도구를 실행하지 않는다.
- Core Skill builder/bootstrap/version 검증을 Workbench 자산으로 통합하고 quoted-empty/nested-only 반례를 보존했다.
- Windows privilege 전환은 Workbench `PrivilegeTransition`과 `TaskSecurityDescriptor`가 담당한다. 사용하지 않는 C# `KillOnCloseJob`은 제거했다. native Go의 이미 존재하는 process owner를 사용하며 두 번째 호스트를 만들지 않았다.
- MCP manager/catalog, Plugin manager·snapshot·format·types·tool service, managed Skill state/install/uninstall/tool, app runtime, permission, 위 전환/SD owner는 실제 `git diff WB_BASE -- <scope>`가 0이다. 경로와 결과는 `completion-132145/owner-convergence-review.json`에 있다.

이전 중단 기록의 “Workbench Plugin은 versioned storage”라는 포괄 서술을 정정한다. Plugin은 `plugins/<name>` 직접 저장과 host state를 사용하고, managed Skill의 versioned state와 구분한다. 주 추종의 관리·수명·rollback owner 자체가 그대로 채택되어 있으므로 별도 uvwt 저장 runtime이나 형식 전환을 추가하지 않는다. 필요한 복사 보완만 아래의 독립 함수 집합으로 연결했다.

### 4.2. 보조 upstream의 실제 적용 범위

| 원본 repository·commit | 채택 범위 | 의존성/검증 |
|---|---|---|
| uvwt/agentdock `18e12e5592efea3d83c85de6ae23f20d611e519b` | browser WebSocket URL 대기를 요청 timeout에 결합 | 원 소스/Unix 회귀의 두 blob 정확히 일치; Windows disposable child timeout·회수 회귀 추가 |
| uvwt/agentdock `a4848d1908638aa7b0a0f9fae859f03d6251eb4c` | managed/Plugin/workspace Skill 설명 원문 보존 | WB context snapshot/filesystem owner에 adaptation. 수정 전 3개 truncation FAIL 재현; 수정 후 app PASS. shared 50개/120 bytes 기존 제한은 별도 유지 |
| uvwt/agentdock `120527c7a385e24fcbe98ef1786c824b086e7aae` | `snapshotPluginTree`, `validatePluginSnapshotPath`, `snapshotPluginRegularFile` | 세 함수 원문 동일, 의존성 완결 단위. 기존 49행 copy loop는 한 호출로 대체. os.Root·열기 전후 identity·실제 byte 상한·O_EXCL을 채택; 새 snapshot 11개 PASS |

Plugin 복사 상한 20,000은 upstream처럼 디렉터리를 포함한 항목 수로 바뀌었다. 512 MiB는 유지한다. 후속 metadata-ignore 정책이나 새 저장 schema는 도입하지 않으며 원문 metadata를 보존한다. 이 제한의 차이는 숨기지 않는다.

24개 계보상 고유 commit은 기능 수가 아니다. 정확한 변경 경로/patch는 `uvwt-review-history.txt`, `uvwt-path-review.json`과 원 commit별 review patch에 남겼다. 처리는 다음과 같다.

| 나머지 변경군 | 이번 결정 |
|---|---|
| Plugin/Skill 수명·상태/저장/계약 (`c3db3f15`, `120527c7`, `e121e69b`, `b36b40c2`, `436d0b0c`) | 주 추종의 기존 완결 관리/수명 구현을 채택하고 실제 Plugin/Skill/app 회귀로 확인. 별도 저장 모델/중복 runtime은 도입하지 않음. 안전한 snapshot 함수 집합만 위처럼 채택 |
| 새 chat/workspace 카드, 이미지 context, Remote MCP OAuth | 필수 제품 요구의 누락으로 확인되지 않은 추가 기능이므로 이번 안정화에 도입·활성화하지 않음 |
| SignPath/저자 Release 파이프라인, 자체 버전 bump | eerraa의 수동 승인 배포 pipeline과 식별을 대체하지 않음. 기존 SHA·hash·tag 검증과 명시적 게시 승인을 유지 |
| Windows 후속 cold-start/migration/cleanup (`5b354d94`, `37e7e4b6`, `e026a712`, `8625d8c6`) | 추가 변경 수용은 보류. 공개 patch는 확인했으나 현재 파일 조합의 추가 조회가 차단되어 해당 의존성 전체의 검증을 완료했다고 하지 않음. 고정 WB의 이번 통합/검증 범위를 유지하며 불완전한 부분 이식은 하지 않음 |

보류된 변경은 이미 포함되었다거나 불필요성이 모두 입증되었다고 기록하지 않는다. 특히 추가 migration 변경의 수용은 별도 후속 검토 대상이며, 이번 소스 통합 성공을 전체 uvwt HEAD의 수용 또는 운영 migration 성공으로 확대하지 않는다.

### 4.3. 남기는 제품 차이

| 실제 사용자 요구 | 선택한 WB만으로 부족한 경계 | 최소 차이 | 관련 검증 | 제거 조건 |
|---|---|---|---|---|
| 자체 배포 식별·오프라인 수동 업데이트·명시적 게시 승인 | 저자 배포와 eerraa의 소유자/승인 정책이 다름 | distribution·고정 다운로드·수동 gate와 immutable build evidence 연결 | scripts/test 92 PASS, cache-only component 8 PASS | 사용자의 배포 정책 변경 또는 upstream의 동일 정책 지원 |
| 한국어 UI와 실제 입력/출력 원문 보존 | 한국어 리소스·원문에 hash 결합한 표시 metadata가 없음 | 기존 resources/OwnedText를 WB formatter/timeline/HTTP UI에 연결 | 3개 locale key/format parity·원문/receipt 검사, policy 7,156 assertions, 실제 WPF 63 samples | upstream이 같은 한국어/원문 계약 제공 |
| 잘못된 discovery가 Core를 중단하거나 last-good를 지우지 않음 | SDK v1.7.0 ListTools는 conversion guard보다 먼저 null 항목을 참조함 | sending middleware의 최소 validation/구조화 오류 보존 | hook 제외 overlay에서 panic, 실제 client 42 PASS; manager/catalog WB 동일 | SDK/상류가 참조 이전 검증과 같은 회귀 제공 |
| 실행/복구 결과의 nonce·오류·exit 원문 정확성 | 유효 receipt 재시도와 잘못된 identity/누락 exit를 구분해야 함 | 기존 receipt reader에 nonce·bounded JSON·pending sharing/lock 구분 | native receipt/broker 선택 14 PASS, exit·잘못된 nonce·원문 반례 | 같은 읽기/복구 계약이 상류에 제공 |
| Windows 상태가 unknown/stale를 성공/종료로 오인하지 않음 | 연결 origin·프로세스/파일 세대·provider 변경과 부정확한 JSON에 대한 제품 계약 필요 | 기존 strict health/passive version 및 native/Tailscale observation 경계, UI는 한 display owner 사용 | health/version 95, native status 39, display 27, Tailscale 31+11 assertions | 동일 identity·unknown·invalid/cancellation 계약 제공 시 해당 보완 제거 |
| Core 종료 시 첫 명령 이전부터 자식 수명 소유, Tunnel과 분리 | 실행 뒤 Job attach만으로는 생성 직후 gap을 보장할 수 없음 | 기존 native creation-time Job/supervised host만 유지; 미사용 C# Job wrapper 제거 | 실제 child/forced host 종료·Tunnel 분리·shim 회귀 PASS | upstream이 creation-time 소유와 해당 동작 제공 |
| 사용자/설치 root에 귀속된 기존 복구 자료 보존 | 로컬 schema 1의 root/SID와 새 schema 2 무결성을 함께 지원해야 함 | TaskAdmin의 좁은 binding/schema adapter; 전환·SD 비교는 WB 동일 | 실제 COM/NTFS 21 scenarios·223 assertions, backup/SD 78, ownership 29 | upstream이 같은 과거 데이터와 소유권 guard 수용 |
| 재현 가능한 도구 공급 및 검증 중 변경 금지 | PATH의 외부 rg 선택과 자체 공급/rollback 요구가 다름 | 실행 세대에 결합한 pinned rg·license/manifest/hash·read-lock 선택 | 실제 번들 4개 패키지 370 PASS, cache-only 8 PASS, NTFS 11 PASS | upstream이 동일 번들/무결성/rollback 계약 제공 |

`--local-isolated`와 `AGENTDOCK_TEST_LOCAL_METADATA=1`은 격리 테스트 전용 명시 opt-in이다. GitHub 환경을 위조하지 않고 기존 CI 조건과 assertion을 유지한다. 실제 task는 고유 acceptance 이름과 임시 root만 쓰며 새 UAC·생산 Core·설치/서비스 설정을 변경하지 않는다. TaskDefinitionPolicy는 이름만으로 권한을 부여하지 않고 정확한 root/실행파일/인자/interactive SID를 계속 검사한다.

### 4.4. 데이터와 rollback 경계

Task Scheduler가 기본 LeastPrivilege의 RunLevel을 생략하는 실제 직렬화를 확인했다. 생략은 LeastPrivilege로만 해석하며, 일반 최고권한 검사는 여전히 이를 거절한다. 명시적인 기존 task 검증/복원에서만 standard/highest를 허용하며, 다른 SID/root/추가 인자/다른 action/Password logon은 거절한다. Core 상태가 unknown이면 privilege 변경 전에 실패한다.

새 Task backup은 schema 2의 XML digest/원 SD에 root·task name·SID를 함께 기록한다. 기존 bound schema 1은 원문을 변경하지 않고 읽고 검증한다. 다른 root/user의 backup, 누락/null 식별, tampered XML, unbound absent-task 기록을 거절한다. Native unknown에서는 충돌하는 복구를 시작하지 않고 evidence를 보존한다. 복원 후 정의와 SD를 검증하며 명시적 deny 또는 추가 권한을 허용 목록 정규화로 숨기지 않는다.

**실행파일만 되돌리는 것은 데이터 rollback이 아니다.** schema 2를 생성한 뒤 이전 1.1.16102 helper가 그 자료를 읽는다고 가정하지 않는다. 전환 완료 또는 현재 버전 helper의 검증된 복구를 먼저 마치고, 현재 schema 2/transition 기록과 원 schema 1을 별도로 보존한 후 실행파일을 되돌려야 한다. 이번에는 운영 전환/Setup을 실행하지 않아 생산 복구 자료를 새 schema로 변환하지 않았다. 일반 installer journal/schema 변경도 선택한 WB 단위로 수용하며 실제 운영 downgrade는 별도 배포 검증 범위다.

## 5. 실제 검증 결과

외부 원시 근거: `D:\Engineering\archives\agentdock\2026-09-27\completion-132145`.
숫자는 test/subtest 이벤트 또는 명시한 assertion 수다. 부분/반복 실행을 더해서 시험 수를 부풀리지 않는다.

| 검사 | 실제 결과 |
|---|---|
| 최종 `go test -json -p 1 -count=1 -timeout=180s ./...` | **PASS / exit 0**. 57개 시험 패키지, test/subtest 2,217 PASS, 기존 조건 SKIP 81. 시험 없는 패키지 7개는 별도 집계 |
| 최종 `go vet ./...` / `go build ./...` | 각각 **PASS / exit 0** |
| main checkout `scripts/test` / `go build ./...` | 병합된 실제 main 경로에서도 **92 PASS / build PASS**; integration과 제품 tree 동일 확인 |
| 변경 Go formatting / 전체 staged diff | **182개 Go 파일 gofmt PASS / git diff --cached --check PASS** |
| scripts/test 전체 | **92 PASS, 0 FAIL, 0 SKIP**. 이전 6개 실패 해소; 제거된 내부 helper assertion은 실제 owner 및 행동 검사로 이관 |
| Windows desktop/native/layout/policy Release win-x64 compile | **PASS**, 최종 compiler 경고/오류 0 |
| 다국어 pure-policy | **7,156 assertions PASS**. en/zh-CN/ko-KR key·format slot parity, receipt/provenance/Unicode·CRLF·공백 원문 확인 |
| task ownership / backup·SD 호환 | **29 / 78 assertions PASS**. schema1/2, 원문 불변, 다른 root/SID/name·XML변조·unbound absence 거절 |
| 실제 격리 Windows task/권한 전환·복원 | **21 scenarios, 223 assertions PASS**. native COM/NTFS, standard/elevated 양 방향, prepare/apply/verify/cancel/rollback/unknown/tamper/absent/deny 검증 |
| 실제 offscreen WPF | **29,032 assertions, 63 rendered samples PASS**. 기존 input/layout과 세 locale 표시. 운영 Core/tray/visible window 실행 없음 |
| 상태·health·Tailscale·summary | health/version 95, summary 34, native status 39, display 27, Tailscale 31 및 실제 runtime 연결 11 assertions PASS |
| 실제 native process/shim 수명 | first-instruction Job·부모 종료·descendant·Tunnel 분리 등 선택 회귀 PASS; 원 도구 재실행 없는 session/receipt 기존 전체 Go 회귀 포함 |
| 실제 pinned rg 통합 | **4개 패키지 370 PASS**, 플랫폼/helper 조건 SKIP 37. 설치된 구성 요소를 source spec 4개 파일 hash/크기로 검증하여 임시 복사한 뒤 실행; 생산 Core 실행 아님 |
| 캐시 전용 rg 구성 요소 검증 | **8 PASS, 0 SKIP**. 다운로드/다른 process를 금지한 유효/변조/partial/license/architecture 반례. AgentDock 설치파일 생성 아님 |
| 실제 NTFS backup/복구 | **11 PASS, 0 SKIP**. owner/ACL/attributes·ADS/EA 비파괴 거절·sharing lock·변조 metadata 거절 |
| Race detector 추가 검사 | **실행 환경 실패**. 시험 진입 전에 모든 test process가 `0xc0000139`로 종료, 0 test PASS. 설치된 gcc 외 대체 compiler는 확인되지 않음. race PASS라고 하지 않으며 제품 assertion 실패와 구분 |
| Setup/생산 Core 교체·재시작/운영 설치·migration/게시 | **NOT-RUN**, 사용자 허용 범위 밖. native 격리 복구 성공과 구분 |
| ARM64/macOS/Linux 실기기 및 물리 DPI·실제 키보드 | **NOT-RUN**, 현재 Windows x64 로컬 결과를 다른 플랫폼 성공으로 확대하지 않음 |

최초 마지막 묶음 Go 실행에는 기존 `TestWorkspaceContextSkillIndexTruncatesWithWarning`가 9.99초 후 context timeout으로 실패했다(2,216 PASS/1 FAIL). 당시 다른 verification이 겹쳤다. 해당 실행을 그대로 보존하고 모든 다른 verification 종료를 확인한 뒤 소스·assertion·timeout 불변으로 전체 suite를 단독 실행하여 위 PASS를 얻었다. 리소스 경합이 원인이라고 확정하지 않으며 이 간헐적 timeout 기록을 숨기지 않는다.

초기 Windows XML/connection 오류 및 첫 native 실행의 10개 실패도 덮지 않았다. 충돌/남은 connection 변수는 실제 owner에 맞게 수정했고, native 실패는 Scheduler 기본 RunLevel 생략을 lowest로 처리하여 원래 반례들을 통과시켰다. 안전 검사 차단과 compiler/test 오류를 구분했다. 이번 두 추가 소스 조회의 pre-dispatch 차단은 해당 요청을 재전송하지 않고 그대로 기록한다. 이 통합이 OpenAI pre-dispatch 차단을 해결했다는 주장은 하지 않는다.

과거 검증은 `partial-validation`, `continuation-validation`, `upstream-fixes-validation`, `mcp-convergence-121938`, `compatibility-convergence-122052`, `plugin-snapshot-convergence-122456`, `windows-completion-gate-131137`에 그대로 있다. 최신 PASS는 이전 FAIL의 기록을 소급 변경하지 않는다. 새로운 240단계 campaign은 없다.

## 6. 소스 반영과 배포 경계

통합 source의 필수 로컬 compile/회귀·권한/원문/데이터 복원 관문을 통과했다. main은 병합 직전에도 최초 LOCAL_BASE 그대로 clean이었다. 공유 이력 reset/rebase/force 또는 ours merge 없이 정상 반영했다.

- 제품 통합 merge commit: `4af02758995634861d3a9e0d93bd60e85cff1098`.
- merge 부모: `23b19952ff4196bc866aa3d87006fb8a27c3df1e`, `b367eaab95202873fb213b8713440bf7822878c4`.
- 검증한 제품 통합 tree: `12fcce09efab6efbd332625b5761e5df11cbe758`.
- 최초 main 반영: 2026-09-27 14:02 KST, `git merge --ff-only`로 위 commit을 그대로 수용.
- 이 완료 기록은 integration branch의 문서 전용 후속 commit으로 남기고 main에 정상 fast-forward한다. 제품 source/test blob은 위 검증 commit과 같다. 최종 main SHA와 Git 복원/파일 해시 결과는 외부 `completion-132145/final-state.json`에 기록한다.
- main/integration은 clean 상태로 종료한다. product tree에 실험 파일·raw 로그·임시 설치파일을 추가하지 않았다. 기존 product PR/worktree와 정상 설치/rollback 자료를 보존한다.

원시 테스트/실패·재검증·build·screenshots와 정확한 source manifest는 외부 archive에 있다. 이번 native fixture의 원래 임시 root 두 개는 종료 점검에서 이미 존재하지 않았고, 해당 root를 참조하는 task/writer는 0이었다. 보존된 JSON/XML 복구/변조 기록은 외부 native-recovery 결과에 남아 있다. 원본 삭제를 이번 점검에서 수행했다고 기록하지 않는다. 생산 oai_fixture/공유 journal와 완료되지 않은 과거 실험 task는 §3의 보존 상태 그대로다.

반영 상태: **로컬 소스 통합·검증·main 반영 완료**. 이 상태는 운영 설치나 모든 후속 upstream 수정의 수용 완료를 뜻하지 않는다.

새 배포 버전은 미정이다. 기존 1.1.16102는 마지막 배포 식별이며 새 소스의 설치파일을 같은 이름으로 만들지 않았다. 로컬 소스 통합, 설치파일 패키징, 실제 설치/rollback, 원격 게시는 서로 다른 상태다. 다음 배포 단계에서는 새 버전과 정확한 source commit을 먼저 명시하고, 보류된 후속 Windows 변경/운영 migration 및 필요한 플랫폼/race 환경 검증 범위를 확정한다. 사용자 승인 없이 Setup·Core 재시작·push·PR·dispatch·tag·Release를 수행하지 않는다.
