# v0.0.47 Phase 1 - Header And Response Adapter

## Tasks

- [x] Add `X-CSRF-Token` and `X-Requested-With` to non-GET drawer fetches.
- [x] Keep `X-CSRFToken` from v0.0.46.
- [x] Parse OpenAI-style `choices[].message.content` responses.
- [x] Parse internal LLM `response` field responses.
- [x] Preserve existing `answer` field support.
- [x] Update regression tests.

## Verification

- `go test ./...` passed.
- Tests verify OpenAI-style `choices[].message.content` response parsing.
- Tests verify internal LLM `response` field parsing.
- Tests verify the Console plugin entry includes `X-CSRFToken`, `X-CSRF-Token`, and `X-Requested-With`.

## Remaining Scope

- OLM packaging, image publication, CRC upgrade, and browser chat smoke remain for later phases.
