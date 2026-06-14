# Phase 2 - manager 및 reconciler 등록

## 작업 내용

- [x] `cmd/manager`를 controller-runtime manager entrypoint로 전환한다.
- [x] `OpsMateConfig` reconciler를 manager에 등록한다.
- [x] reconcile 함수는 안전한 no-op 또는 기본 status 갱신 경계까지만 구현한다.

## 검증

- [x] `go fmt ./...`
- [x] `go test ./...`
- [x] `go build -o .cache/manager.exe ./cmd/manager`

## 남은 범위

- [ ] 실제 Kubernetes 리소스 생성은 후속 phase 또는 후속 버전에서 구현한다.
