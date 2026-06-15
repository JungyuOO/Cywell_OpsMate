# v0.0.22 Phase 1 - OpenShift Admin Auth Smoke

## 작업 내용

- [x] Added `deploy/scripts/openshift-v022-smoke.ps1`.
- [x] Added `oc whoami` and namespace checks.
- [x] Added OAuth cookie Secret create/apply step.
- [x] Added admin Route OAuth redirect check.

## 검증

- PowerShell parser validation for `deploy/scripts/openshift-v022-smoke.ps1`.

## 남은 범위

- Live browser login through the generated Route should be run on the target OpenShift cluster.

## 연결 이슈

- #105
