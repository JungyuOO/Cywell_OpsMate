# v0.0.25 Phase 2 - Local Diagnostics Smoke

## Work

- [x] Added `deploy/scripts/local-v025-diagnostics-smoke.ps1`.
- [x] Checked diagnostics HTML, JS, API, and schema.
- [x] Checked the console diagnostics JS does not contain OAuth handling.
- [x] Added loopback-only `CYOPS_DEV_ADMIN_USER` support for direct browser verification.

## Verification

- PowerShell smoke execution against `cmd/appserver`.
- Unit tests for loopback and non-loopback dev user behavior.

## Remaining Scope

- Extend smoke to PostgreSQL-backed diagnostics once live cluster evidence is collected.

## Linked Issue

- #121
