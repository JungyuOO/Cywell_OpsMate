# Phase 3 - Go Operator 프로젝트 구조 초안

## 작업 내용

- [x] Go module을 생성했다.
- [x] Kubebuilder/Operator SDK 구조 방향을 `PROJECT`에 기록했다.
- [x] `cmd/manager` 진입점을 추가했다.
- [x] `internal/controller` 하위에 appserver, postgres, console, aiops, rag, reconciler, utils 경계를 만들었다.
- [x] `config`와 `bundle` 기본 디렉터리 초안을 만들었다.

## 검증

- [x] `Makefile`에 `fmt`, `test`, `build` 명령을 정의했다.
- [ ] `go fmt ./...`는 로컬 Go 설치 후 수행한다.
- [ ] `go test ./...`는 로컬 Go 설치 후 수행한다.
- [ ] `go build -o bin/manager ./cmd/manager`는 로컬 Go 설치 후 수행한다.

## 남은 범위

- [ ] 실제 controller-runtime manager wiring은 다음 버전으로 이관한다.
- [ ] CRD generation과 RBAC generation은 다음 버전으로 이관한다.
- [ ] OpenShift ConsolePlugin 배포 리소스 구현은 다음 버전으로 이관한다.
