# v0.0.40 Planner

## 1. 목표

- Make the CYOps launcher visibly distinct from the Red Hat OpenShift Lightspeed button.
- Avoid bottom-right overlap by placing CYOps to the left of the Lightspeed launcher.
- Add inline launcher styles so CYOps remains visible even if Console CSS ordering changes.

## 2. 아키텍처 개요

- The appserver still serves the dynamic plugin manifest and callback entry bundle.
- The entry bundle injects a `CYOps` launcher button with class styles and inline fallback styles.
- The drawer behavior and backend API paths remain unchanged.

## 3. 기술 스택

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | OpenShift callback dynamic plugin | Existing v0.0.39 manifest shape |
| Launcher UI | Vanilla JS + CSS | Inline fallback styles for visibility |
| Packaging | OLM bundle/catalog | v0.0.40 replaces v0.0.39 |

## 4. 구현 단계

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Launcher visibility and collision fix | done | appserver tests |
| Phase 2 | Packaging and CRC smoke | pending | OLM and endpoint evidence |
| Phase 3 | Issue/PR handoff | pending | GitHub evidence |

Tracking issue: #190

## 5. 마이그레이션 또는 운영 전략

- Upgrade from `cywell-opsmate.v0.0.39` to `cywell-opsmate.v0.0.40`.
- Keep the existing `cyops-console` plugin name and appserver Service unchanged.

## 6. 메시지/통신/데이터 프로토콜

- No API contract changes.
- The launcher still calls `/api/chat` and `/api/documents` through the Console plugin backend path.

## 7. 보안 고려사항

- No credential or OAuth changes.
- Red Hat Lightspeed remains enabled and untouched.

## 8. 완료 기준

- [x] Entry bundle renders a `CYOps` launcher label.
- [x] Launcher is positioned left of the Lightspeed button.
- [x] Launcher has inline fallback styles.
- [x] Go tests pass.
- [ ] v0.0.40 packaging validates and installs on CRC.
