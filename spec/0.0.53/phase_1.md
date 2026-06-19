# v0.0.53 Phase 1 - Chat Drawer Cleanup

## Tasks

- [x] Remove document upload/list UI from the chat drawer.
- [x] Keep the chat drawer focused on messages and input.
- [x] Add `GET /api/chat/message/<encoded-message>` for proxy-safe chat.
- [x] Update Console plugin chat calls to use the path-based fallback.

## Verification

- `go test ./...` passed.

## Remaining Scope

- Package and verify on CRC.
