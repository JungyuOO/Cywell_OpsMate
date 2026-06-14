# Phase 3 - CRD/RBAC/manager Kustomize 기본 매니페스트

## 작업 내용

- [x] `config/crd`에 `OpsMateConfig` CRD 산출물을 둔다.
- [x] `config/rbac`에 manager와 API 접근 권한을 정의한다.
- [x] `config/manager`, `config/default`, `config/samples`의 최소 Kustomize 구조를 작성한다.

## 검증

- [x] `kubectl kustomize config/default`
- [x] `kubectl kustomize config/samples`
- [x] `go test ./...`

## 남은 범위

- [x] Operator bundle/catalog 정리는 후속 버전으로 이관한다.
