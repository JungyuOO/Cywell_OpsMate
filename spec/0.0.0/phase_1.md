# Phase 1 - 프로젝트 규칙 및 버전 명세 체계 정리

## 작업 내용

- [x] OpenShift Web Console Operator 프로젝트의 개발 규칙을 루트 `AGENTS.md`에 기록했다.
- [x] 모든 기능 구현은 `spec/<version>/planner.md`에 적힌 범위만 수행하도록 규칙을 정리했다.
- [x] phase 완료 시 `planner.md`와 `phase_N.md`에 체크 표시를 남기도록 기준을 정리했다.
- [x] GitHub Issue, commit, PR merge 또는 commit 연결을 통해 Closed 이력이 남도록 운영 기준을 정리했다.
- [x] `spec/0.0.0/planner.md`를 목표, 아키텍처, 기술 스택, 구현 단계, 보안, 완료 기준 형식으로 정리했다.

## 검증

- [x] 로컬에서 `AGENTS.md`와 `spec/0.0.0` 문서 구조를 확인했다.
- [x] `planner.md`의 Phase 1 상태가 완료로 표시되는지 확인했다.

## 남은 범위

- [ ] Operator 기본 구조 설계는 Phase 2에서 진행한다.
- [ ] OpenShift Web Console Plugin 화면/메뉴 범위는 Phase 3에서 진행한다.
- [ ] Lightspeed API 연동 경계는 Phase 4에서 진행한다.
- [ ] AIOps/RAG 세부 구현 범위는 Phase 5 또는 후속 버전에서 분리한다.
