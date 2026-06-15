package appserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresDocumentRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.now = func() time.Time { return time.Date(2026, 6, 14, 8, 0, 0, 0, time.UTC) }
	repository.newID = func() (string, error) { return "00000000-0000-4000-8000-000000000001", nil }

	created, err := repository.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   "runbook.pdf",
		SizeBytes:  42,
		ObjectURI:  "/var/lib/cyops/documents/runbook.pdf",
		UploadedBy: "admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Status != "uploaded" {
		t.Fatalf("created status = %q, want uploaded", created.Status)
	}
	if created.ObjectURI == "" {
		t.Fatal("created object uri is empty")
	}

	list, err := repository.ListContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list len = %d, want 1", len(list))
	}

	found, err := repository.GetContext(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.Filename != "runbook.pdf" {
		t.Fatalf("filename = %q, want runbook.pdf", found.Filename)
	}
	if found.ObjectURI != "/var/lib/cyops/documents/runbook.pdf" {
		t.Fatalf("object uri = %q", found.ObjectURI)
	}

	deleting, err := repository.MarkDeletingContext(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deleting.Status != "deleting" {
		t.Fatalf("deleting status = %q, want deleting", deleting.Status)
	}
}

func TestPostgresStorageUploadIntegration(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	server := NewServerWithOptions(ServerOptions{
		Provider:  MockProvider{},
		Documents: NewPostgresDocumentRepository(db, "opsmate"),
		Storage:   LocalDocumentStorage{BasePath: t.TempDir()},
	})

	upload := multipartRequest(t, "file", "../runbook.txt", "check pod status")
	upload.Header.Set("X-Forwarded-User", "admin")
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, upload)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("upload status = %d, want %d: %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	var filename string
	var objectURI string
	var sizeBytes int64
	if err := db.QueryRowContext(ctx, "SELECT filename, object_uri, size_bytes FROM cyops_documents WHERE namespace = $1", "opsmate").Scan(&filename, &objectURI, &sizeBytes); err != nil {
		t.Fatal(err)
	}
	if filename != "runbook.txt" {
		t.Fatalf("filename = %q, want multipart filename", filename)
	}
	if objectURI == "" || strings.Contains(objectURI, "..") {
		t.Fatalf("object uri = %q, want sanitized stored path", objectURI)
	}
	if sizeBytes != int64(len("check pod status")) {
		t.Fatalf("size = %d, want uploaded content size", sizeBytes)
	}
}

func TestPostgresIngestionIntegration(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	basePath := t.TempDir()
	stored, err := LocalDocumentStorage{BasePath: basePath}.Store(ctx, "00000000-0000-4000-8000-000000000001", "runbook.md", strings.NewReader("# Runbook\n\nCheck pod status before restart."))
	if err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.newID = fixedIDs(
		"00000000-0000-4000-8000-000000000001",
		"00000000-0000-4000-8000-000000000002",
		"00000000-0000-4000-8000-000000000003",
		"00000000-0000-4000-8000-000000000004",
	)
	document, err := repository.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   "runbook.md",
		SizeBytes:  stored.SizeBytes,
		ObjectURI:  stored.URI,
		UploadedBy: "admin",
	})
	if err != nil {
		t.Fatal(err)
	}

	ingested, err := IngestionService{
		Repository: repository,
		Chunker:    FixedRuneChunker{MaxRunes: 20},
	}.IngestDocument(ctx, document.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ingested.Status != "ready" {
		t.Fatalf("status = %q, want ready", ingested.Status)
	}
	if ingested.ChunkCount != 3 {
		t.Fatalf("chunk count = %d, want 3", ingested.ChunkCount)
	}

	chunks, err := repository.ListChunksContext(ctx, document.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 3 {
		t.Fatalf("chunks len = %d, want 3", len(chunks))
	}
	if chunks[0].Text != "# Runbook\n\nCheck pod" {
		t.Fatalf("chunk text = %q", chunks[0].Text)
	}
}

func TestPostgresEmbeddingIntegration(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.newID = fixedIDs(
		"00000000-0000-4000-8000-000000000011",
		"00000000-0000-4000-8000-000000000012",
	)
	document, err := repository.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   "runbook.md",
		SizeBytes:  20,
		ObjectURI:  "/tmp/runbook.md",
		UploadedBy: "admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.ReplaceChunks(ctx, document.ID, []DocumentChunk{{
		Text:        "check pod status",
		TokenCount:  3,
		SourceStart: 0,
		SourceEnd:   16,
	}}); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.CompleteIngestion(ctx, document.ID); err != nil {
		t.Fatal(err)
	}

	embedded, err := EmbeddingService{
		Repository: repository,
		Provider:   DeterministicEmbeddingProvider{Dimensions: 8},
	}.EmbedDocument(ctx, document.ID)
	if err != nil {
		t.Fatal(err)
	}
	if embedded.EmbeddingStatus != "ready" {
		t.Fatalf("embedding status = %q, want ready", embedded.EmbeddingStatus)
	}

	vectors, err := repository.ListEmbeddingsContext(ctx, document.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(vectors) != 1 {
		t.Fatalf("vectors len = %d, want 1", len(vectors))
	}
	if vectors[0].Model != MockEmbeddingModel {
		t.Fatalf("model = %q", vectors[0].Model)
	}
	if vectors[0].Dimensions != 8 || len(vectors[0].Embedding) != 8 {
		t.Fatalf("vector = %+v, want 8 dimensions and 8 bytes", vectors[0])
	}
}

