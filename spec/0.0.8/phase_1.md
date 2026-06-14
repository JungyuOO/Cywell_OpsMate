# Phase 1 - parser and chunker boundary

## 작업 항목

- [x] parser interface를 추가한다.
- [x] plain text/markdown parser를 구현한다.
- [x] deterministic chunker를 구현한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- chunk persistence는 Phase 2에서 구현했다.

## 작업 내용

- `DocumentParser`, `TextFileParser`, `FixedRuneChunker`를 추가했다.
- plain text/markdown 파일을 dependency 없이 읽고 deterministic rune chunks로 분할한다.
