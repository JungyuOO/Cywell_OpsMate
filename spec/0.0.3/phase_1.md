# Phase 1 - ConsolePlugin desired resource helper

## 작업 내용

- [ ] `internal/controller/console`에 ConsolePlugin desired object builder를 추가한다.
- [ ] appserver Service 이름과 namespace를 ConsolePlugin backend service에 연결한다.
- [ ] `ConsoleSpec.Enabled`가 false일 때 reconcile에서 건너뛸 수 있는 경계를 둔다.
- [ ] desired object shape를 단위 테스트로 고정한다.

## 검증

- [ ] `go test ./internal/controller/console`
- [ ] `go test ./...`

## 남은 범위

- [ ] reconciler create/update 연결은 Phase 2에서 구현한다.
- [ ] 실제 plugin frontend asset은 후속 버전으로 이관한다.
