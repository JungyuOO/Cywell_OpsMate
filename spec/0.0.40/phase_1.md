# v0.0.40 Phase 1 - Launcher Visibility

## 작업 항목

- [x] Move launcher from `right: 22px` to `right: 92px`.
- [x] Change visible launcher text from `CY` to `CYOps`.
- [x] Add inline fallback styles to the launcher element.
- [x] Add regression checks for the visible label and offset.

## 검증

- `go test ./...` passed.
- Local appserver entry smoke returned:
  - `hasCYOpsLabel: true`
  - `hasInlineRight: true`
  - `hasImportantOffset: true`
  - `hasMaxZIndex: true`

## 남은 범위

- Packaging and CRC smoke remain for Phase 2.
