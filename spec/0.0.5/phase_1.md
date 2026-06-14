# Phase 1 - Backend API DTO 및 handler skeleton

## 작업 내용

- [ ] `/api/chat` DTO와 handler를 추가한다.
- [ ] `/api/documents` list/upload handler를 추가한다.
- [ ] `/api/documents/{documentId}` detail/delete handler를 추가한다.
- [ ] mocked provider와 in-memory document repository로 handler tests를 고정한다.

## 검증

- [ ] `go test ./internal/appserver`
- [ ] `go test ./...`

## 남은 범위

- [ ] PostgreSQL repository 연결은 Phase 2 이후로 이관한다.
- [ ] 실제 Lightspeed API client는 Phase 3 이후로 이관한다.
