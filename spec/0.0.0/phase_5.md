# Phase 5 - v0.0.0 검증 및 후속 버전 이관

## 작업 내용

- [x] `planner.md`를 5개 phase 기준으로 구체화했다.
- [x] `phase_1.md`부터 `phase_5.md`까지 작성했다.
- [x] `v0.0.0` 범위를 프로젝트 구조 초안으로 제한했다.
- [x] 실제 Operator reconcile, ConsolePlugin, PostgreSQL 배포, Lightspeed API client, AIOps/RAG 구현은 다음 버전으로 이관했다.

## 검증

- [x] `go fmt ./...`
- [x] `go test ./...`
- [x] `go build -o .cache/manager.exe ./cmd/manager`
- [x] `git diff --cached` 검토

## 남은 범위

- [x] Go toolchain이 설치된 환경에서 `go fmt ./...`, `go test ./...`, `go build -o .cache/manager.exe ./cmd/manager`를 수행한다.
- [ ] 다음 버전 폴더를 생성한다.
- [ ] 다음 버전 브랜치로 전환한다.
- [ ] controller-runtime/Kubebuilder dependency를 도입한다.
- [ ] 실제 CRD, RBAC, manager, reconciler 구현을 시작한다.
