# v0.0.38 Phase 2 - CYOps Launcher And Chat UI

## 작업 항목

- [x] Inject a bottom-right CYOps launcher into OpenShift Web Console.
- [x] Open a chat drawer branded as CYOps.
- [x] Call `/api/chat` through the ConsolePlugin backend path.
- [x] List and upload customer documents through `/api/documents`.

## 검증

- `go test ./...` passed.
- Local HTTP smoke confirmed `/api/chat` returned provider `lightspeed` through the appserver.
- Playwright smoke opened `/console-plugin/diagnostics`, injected `/plugin-entry.js`, clicked the CYOps launcher, and confirmed:
  - `launcher: true`
  - `drawerOpen: true`
  - `title: CYOps`
  - `textarea: true`
  - `upload: true`

## 남은 범위

- CRC browser smoke and final OLM upgrade validation remain for Phase 3.
