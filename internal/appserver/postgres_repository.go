package appserver

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type PostgresDocumentRepository struct {
	db        *sql.DB
	namespace string
	now       func() time.Time
	newID     func() (string, error)
}

type RankedChunk struct {
	Document Document
	Chunk    DocumentChunk
	Score    float64
}

func NewPostgresDocumentRepository(db *sql.DB, namespace string) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{
		db:        db,
		namespace: namespace,
		now:       time.Now,
		newID:     newUUID,
	}
}

func (r *PostgresDocumentRepository) List() []Document {
	documents, err := r.ListContext(context.Background())
	if err != nil {
		return nil
	}
	return documents
}

func (r *PostgresDocumentRepository) Create(filename string, sizeBytes int64, uploadedBy string) Document {
	document, err := r.CreateStored(context.Background(), CreateStoredDocumentInput{
		Filename:   filename,
		SizeBytes:  sizeBytes,
		UploadedBy: uploadedBy,
	})
	if err != nil {
		return Document{}
	}
	return document
}

func (r *PostgresDocumentRepository) Get(id string) (Document, bool) {
	document, err := r.GetContext(context.Background(), id)
	return document, err == nil
}

func (r *PostgresDocumentRepository) MarkDeleting(id string) (Document, bool) {
	document, err := r.MarkDeletingContext(context.Background(), id)
	return document, err == nil
}

func (r *PostgresDocumentRepository) ListContext(ctx context.Context) ([]Document, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT d.id, d.filename, d.status, d.size_bytes, d.object_uri, d.embedding_status, d.uploaded_by, d.created_at, d.last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = d.id) AS chunk_count
FROM cyops_documents d
WHERE d.namespace = $1 AND d.deleted_at IS NULL
ORDER BY d.created_at, d.id`, r.namespace)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		document, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, rows.Err()
}

func (r *PostgresDocumentRepository) ListReadyDocumentsForReembedding(ctx context.Context, limit int) ([]Document, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT d.id, d.filename, d.status, d.size_bytes, d.object_uri, d.embedding_status, d.uploaded_by, d.created_at, d.last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = d.id) AS chunk_count
FROM cyops_documents d
WHERE d.namespace = $1
  AND d.deleted_at IS NULL
  AND d.status = 'ready'
  AND d.embedding_status IN ('pending', 'ready', 'failed')
ORDER BY d.updated_at, d.id
LIMIT $2`, r.namespace, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		document, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, rows.Err()
}

func (r *PostgresDocumentRepository) CreateContext(ctx context.Context, filename string, sizeBytes int64, uploadedBy string) (Document, error) {
	return r.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   filename,
		SizeBytes:  sizeBytes,
		UploadedBy: uploadedBy,
	})
}

func (r *PostgresDocumentRepository) CreateStored(ctx context.Context, input CreateStoredDocumentInput) (Document, error) {
	id := input.ID
	if id == "" {
		generatedID, err := r.newID()
		if err != nil {
			return Document{}, err
		}
		id = generatedID
	}
	createdAt := r.now().UTC()
	document := Document{
		ID:              id,
		Filename:        input.Filename,
		Status:          "uploaded",
		SizeBytes:       input.SizeBytes,
		ObjectURI:       input.ObjectURI,
		EmbeddingStatus: "pending",
		UploadedBy:      input.UploadedBy,
		CreatedAt:       createdAt,
	}

	_, err := r.db.ExecContext(ctx, `
INSERT INTO cyops_documents (
	id, namespace, filename, size_bytes, object_uri, status, embedding_status, uploaded_by, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)`,
		document.ID,
		r.namespace,
		document.Filename,
		document.SizeBytes,
		document.ObjectURI,
		document.Status,
		document.EmbeddingStatus,
		document.UploadedBy,
		document.CreatedAt,
	)
	if err != nil {
		return Document{}, err
	}
	return document, nil
}

