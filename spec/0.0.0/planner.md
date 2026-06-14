# OpenShift Web Console + AIOps Operator 구축 계획

## 1) 목표

- OpenShift Web Console에서 사용할 수 있는 Operator를 구축한다.
- OpenShift Web Console 내부에서 Lightspeed와 유사한 사용자 경험을 제공한다.
- Lightspeed API를 받아 사용할 수 있는 연동 구조를 마련한다.
- 별도의 AIOps 기능을 구축할 수 있는 확장 지점을 정의한다.
- 필요 시 별도의 RAG 기능을 구축할 수 있도록 버전별 범위를 분리한다.
- 모든 기능은 `spec/0.0.0`에 명시된 범위 안에서만 구현한다.

## 2) 아키텍처 개요

통신 흐름

1. OpenShift Web Console 사용자가 Operator UI를 통해 기능을 호출한다.
2. Operator backend가 Lightspeed API 또는 내부 AIOps 기능으로 요청을 전달한다.
3. 필요 시 RAG 계층이 사내 문서, 운영 지식, 장애 이력 데이터를 조회한다.
4. 결과는 Web Console에서 운영자가 확인할 수 있는 형태로 반환한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| Operator | Kubernetes Operator 구조 | OpenShift 배포와 운영 표준에 맞춤 |
| Web Console Extension | OpenShift Console Plugin | Web Console 내부 UX 제공 |
| Lightspeed 연동 | HTTP API client | 외부 Lightspeed API 호출 경계 분리 |
| AIOps 기능 | 내부 service module | Operator 기능과 분리해 버전별 확장 가능 |
| RAG | 후속 버전에서 확정 | `0.0.0`에서는 확장 지점만 정의 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | 프로젝트 규칙 및 버전 명세 체계 정리 | ✅ 완료 | `AGENTS.md`, `spec/0.0.0/planner.md`, `phase_1.md` |
| Phase 2 | Operator 기본 구조 설계 | ⬜ 예정 | Operator scaffold 설계 문서 |
| Phase 3 | Web Console Plugin 범위 정의 | ⬜ 예정 | Console UI 메뉴/화면 명세 |
| Phase 4 | Lightspeed API 연동 경계 정의 | ⬜ 예정 | API client 인터페이스 명세 |
| Phase 5 | AIOps/RAG 후속 범위 분리 | ⬜ 예정 | 후속 버전 후보 범위 |

## 5) 운영 전략

- `spec/0.0.0`의 기능만 구현한다.
- phase가 완료되면 `planner.md`와 해당 `phase_N.md`에 체크 표시를 남긴다.
- 다음 기능 범위는 새 버전 폴더를 생성한 뒤 새 브랜치에서 진행한다.
- GitHub Issue, commit, PR merge 이력을 통해 작업 완료와 Closed 기록을 남긴다.

## 6) 통신/데이터 프로토콜

`0.0.0`에서는 구현 프로토콜을 확정하지 않고 경계만 정의한다.

| 경로 | 용도 | 상태 |
| --- | --- | --- |
| Web Console Plugin → Operator backend | 사용자 요청 전달 | 후속 phase에서 정의 |
| Operator backend → Lightspeed API | Lightspeed API 호출 | 후속 phase에서 정의 |
| Operator backend → AIOps module | 내부 분석 기능 호출 | 후속 phase에서 정의 |
| Operator backend → RAG module | 문서/지식 조회 | 후속 버전에서 정의 |

## 7) 보안 고려사항

- Lightspeed API 인증 정보는 코드에 하드코딩하지 않는다.
- OpenShift Secret 또는 동등한 보안 저장소를 사용한다.
- 고객사별 데이터, 운영 로그, 장애 이력은 테넌트 경계를 분리한다.
- RAG 기능은 후속 버전에서 데이터 반출, 색인 권한, 감사 로그 기준을 별도로 정의한다.

## 8) 완료 기준

- [x] 루트 개발 규칙이 `AGENTS.md`에 기록된다.
- [x] `spec/0.0.0/planner.md`가 버전별 계획 형식으로 정리된다.
- [x] `spec/0.0.0/phase_1.md`가 작성된다.
- [ ] Operator 기본 구조 설계가 문서화된다.
- [ ] Web Console Plugin 범위가 문서화된다.
- [ ] Lightspeed API 연동 경계가 문서화된다.
- [ ] AIOps/RAG 후속 범위가 분리된다.

