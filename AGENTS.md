# Cywell_OpsMate Agent Rules

이 저장소는 OpenShift Web Console에서 사용할 수 있는 Operator를 구축한다.
Operator는 Lightspeed API를 연동할 수 있어야 하며, 별도의 AIOps 기능과 필요 시 별도 RAG 기능을 확장할 수 있어야 한다.

## 프로젝트 범위

- OpenShift Web Console에서 동작하는 Operator를 기준으로 개발한다.
- Lightspeed API 연동은 외부 API를 받아 사용하는 구조로 설계한다.
- AIOps 기능은 Operator 기능과 분리 가능한 내부 기능 단위로 구현한다.
- RAG 기능은 필요한 버전에서 별도 범위로 명시된 경우에만 구현한다.

## 버전별 작업 규칙

- 모든 기능 구현은 `spec/<version>/` 하위 문서를 기준으로만 수행한다.
- 예: `spec/0.0.0/`, `spec/0.0.1/`, `spec/0.1.0/`.
- 현재 버전의 `planner.md`에 적힌 기능만 구현한다.
- 현재 버전 범위가 완료되기 전에는 다음 버전 기능을 구현하지 않는다.
- 다음 버전 작업을 시작할 때는 새 `spec/<next-version>/` 폴더를 만들고, 해당 버전에 맞는 새 브랜치로 변경한 뒤 작업한다.

## 필수 문서

각 버전 폴더에는 최소한 다음 문서를 둔다.

- `planner.md`: 해당 버전의 목표, 아키텍처, 기술 스택, 구현 단계, 보안 고려사항, 완료 기준을 기록한다.
- `phase_1.md`, `phase_2.md` 등: phase별 작업 내용, 검증 방법, 남은 범위를 기록한다.

`planner.md`는 다음 구조를 따른다.

1. 목표
2. 아키텍처 개요
3. 기술 스택
4. 구현 단계
5. 마이그레이션 또는 운영 전략
6. 메시지/통신/데이터 프로토콜
7. 보안 고려사항
8. 완료 기준

## Phase 완료 규칙

각 phase가 완료되면 반드시 다음을 수행한다.

- `planner.md`의 해당 phase 상태를 체크 표시로 갱신한다.
- 해당 `phase_N.md`의 완료 항목을 체크 표시로 갱신한다.
- 작업 내용, 검증, 남은 범위를 문서에 남긴다.
- 관련 GitHub Issue를 생성하거나 기존 Issue를 연결한다.
- 커밋 메시지에 작업 의도와 검증 결과를 기록한다.
- 브랜치를 원격에 푸시한다.
- PR merge 또는 commit 연결을 통해 GitHub Issue가 이력과 함께 Closed 되도록 한다.

## GitHub Issue 기록 형식

Issue 본문은 아래 형식을 따른다.

```markdown
## 작업 내용
- 수행한 작업을 기록합니다.

## 검증
- 수행한 검증을 기록합니다.

## 남은 범위
- 후속 버전 또는 후속 phase로 넘기는 범위를 기록합니다.
```

## 커밋 규칙

커밋 메시지는 저장소의 Lore Commit Protocol을 따른다.
첫 줄은 변경 내용이 아니라 변경한 이유를 설명한다.
검증한 항목과 검증하지 못한 항목을 `Tested:` / `Not-tested:` trailer로 남긴다.

## 검증 규칙

- 프론트엔드 변경 시 `frontend check`와 build에 해당하는 저장소 명령을 수행한다.
- Operator, API, AIOps, RAG 기능 변경 시 가능한 단위 테스트, 통합 테스트, 빌드 검증을 수행한다.
- 검증하지 못한 항목은 완료 보고와 커밋 메시지에 명시한다.
