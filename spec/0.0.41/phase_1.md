# v0.0.41 Phase 1 - Launcher Over Lightspeed

## 작업 항목

- [x] Move CYOps launcher from left-of-Lightspeed to the same bottom-right slot as Lightspeed.
- [x] Keep visible `CYOps` label.
- [x] Keep max z-index so CYOps renders over the Lightspeed icon.
- [x] Update regression checks for `right: 22px`, `bottom: 22px`, and max z-index.

## 검증

- `go test ./...` passed.
- Local appserver entry smoke returned:
  - `hasCYOpsLabel: true`
  - `hasInlineRight: true`
  - `hasInlineBottom: true`
  - `hasImportantBottom: true`
  - `hasMaxZIndex: true`

## 남은 범위

- Packaging and CRC smoke remain for Phase 2.
