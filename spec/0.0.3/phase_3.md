# Phase 3 - appserver TLS/service-ca 경계

## 작업 내용

- [ ] appserver Service에 OpenShift service-ca serving cert annotation을 추가한다.
- [ ] appserver Deployment에 serving cert Secret mount/env 경계를 추가한다.
- [ ] 단위 테스트로 TLS Secret 이름과 mount path를 고정한다.

## 검증

- [ ] `go test ./internal/controller/appserver`
- [ ] `go test ./...`

## 남은 범위

- [ ] 실제 TLS listener 적용은 backend skeleton 또는 후속 버전에서 구현한다.
