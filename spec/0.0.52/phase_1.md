# v0.0.52 Phase 1 - Chat Input And Proxy-Safe Chat

## Tasks

- [x] Add GET-compatible `/api/chat` handling.
- [x] Keep existing POST `/api/chat` behavior.
- [x] Make Enter submit the drawer textarea and Shift+Enter preserve multiline input.
- [x] Use GET chat from the Console plugin proxy path.

## Verification

- `go test ./...` passed.

## Remaining Scope

- Package and verify on CRC.
