# Phase 3 - Data model 및 RAG pipeline contract

## 작업 내용

- [ ] PostgreSQL metadata table 초안을 정의한다.
- [ ] pgvector 우선 전략과 별도 vector service 분리 조건을 정의한다.
- [ ] 문서 원본 저장소(PVC/object store) 경계를 정의한다.
- [ ] ingestion worker/server 필요성과 배포 경계를 정의한다.

## 검증

- [ ] 데이터 모델이 문서 목록/업로드/삭제/질의 흐름을 모두 설명하는지 확인
- [ ] BYOKnowledge 없이 동작 가능한지 확인

## 남은 범위

- [ ] migration SQL과 실제 ingestion 구현은 후속 버전에서 진행한다.
