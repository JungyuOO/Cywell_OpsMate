# v0.0.7 Planner

## 1. 목표

- appserver가 환경변수 기반으로 PostgreSQL DSN, namespace, 문서 저장 경로, Lightspeed endpoint를 구성할 수 있게 한다.
- `/api/documents` 업로드가 원본 파일을 `DocumentStorage`에 저장하고, PostgreSQL metadata에 `object_uri`와 size를 함께 기록할 수 있게 한다.
- migration 실행 전략을 appserver 코드 경계로 명시하고 idempotent하게 검증한다.
- ingestion/parser/embedding 구현은 다음 버전으로 넘기되, v0.0.8의 우선순위를 문서화한다.

## 2. 아키텍처 개요

1. appserver는 `LoadConfigFromEnv`로 runtime 설정을 읽는다.
2. PostgreSQL DSN이 있으면 `pgx` stdlib driver로 DB를 열고 migration을 적용한다.
3. 문서 저장 경로가 있으면 `LocalDocumentStorage`가 PVC/local path에 원본 파일을 저장한다.
4. 업로드 handler는 저장된 object URI와 size를 repository에 전달한다.
5. Lightspeed endpoint가 있으면 `LightspeedProvider`가 외부 REST API를 호출하고, 없으면 skeleton 응답을 유지한다.

## 3. 기술 스택

| 영역 | 선택 | 이유 |
| --- | --- | --- |
| Runtime config | environment variables | Operator Deployment env와 직접 연결 가능 |
| PostgreSQL driver | `github.com/jackc/pgx/v5/stdlib` | v0.0.6 integration test와 일관 |
| Migration | embedded SQL + `database/sql` | 별도 migration dependency 없이 idempotent SQL 실행 |
| Storage | local/PVC filesystem | OpenShift PVC에 바로 매핑 가능한 최소 구현 |
| HTTP API | Go `net/http` | 기존 appserver skeleton 유지 |

## 4. 구현 단계

| Phase | 범위 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | appserver runtime config | 완료 | env loader, dependency builder tests |
| Phase 2 | upload storage + metadata persistence | 완료 | stored upload path, handler tests |
| Phase 3 | runtime migration strategy | 완료 | migration runner, integration test 적용 |
| Phase 4 | v0.0.8 ingestion/embedding handoff | 완료 | parser/chunker/embedding scope |

## 5. 마이그레이션 또는 운영 전략

- appserver startup path에서 PostgreSQL DSN이 설정된 경우 migration SQL을 적용한다.
- migration SQL은 `CREATE TABLE IF NOT EXISTS`와 `CREATE INDEX IF NOT EXISTS`만 사용해 반복 실행 가능해야 한다.
- migration 실패 시 appserver dependency construction이 error를 반환하고 서버를 시작하지 않는다.
- 운영 배포에서 문서 저장 경로는 PVC mount path로 전달한다.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 데이터 | 상태 |
| --- | --- | --- |
| env -> appserver | `CYOPS_POSTGRES_DSN`, `CYOPS_NAMESPACE`, `CYOPS_DOCUMENT_STORAGE_PATH`, `CYOPS_LIGHTSPEED_ENDPOINT` | 구현 예정 |
| `POST /api/documents` -> storage | multipart file stream | 구현 예정 |
| storage -> repository | `object_uri`, `size_bytes`, filename, user | 구현 예정 |
| repository -> PostgreSQL | `cyops_documents` row | 구현 예정 |

## 7. 보안 고려사항

- 업로드 파일명은 path traversal 방지를 위해 저장 시 `filepath.Base`로 축약한다.
- 원본 문서는 `0600`, 문서별 디렉터리는 `0700` 권한으로 저장한다.
- DSN과 provider credentials는 로그나 응답에 노출하지 않는다.
- namespace는 repository scope로 사용해 문서 목록과 조회를 격리한다.

## 8. 완료 기준

- [x] appserver runtime config loader와 dependency builder가 있다.
- [x] 업로드 handler가 storage write와 metadata create를 연결한다.
- [x] PostgreSQL integration test가 repository + storage upload flow를 검증한다.
- [x] migration runner가 idempotent하게 SQL을 적용한다.
- [x] v0.0.8 ingestion/embedding 구현 범위가 handoff 문서로 정리된다.
