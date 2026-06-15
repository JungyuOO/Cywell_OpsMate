# v0.0.21 Phase 1 - OpenShift Admin Auth Proxy Resources

## 작업 내용

- [x] Added `internal/controller/authproxy`.
- [x] Added OAuth proxy ServiceAccount with OpenShift redirect annotation.
- [x] Added OAuth proxy Deployment that proxies to the appserver Service.
- [x] Added Service and Route that target only the auth proxy.

## 검증

- `go test ./internal/controller/authproxy`
- `go test ./internal/controller`

## 남은 범위

- Live OpenShift login flow smoke should be performed against a real cluster after image and cookie Secret values are configured.

## 연결 이슈

- #100
