# Phase 1 - controller-runtime 의존성 및 API scheme 정리

## 작업 내용

- [x] controller-runtime 기반 의존성을 `go.mod`에 추가한다.
- [x] `api/v1alpha1`에 group/version과 scheme 등록 코드를 추가한다.
- [x] `OpsMateConfig` 타입에 CRD 생성을 위한 marker와 status/spec 구조를 보강한다.
- [x] 다음 phase에서 manager와 reconciler를 연결할 수 있도록 API surface를 고정한다.

## 검증

- [x] `go mod tidy`
- [x] `go fmt ./...`
- [x] `go test ./...`

## 남은 범위

- [x] manager 실행 구조는 Phase 2에서 구현한다.
- [ ] CRD/RBAC/Kustomize manifest는 Phase 3에서 작성한다.
- [ ] 실제 appserver/postgres/console 리소스 생성은 후속 버전으로 이관한다.
