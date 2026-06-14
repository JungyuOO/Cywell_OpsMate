# Phase 3 - observability and auth operations

## 작업 항목

- [x] retrieval observer interface를 추가한다.
- [x] slow retrieval threshold와 failure reason boundary를 추가한다.
- [x] embedding token rotation 운영 조건을 문서화한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- metrics backend 연동은 후속 버전에서 진행한다.

## 작업 내용

- `RetrievalObserver`와 `RetrievalObservation`을 추가했다.
- observation은 mode, duration, slow 여부, failure reason, result count만 포함하고 query/chunk text는 포함하지 않는다.
- token rotation은 Secret update 후 appserver rollout restart를 기본 절차로 문서화했다.
