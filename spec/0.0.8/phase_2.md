# Phase 2 - PostgreSQL chunk repository

## 작업 항목

- [x] document status를 processing/ready/failed로 갱신하는 repository method를 추가한다.
- [x] document chunks를 교체 저장하는 repository method를 추가한다.
- [x] chunk 조회 또는 검증용 method를 추가한다.

## 검증

- [x] Docker PostgreSQL integration test

## 남은 범위

- embedding persistence는 v0.0.9로 넘긴다.

## 작업 내용

- `BeginIngestion`, `CompleteIngestion`, `FailIngestion`, `ReplaceChunks`, `ListChunksContext`를 추가했다.
- document 조회 시 chunk table count를 `chunkCount`로 반영한다.