func (r *PostgresDocumentRepository) GetContext(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT d.id, d.filename, d.status, d.size_bytes, d.object_uri, d.embedding_status, d.uploaded_by, d.created_at, d.last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = d.id) AS chunk_count
FROM cyops_documents d
WHERE d.id = $1 AND d.namespace = $2 AND d.deleted_at IS NULL`, id, r.namespace)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) MarkDeletingContext(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET status = 'deleting', updated_at = $3
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error, 0`,
		id,
		r.namespace,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) BeginIngestion(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET status = 'processing', embedding_status = 'pending', last_error = '', updated_at = $3
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = cyops_documents.id)`,
		id,
		r.namespace,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) CompleteIngestion(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET status = 'ready', embedding_status = 'pending', last_error = '', updated_at = $3
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = cyops_documents.id)`,
		id,
		r.namespace,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) FailIngestion(ctx context.Context, id string, message string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET status = 'failed', embedding_status = 'failed', last_error = $3, updated_at = $4
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = cyops_documents.id)`,
		id,
		r.namespace,
		message,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) BeginEmbedding(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET embedding_status = 'processing', last_error = '', updated_at = $3
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = cyops_documents.id)`,
		id,
		r.namespace,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) CompleteEmbedding(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET embedding_status = 'ready', last_error = '', updated_at = $3
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = cyops_documents.id)`,
		id,
		r.namespace,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) FailEmbedding(ctx context.Context, id string, message string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET embedding_status = 'failed', last_error = $3, updated_at = $4
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error,
	(SELECT COUNT(*) FROM cyops_document_chunks c WHERE c.document_id = cyops_documents.id)`,
		id,
		r.namespace,
		message,
		r.now().UTC(),
	)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) ReplaceChunks(ctx context.Context, documentID string, chunks []DocumentChunk) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM cyops_document_chunks WHERE document_id = $1", documentID); err != nil {
		return err
	}
	for index, chunk := range chunks {
		chunkID := chunk.ID
		if chunkID == "" {
			generatedID, err := r.newID()
			if err != nil {
				return err
			}
			chunkID = generatedID
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO cyops_document_chunks (
	id, document_id, chunk_index, text, token_count, source_start, source_end
) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			chunkID,
			documentID,
			index,
			chunk.Text,
			chunk.TokenCount,
			chunk.SourceStart,
			chunk.SourceEnd,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostgresDocumentRepository) ListChunksContext(ctx context.Context, documentID string) ([]DocumentChunk, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, document_id, chunk_index, text, token_count, source_start, source_end
FROM cyops_document_chunks
WHERE document_id = $1
ORDER BY chunk_index`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []DocumentChunk
	for rows.Next() {
		var chunk DocumentChunk
		if err := rows.Scan(
			&chunk.ID,
			&chunk.DocumentID,
			&chunk.ChunkIndex,
			&chunk.Text,
			&chunk.TokenCount,
			&chunk.SourceStart,
			&chunk.SourceEnd,
		); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, rows.Err()
}

func (r *PostgresDocumentRepository) ReplaceEmbeddings(ctx context.Context, vectors []EmbeddingVector) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, vector := range vectors {
		if _, err := tx.ExecContext(ctx, "DELETE FROM cyops_document_embeddings WHERE chunk_id = $1", vector.ChunkID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO cyops_document_embeddings (
	chunk_id, model, dimensions, embedding
) VALUES ($1, $2, $3, $4)`,
			vector.ChunkID,
			vector.Model,
			vector.Dimensions,
			vector.Embedding,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostgresDocumentRepository) ListEmbeddingsContext(ctx context.Context, documentID string) ([]EmbeddingVector, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT e.chunk_id, e.model, e.dimensions, e.embedding
FROM cyops_document_embeddings e
JOIN cyops_document_chunks c ON c.id = e.chunk_id
WHERE c.document_id = $1
ORDER BY c.chunk_index`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vectors []EmbeddingVector
	for rows.Next() {
		var vector EmbeddingVector
		if err := rows.Scan(
			&vector.ChunkID,
			&vector.Model,
			&vector.Dimensions,
			&vector.Embedding,
		); err != nil {
			return nil, err
		}
		vectors = append(vectors, vector)
	}
	return vectors, rows.Err()
}

func (r *PostgresDocumentRepository) ListRankedChunksPGVector(ctx context.Context, documentIDs []string, queryEmbedding []byte, limit int) ([]RankedChunk, error) {
	query, args, err := rankedChunksPGVectorSQL(r.namespace, documentIDs, bytesToVectorLiteral(queryEmbedding), limit)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranked []RankedChunk
	for rows.Next() {
		var item RankedChunk
		if err := rows.Scan(
			&item.Document.ID,
			&item.Document.Filename,
			&item.Chunk.ID,
			&item.Chunk.Text,
			&item.Score,
		); err != nil {
			return nil, err
		}
		ranked = append(ranked, item)
	}
	return ranked, rows.Err()
}

func rankedChunksPGVectorSQL(namespace string, documentIDs []string, vectorLiteral string, limit int) (string, []any, error) {
	if namespace == "" {
		return "", nil, fmt.Errorf("namespace is required")
	}
	if len(documentIDs) == 0 {
		return "", nil, fmt.Errorf("document ids are required")
	}
	if vectorLiteral == "" {
		return "", nil, fmt.Errorf("query vector is required")
	}
	if limit <= 0 {
		limit = 3
	}

	args := []any{namespace, vectorLiteral}
	placeholders := make([]string, 0, len(documentIDs))
	for _, documentID := range documentIDs {
		args = append(args, documentID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}
	args = append(args, limit)

	query := fmt.Sprintf(`
SELECT d.id, d.filename, c.id, c.text, 1 - (e.embedding <=> $2::vector) AS score
FROM cyops_document_chunks c
JOIN cyops_documents d ON d.id = c.document_id
JOIN cyops_document_embeddings e ON e.chunk_id = c.id
WHERE d.namespace = $1
  AND d.deleted_at IS NULL
  AND d.status = 'ready'
  AND d.embedding_status = 'ready'
  AND d.id IN (%s)
ORDER BY e.embedding <=> $2::vector, c.chunk_index
LIMIT $%d`, strings.Join(placeholders, ", "), len(args))
	return query, args, nil
}

func bytesToVectorLiteral(value []byte) string {
	if len(value) == 0 {
		return ""
	}
	parts := make([]string, len(value))
	for i, item := range value {
		parts[i] = fmt.Sprintf("%d", item)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

type documentScanner interface {
	Scan(dest ...any) error
}

func scanDocument(scanner documentScanner) (Document, error) {
	var document Document
	err := scanner.Scan(
		&document.ID,
		&document.Filename,
		&document.Status,
		&document.SizeBytes,
		&document.ObjectURI,
		&document.EmbeddingStatus,
		&document.UploadedBy,
		&document.CreatedAt,
		&document.LastError,
		&document.ChunkCount,
	)
	return document, err
}

func newUUID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	encoded := hex.EncodeToString(bytes[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s", encoded[0:8], encoded[8:12], encoded[12:16], encoded[16:20], encoded[20:32]), nil
}
