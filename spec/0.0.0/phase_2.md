# Phase 2 - Lightspeed Operator 전수 구조 분석

## 작업 내용

- [x] `openshift/lightspeed-operator` 공개 저장소를 로컬 참고자료로 clone했다.
- [x] 최상위 파일, API, controller, component package, config, bundle, test, hack, catalog 폴더를 확인했다.
- [x] 원본이 Go 1.25.9와 Kubebuilder v4 layout을 사용한다는 점을 확인했다.
- [x] 원본이 `OLSConfig` CRD와 `appserver`, `postgres`, `console` 컴포넌트별 reconcile 패키지를 사용한다는 점을 확인했다.
- [x] conversation cache DB가 PostgreSQL이며 Operator가 PostgreSQL 리소스를 관리한다는 점을 확인했다.
- [x] 검토 결과를 `lightspeed_operator_review.md`에 기록했다.

## 검증

- [x] `rg --files`로 원본 저장소 파일 목록을 확인했다.
- [x] `ARCHITECTURE.md`, `PROJECT`, `go.mod`, `api/v1alpha1/olsconfig_types.go`, `cmd/main.go`를 확인했다.
- [x] `internal/controller/postgres`, `console`, `appserver` 패키지를 확인했다.

## 남은 범위

- [ ] 원본 코드 복사는 하지 않는다.
- [ ] controller-runtime dependency 추가와 실제 reconcile 구현은 다음 버전에서 진행한다.
