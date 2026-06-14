-- CYOps RAG metadata schema skeleton.
-- pgvector extension is expected before embeddings are enabled:
-- CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS cyops_documents (
    id UUID PRIMARY KEY,
    namespace TEXT NOT NULL,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT '',
    size_bytes BIGINT NOT NULL DEFAULT 0,
    object_uri TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    embedding_status TEXT NOT NULL,
    uploaded_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    last_error TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS cyops_document_chunks (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES cyops_documents(id),
    chunk_index INTEGER NOT NULL,
    text TEXT NOT NULL,
    token_count INTEGER NOT NULL DEFAULT 0,
    source_start INTEGER NOT NULL DEFAULT 0,
    source_end INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (document_id, chunk_index)
);

CREATE TABLE IF NOT EXISTS cyops_document_embeddings (
    chunk_id UUID PRIMARY KEY REFERENCES cyops_document_chunks(id),
    model TEXT NOT NULL,
    dimensions INTEGER NOT NULL,
    -- Replace BYTEA with pgvector VECTOR(dimensions) when the target
    -- embedding model dimensions are fixed in v0.0.6.
    embedding BYTEA NOT NULL DEFAULT ''::bytea,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cyops_chat_sessions (
    id UUID PRIMARY KEY,
    namespace TEXT NOT NULL,
    user_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cyops_chat_messages (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES cyops_chat_sessions(id),
    role TEXT NOT NULL,
    provider TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    citations_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS cyops_documents_namespace_status_idx
    ON cyops_documents (namespace, status);

CREATE INDEX IF NOT EXISTS cyops_document_chunks_document_idx
    ON cyops_document_chunks (document_id, chunk_index);

CREATE INDEX IF NOT EXISTS cyops_chat_messages_session_idx
    ON cyops_chat_messages (session_id, created_at);
