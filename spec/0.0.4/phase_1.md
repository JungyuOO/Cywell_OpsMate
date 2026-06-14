# Phase 1 - CYOps UX/API 아키텍처 결정

## 작업 내용

- [x] CYOps UI가 Red Hat Lightspeed Viewer를 그대로 쓰지 않는다는 결정을 기록한다.
- [x] Lightspeed Operator/API는 provider로만 사용한다는 경계를 기록한다.
- [x] CYOps floating launcher, chat drawer, 상단 `CYOps` brand, 문서 drawer 위치를 정의한다.
- [x] 고객 문서 업로드/목록/관리 UI가 chat drawer 안에서 동작한다는 UX 계약을 정의한다.

## 검증

- [x] `planner.md`와 architecture decision 문서 검토
- [x] v0.0.5 구현 단계로 분해 가능한지 확인

## 남은 범위

- [ ] 실제 ConsolePlugin frontend bundle 구현은 후속 버전에서 진행한다.
- [ ] 실제 Lightspeed API 호출은 후속 버전에서 구현한다.
