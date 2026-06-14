package appserver

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

type PostgresDocumentRepository struct {
	db        *sql.DB
	namespace string
	now       func() time.Time
	newID     func() (string, error)
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
SELECT id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error
FROM cyops_documents
WHERE namespace = $1 AND deleted_at IS NULL
ORDER BY created_at, id`, r.namespace)
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
SELECT id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error
FROM cyops_documents
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL`, id, r.namespace)
	return scanDocument(row)
}

func (r *PostgresDocumentRepository) MarkDeletingContext(ctx context.Context, id string) (Document, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE cyops_documents
SET status = 'deleting', updated_at = $3
WHERE id = $1 AND namespace = $2 AND deleted_at IS NULL
RETURNING id, filename, status, size_bytes, object_uri, embedding_status, uploaded_by, created_at, last_error`,
		id,
		r.namespace,
		r.now().UTC(),
	)
	return scanDocument(row)
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
