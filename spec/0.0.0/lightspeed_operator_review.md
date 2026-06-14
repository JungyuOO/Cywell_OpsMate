# OpenShift Lightspeed Operator 검토 기록

검토 대상: https://github.com/openshift/lightspeed-operator

## 전수 확인 범위

로컬 참고 경로 `.omx/external/lightspeed-operator`에 공개 저장소를 clone한 뒤 다음 범위를 확인했다.

- 최상위 프로젝트 파일: `go.mod`, `PROJECT`, `Makefile`, `Dockerfile`, `bundle.Dockerfile`, `README.md`, `ARCHITECTURE.md`
- API 정의: `api/v1alpha1/olsconfig_types.go`, generated deepcopy 파일
- 컨트롤러 진입점: `cmd/main.go`
- 메인 컨트롤러: `internal/controller/olsconfig_controller.go`, helpers, finalizer, operator assets
- 컴포넌트 컨트롤러: `internal/controller/appserver`, `console`, `postgres`
- 공통 계층: `internal/controller/reconciler`, `utils`, `watchers`, `internal/tls`, `internal/relatedimages`
- 배포 구성: `config/crd`, `config/default`, `config/manager`, `config/rbac`, `config/manifests`, `config/prometheus`, `config/samples`, `config/scorecard`, `config/user-access`
- OLM 산출물: `bundle/manifests`, `bundle/metadata`, `bundle/tests/scorecard`
- 테스트: `test/e2e`, `test/utils`, 컴포넌트별 unit test
- 운영/릴리스 보조 파일: `hack`, `.tekton`, `konflux-integration`, `lightspeed-catalog*`, `related_images.json`

## 적용할 구조 결정

- 우리 프로젝트도 Go 기반 Operator로 구성한다.
- Kubebuilder v4 / Operator SDK 형식의 `PROJECT`, `api`, `cmd`, `internal/controller`, `config`, `bundle` 구조를 따른다.
- CRD 이름은 `OpsMateConfig`, API group은 `opsmate.cywell.io`, 버전은 `v1alpha1`로 시작한다.
- 원본의 `OLSConfig`처럼 cluster singleton 성격을 유지하며 기본 리소스 이름은 `cluster`로 둔다.
- Lightspeed Operator가 conversation cache로 PostgreSQL을 관리하므로, 우리 프로젝트도 기본 DB를 PostgreSQL로 둔다.
- PostgreSQL은 후속 버전에서 Deployment, Service, PVC, Secret, ConfigMap, NetworkPolicy로 Operator가 reconcile한다.
- Console Plugin은 후속 버전에서 OpenShift `ConsolePlugin` CR과 console operator plugin 활성화 흐름으로 구현한다.
- AIOps와 RAG는 Lightspeed API 연동과 분리된 패키지 경계로 둔다.

## v0.0.0에서 구현한 범위

- Go module 초안
- Kubebuilder/Operator SDK 구조를 나타내는 `PROJECT`
- `cmd/manager` 진입점 초안
- `api/v1alpha1/OpsMateConfig` 타입 초안
- `internal/controller` 하위 컴포넌트 패키지 경계
- `config`와 `bundle` 기본 디렉터리

## 다음 버전으로 이관할 범위

- 실제 Kubebuilder markers와 CRD 생성
- controller-runtime manager 연결
- Kubernetes/OpenShift API dependency 추가
- PostgreSQL 리소스 생성 로직
- Lightspeed API client 구현
- OpenShift ConsolePlugin 구현
- AIOps 분석 로직
- RAG index image/BYOK 연동
- unit/e2e/scorecard 테스트
