# v0.0.41 Planner

## 1. 목표

- Place the CYOps launcher over the Red Hat OpenShift Lightspeed launcher in the same bottom-right slot.
- Use z-index layering instead of moving CYOps left or vertically above the Lightspeed icon.
- Keep inline styles so the launcher remains visible in the Console runtime.

## 2. 아키텍처 개요

- The appserver-served callback entry bundle still injects the CYOps launcher.
- CYOps launcher uses `right: 22px` and `bottom: 22px`.
- CYOps launcher uses the maximum z-index layer so it visually replaces the Lightspeed launcher.
- CYOps drawer opens above the launcher with `bottom: 86px`.

## 3. 기술 스택

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | OpenShift callback dynamic plugin | Existing v0.0.39 manifest shape |
| Launcher UI | Vanilla JS + CSS | Same right-side stack above Lightspeed |
| Packaging | OLM bundle/catalog | v0.0.41 replaces v0.0.40 |

## 4. 구현 단계

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Layer launcher over Lightspeed | done | appserver tests |
| Phase 2 | Packaging and CRC smoke | done | CRC CSV/appserver endpoint evidence |
| Phase 3 | Issue/PR handoff | done | PR #193 merged and issue #192 closed |

Tracking issue: #192

## 5. 마이그레이션 또는 운영 전략

- Upgrade from `cywell-opsmate.v0.0.40` to `cywell-opsmate.v0.0.41`.
- Keep `cyops-console`, appserver Service, and API paths unchanged.

## 6. 메시지/통신/데이터 프로토콜

- No API contract changes.

## 7. 보안 고려사항

- No credential or OAuth changes.
- Red Hat Lightspeed remains enabled and untouched.

## 8. 완료 기준

- [x] Launcher uses `right: 22px`.
- [x] Launcher uses `bottom: 22px`.
- [x] Launcher uses max z-index above Lightspeed.
- [x] Go tests pass.
- [x] v0.0.41 packaging validates and installs on CRC.
