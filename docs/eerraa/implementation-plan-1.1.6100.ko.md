# AgentDock — 업스트림 우선 수렴과 로컬 통합

정본: `docs/eerraa/implementation-plan-1.1.6100.ko.md`  
실행일: 2026-09-27 KST  
상태: 실험 독립 보관·검증·선택적 정리 완료. 제품 통합과 회귀 검증 진행 중. 아직 main에 병합하지 않았다.

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

## 4. 통합 원칙과 구현 범위

LOCAL_BASE에서 만든 외부 worktree에 WB_BASE를 정상 merge한다. 공통 기능은 필요한 동작·권한·데이터 호환·복구 조건을 충족하는 upstream 구현을 채택하고 중복 자체 helper/캐시/수명/우회책을 제거한다. 자동 병합된 부분도 중복 책임을 확인한다. upstream 구현이 부족한 실제 사용자 요구만 좁은 downstream 변경으로 유지한다. 익숙함·먼저 구현함·과거 PASS만으로 자체 구현을 보호하지 않는다.

discovery nil/last-good, receipt retry, runtime ownership, health/cache, Tailscale, Plugin/Skill lifecycle, 출력·insertion, 설치·migration·rollback 모두 대체 검토 대상이다. 특정 클래스/함수명·기존 timeout 선택을 제품 요구 자체로 간주하지 않는다. 실행·취소·종료·재조회·conversation/call 귀속·unknown 정확성, 한국어 제품 표시와 실제 입력/출력 원문, 재현 가능한 core Skill/도구 공급, 자체 distribution 및 수동 게시 경계는 보존한다.

새 모델 API controller, 실행기/task 시스템, A–G 변형, native typed extension 실험은 추가하지 않는다. 새 승인 방식·자동 reviewer·OAuth provider·업데이트/게시 동작은 소스 도입과 운영 활성화를 분리한다. 저장 schema 변경은 migration/rollback 호환을 확인하며 구 실행파일로 되돌리는 것만으로 새 데이터 복원이 된다고 가정하지 않는다.

수용하지 않은 upstream을 `merge -s ours`로 숨기지 않는다. `-X theirs`나 전체 파일 일괄 덮어쓰기로 충돌을 처리하지 않는다. 필요한 uvwt adaptation은 repository·원commit·변형 이유를 기록한다. 코드 통합 후 남는 자체 차이만 아래 표에 실제 근거로 확정한다.

| 실제 사용자 요구 | upstream만으로 부족한 이유 | 남기는 최소 변경 | 검증 | 제거 조건 |
|---|---|---|---|---|
| 자체 배포 식별·오프라인 수동 업데이트·게시 승인 | Workbench 배포와 eerraa 배포의 소유자가 다름 | 기존 distribution 경계와 승인 gate를 검토 후 유지 | 통합본 검증 전 | 사용자가 배포 정책을 명시적으로 변경할 때 |

다른 자체 패치의 유지 여부는 아직 확정하지 않았다. 미검토 항목을 최종 잔여 차이로 승인한 표가 아니다.

## 5. 검증과 종료 관문

upstream 시험 자산을 우선 재사용한다. 내부 함수명을 고정한 시험은 동작 계약 중심으로 정리하되 assertion 약화·skip·행동 회귀 삭제로 실패를 없애지 않는다. 변경된 공통 admission·저장·직렬화·응답 조립의 영향을 고려해 관련 Go suite와 Windows desktop compile/회귀를 실행한다. 한국어 resource key/placeholder, 사용자 원문, 명시적 deny, root identity, 동시성, 출력 예산, 추가 메시지 전달, stale/last-good, metadata·migration·부분 실패·rollback을 관련 경계에서 검증한다.

native exec/session/task는 격리된 실제 실행·관찰·취소·exit code·output 및 결과 재조회가 재실행하지 않는지 검사한다. COMMON240로 대체하거나 새 240단계 캠페인을 만들지 않는다. 정상 제품/운영 설치에 Setup을 실행하지 않는다. 기존 결과는 동일 소스·의존성·조건일 때만 재사용한다.

현재 실제 상태: 백업·독립 복원·실험 정리 PASS; 공개 기준선 fetch 완료; 제품 merge/compile/회귀/main 통합 not-run. 현재 설치·게시 not-run. 과거 1.1.16102 PASS 기록은 Git/보관본에 남기고 새 조합의 PASS로 복사하지 않는다.

모든 필수 검증 후 main의 현재 상태를 다시 확인하고 정상 merge한다. main이 바뀌었으면 덮어쓰지 않고 변경 조합을 검토한다. 정본·필요한 제품 코드/시험/자산/라이선스만 main에 남기고 원시 실험 자료·임시 빌드·설치파일은 외부에 둔다.

## 6. 현재 재개 지점

외부 integration worktree가 LOCAL_BASE에 만들어졌다. 다음 행동은 고정 WB_BASE 정상 병합과 실제 충돌의 upstream 우선 해결이다. main은 아직 LOCAL_BASE이며 배포 버전 결정 전 패키징을 하지 않는다. 실제 완료 경계와 남은 검사만 이 정본에서 갱신한다.