func TestPostgresReembeddingMarksFailures(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.newID = fixedIDs(
		"00000000-0000-4000-8000-000000000031",
		"00000000-0000-4000-8000-000000000032",
	)
	document, err := repository.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   "runbook.md",
		SizeBytes:  20,
		ObjectURI:  "/tmp/runbook.md",
		UploadedBy: "admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.ReplaceChunks(ctx, document.ID, []DocumentChunk{{
		Text:       "check pod status",
		TokenCount: 3,
	}}); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.CompleteIngestion(ctx, document.ID); err != nil {
		t.Fatal(err)
	}

	result, err := EmbeddingService{
		Repository: repository,
		Provider:   staticEmbeddingProvider{err: errors.New("provider unavailable")},
	}.ReembedReadyDocuments(ctx, ReembeddingRequest{Limit: 10})
	if err == nil {
		t.Fatal("expected re-embedding error")
	}
	if result.Processed != 1 || result.Failed != 1 {
		t.Fatalf("result = %+v, want one failed document", result)
	}

	updated, err := repository.GetContext(ctx, document.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.EmbeddingStatus != "failed" {
		t.Fatalf("embedding status = %q, want failed", updated.EmbeddingStatus)
	}
	if updated.LastError != "provider unavailable" {
		t.Fatalf("last error = %q, want provider unavailable", updated.LastError)
	}
}

func TestPGVectorReadinessReportsMissingExtension(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = CheckPGVectorReady(ctx, db)
	if err != nil && !strings.Contains(err.Error(), "pgvector extension is not ready") {
		t.Fatalf("error = %v, want pgvector readiness message", err)
	}
}

