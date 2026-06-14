# Phase 4 - OpsMateConfig API 초안

## 작업 내용

- [x] `api/v1alpha1/OpsMateConfig` 타입 초안을 작성했다.
- [x] Lightspeed API 설정 경계를 `LightspeedSpec`으로 분리했다.
- [x] PostgreSQL DB 방향을 `DatabaseSpec`에 반영했다.
- [x] Console Plugin 설정 경계를 `ConsoleSpec`으로 분리했다.
- [x] AIOps와 RAG 설정 경계를 별도 spec으로 분리했다.

## 검증

- [x] API 타입이 외부 dependency 없이 작성되었는지 파일 수준으로 확인했다.
- [ ] `go test ./...`로 전체 Go package 컴파일을 확인한다. 현재 로컬 PATH에 Go가 없어 다음 환경 검증으로 이관한다.

## 남은 범위

- [ ] Kubebuilder marker, deepcopy generation, CRD YAML 생성은 다음 버전에서 진행한다.
- [ ] OpenShift/Kubernetes 타입 import는 실제 controller-runtime 도입 버전에서 추가한다.
