# CYOps Backend API Contract

## 공통 규칙

- 모든 API는 appserver가 제공한다.
- ConsolePlugin frontend는 appserver 같은 namespace의 service endpoint를 호출한다.
- 인증/인가 상세는 OpenShift console proxy 또는 service account token 연동 버전에서 확정한다.
- 응답은 JSON을 기본으로 한다.
- Secret 값과 문서 원문 전문은 로그와 오류 응답에 포함하지 않는다.

## POST /api/chat

사용자 질문을 CYOps backend에 전달한다.

### Request

```json
{
  "sessionId": "optional-session-id",
  "message": "Why is this pod not ready?",
  "provider": "lightspeed",
  "rag": {
    "enabled": true,
    "documentIds": ["doc-001"],
    "scope": "selected"
  },
  "clusterContext": {
    "namespace": "openshift-marketplace",
    "resourceRef": {
      "apiVersion": "v1",
      "kind": "Pod",
      "name": "example"
    }
  }
}
```

### Response

```json
{
  "sessionId": "session-001",
  "messageId": "msg-002",
  "answer": "The pod is not ready because...",
  "provider": "lightspeed",
  "citations": [
    {
      "documentId": "doc-001",
      "title": "customer-runbook.pdf",
      "chunkId": "chunk-009"
    }
  ],
  "warnings": [
    "Review AI generated content before use."
  ]
}
```

## GET /api/documents

고객 문서 목록을 조회한다.

### Query

- `namespace`: optional namespace/tenant scope
- `status`: optional `uploaded|processing|ready|failed`

### Response

```json
{
  "items": [
    {
      "id": "doc-001",
      "filename": "customer-runbook.pdf",
      "status": "ready",
      "sizeBytes": 120000,
      "chunkCount": 42,
      "embeddingStatus": "ready",
      "uploadedBy": "admin",
      "createdAt": "2026-06-14T08:00:00Z"
    }
  ]
}
```

## POST /api/documents

고객 문서를 업로드한다.

### Request

- `multipart/form-data`
- field `file`: uploaded document
- field `scope`: optional namespace/tenant scope

### Response

```json
{
  "id": "doc-001",
  "filename": "customer-runbook.pdf",
  "status": "uploaded"
}
```

## DELETE /api/documents/{documentId}

문서 metadata, 원본, chunk, embedding을 삭제 대상으로 표시한다.

### Response

```json
{
  "id": "doc-001",
  "status": "deleting"
}
```

## GET /api/documents/{documentId}

문서 상세와 ingestion 상태를 조회한다.

### Response

```json
{
  "id": "doc-001",
  "filename": "customer-runbook.pdf",
  "status": "ready",
  "chunkCount": 42,
  "embeddingStatus": "ready",
  "lastError": ""
}
```

## POST /api/providers/lightspeed/chat

내부 provider 경계다. ConsolePlugin frontend가 직접 호출하지 않고 `/api/chat`이 provider policy에 따라 호출한다.

### Request

```json
{
  "message": "Why is this pod not ready?",
  "context": [
    {
      "type": "rag_chunk",
      "text": "Minimal relevant text only",
      "source": "doc-001/chunk-009"
    }
  ],
  "clusterContext": {}
}
```

### Response

```json
{
  "answer": "The pod is not ready because...",
  "rawProvider": "lightspeed"
}
```

## UI mapping

- Chat drawer input calls `POST /api/chat`.
- `+` menu calls `POST /api/documents` for upload.
- Document side panel calls `GET /api/documents`.
- Delete action calls `DELETE /api/documents/{documentId}`.
- Document status badges are driven by `status` and `embeddingStatus`.