func TestPostgresRetrieverRanksEmbeddedChunks(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.newID = fixedIDs(
		"00000000-0000-4000-8000-000000000021",
		"00000000-0000-4000-8000-000000000022",
		"00000000-0000-4000-8000-000000000023",
	)
	document, err := repository.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   "runbook.md",
		SizeBytes:  30,
		ObjectURI:  "/tmp/runbook.md",
		UploadedBy: "admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.ReplaceChunks(ctx, document.ID, []DocumentChunk{
		{Text: "restart deployment", TokenCount: 2},
		{Text: "check pod status", TokenCount: 3},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.CompleteIngestion(ctx, document.ID); err != nil {
		t.Fatal(err)
	}
	chunks, err := repository.ListChunksContext(ctx, document.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.ReplaceEmbeddings(ctx, []EmbeddingVector{
		{ChunkID: chunks[0].ID, Model: "test", Dimensions: 2, Embedding: []byte{0, 10}},
		{ChunkID: chunks[1].ID, Model: "test", Dimensions: 2, Embedding: []byte{10, 0}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.CompleteEmbedding(ctx, document.ID); err != nil {
		t.Fatal(err)
	}

	result, err := PostgresRetriever{
		Repository: repository,
		Embedder: staticEmbeddingProvider{
			vectors: []EmbeddingVector{{ChunkID: "query", Model: "test", Dimensions: 2, Embedding: []byte{10, 0}}},
		},
	}.Retrieve(ctx, RetrievalRequest{
		Message:     "pod status",
		DocumentIDs: []string{document.ID},
		Limit:       1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Citations) != 1 {
		t.Fatalf("citations len = %d, want 1", len(result.Citations))
	}
	if result.Citations[0].ChunkID != chunks[1].ID {
		t.Fatalf("chunk id = %q, want %q", result.Citations[0].ChunkID, chunks[1].ID)
	}
	if result.Citations[0].Rank != 1 || result.Citations[0].Score != "1.0000" {
		t.Fatalf("citation = %+v", result.Citations[0])
	}
}

func TestPGVectorLiveRAGSmoke(t *testing.T) {
	dsn := os.Getenv("CYOPS_PGVECTOR_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_PGVECTOR_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}
	if err := ApplyPGVectorEmbeddingMigration(ctx, db, 2); err != nil {
		t.Fatal(err)
	}
	if err := seedPGVectorSmokeData(ctx, db); err != nil {
		t.Fatal(err)
	}

	metrics := NewRetrievalMetrics()
	server := NewServerWithOptions(ServerOptions{
		Provider:  staticAnswerProvider{answer: "Use the cited runbook."},
		Documents: NewPostgresDocumentRepository(db, "opsmate"),
		Retriever: PostgresRetriever{
			Repository: NewPostgresDocumentRepository(db, "opsmate"),
			Embedder: staticEmbeddingProvider{
				vectors: []EmbeddingVector{{ChunkID: "query", Model: "test", Dimensions: 2, Embedding: []byte{10, 0}}},
			},
			Mode:      "pgvector",
			Observer:  metrics,
			SlowAfter: time.Nanosecond,
		},
		Metrics: metrics,
	})

	chat := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":"pod status","rag":{"enabled":true,"documentIds":["00000000-0000-4000-8000-000000000101"]}}`))
	chatRecorder := httptest.NewRecorder()
	server.ServeHTTP(chatRecorder, chat)
	if chatRecorder.Code != http.StatusOK {
		t.Fatalf("chat status = %d, want %d: %s", chatRecorder.Code, http.StatusOK, chatRecorder.Body.String())
	}
	var chatResponse ChatResponse
	if err := json.NewDecoder(chatRecorder.Body).Decode(&chatResponse); err != nil {
		t.Fatal(err)
	}
	if len(chatResponse.Citations) == 0 {
		t.Fatalf("citations len = 0, want pgvector-ranked citation")
	}
	if chatResponse.Citations[0].ChunkID != "00000000-0000-4000-8000-000000000103" {
		t.Fatalf("first citation = %+v, want pod status chunk first", chatResponse.Citations[0])
	}
	if chatResponse.Citations[0].Score != "1.0000" {
		t.Fatalf("first citation score = %q, want 1.0000", chatResponse.Citations[0].Score)
	}

	metricsRequest := httptest.NewRequest(http.MethodGet, "/api/ops/retrieval-metrics", nil)
	metricsRecorder := httptest.NewRecorder()
	server.ServeHTTP(metricsRecorder, metricsRequest)
	if metricsRecorder.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d: %s", metricsRecorder.Code, http.StatusOK, metricsRecorder.Body.String())
	}
	var snapshot RetrievalMetricsSnapshot
	if err := json.NewDecoder(metricsRecorder.Body).Decode(&snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.ByMode["pgvector"] != 1 || snapshot.Total != 1 {
		t.Fatalf("metrics snapshot = %+v, want one pgvector retrieval", snapshot)
	}
	if snapshot.Last.ResultCount != len(chatResponse.Citations) {
		t.Fatalf("last result count = %d, want %d", snapshot.Last.ResultCount, len(chatResponse.Citations))
	}
}

func seedPGVectorSmokeData(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
INSERT INTO cyops_documents (
	id, namespace, filename, size_bytes, object_uri, status, embedding_status, uploaded_by, created_at, updated_at
) VALUES (
	'00000000-0000-4000-8000-000000000101', 'opsmate', 'runbook.md', 32, '/tmp/runbook.md', 'ready', 'ready', 'admin', now(), now()
);
INSERT INTO cyops_document_chunks (
	id, document_id, chunk_index, text, token_count, source_start, source_end
) VALUES
	('00000000-0000-4000-8000-000000000102', '00000000-0000-4000-8000-000000000101', 0, 'restart deployment', 2, 0, 18),
	('00000000-0000-4000-8000-000000000103', '00000000-0000-4000-8000-000000000101', 1, 'check pod status', 3, 19, 35);
INSERT INTO cyops_document_embeddings (
	chunk_id, model, dimensions, embedding
) VALUES
	('00000000-0000-4000-8000-000000000102', 'test', 2, '[0,10]'::vector),
	('00000000-0000-4000-8000-000000000103', 'test', 2, '[10,0]'::vector);
`)
	return err
}

type staticAnswerProvider struct {
	answer string
}

func (p staticAnswerProvider) Chat(ProviderRequest) (ProviderResponse, error) {
	return ProviderResponse{Answer: p.answer, RawProvider: "test"}, nil
}

type staticEmbeddingProvider struct {
	vectors []EmbeddingVector
	err     error
}

func (p staticEmbeddingProvider) Embed(context.Context, []DocumentChunk) ([]EmbeddingVector, error) {
	if p.err != nil {
		return nil, p.err
	}
	if len(p.vectors) == 0 {
		return nil, errors.New("no vectors")
	}
	return p.vectors, nil
}

func fixedIDs(ids ...string) func() (string, error) {
	index := 0
	return func() (string, error) {
		if index >= len(ids) {
			return ids[len(ids)-1], nil
		}
		id := ids[index]
		index++
		return id, nil
	}
}
