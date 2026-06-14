package appserver

import (
	"context"
	"time"
)

type ChatRequest struct {
	SessionID      string         `json:"sessionId,omitempty"`
	Message        string         `json:"message"`
	Provider       string         `json:"provider,omitempty"`
	RAG            RAGRequest     `json:"rag,omitempty"`
	ClusterContext map[string]any `json:"clusterContext,omitempty"`
}

type RAGRequest struct {
	Enabled     bool     `json:"enabled"`
	DocumentIDs []string `json:"documentIds,omitempty"`
	Scope       string   `json:"scope,omitempty"`
}

type ChatResponse struct {
	SessionID string     `json:"sessionId"`
	MessageID string     `json:"messageId"`
	Answer    string     `json:"answer"`
	Provider  string     `json:"provider"`
	Citations []Citation `json:"citations,omitempty"`
	Warnings  []string   `json:"warnings,omitempty"`
}

type Citation struct {
	DocumentID string `json:"documentId"`
	Title      string `json:"title"`
	ChunkID    string `json:"chunkId"`
	Rank       int    `json:"rank,omitempty"`
	Score      string `json:"score,omitempty"`
}

type Document struct {
	ID              string    `json:"id"`
	Filename        string    `json:"filename"`
	Status          string    `json:"status"`
	SizeBytes       int64     `json:"sizeBytes"`
	ObjectURI       string    `json:"objectUri,omitempty"`
	ChunkCount      int       `json:"chunkCount"`
	EmbeddingStatus string    `json:"embeddingStatus"`
	UploadedBy      string    `json:"uploadedBy"`
	CreatedAt       time.Time `json:"createdAt"`
	LastError       string    `json:"lastError,omitempty"`
}

type DocumentListResponse struct {
	Items []Document `json:"items"`
}

type DocumentUploadResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Status   string `json:"status"`
}

type DocumentDeleteResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ProviderRequest struct {
	Message        string
	Context        []ProviderContext
	ClusterContext map[string]any
}

type ProviderContext struct {
	Type   string
	Text   string
	Source string
}

type ProviderResponse struct {
	Answer      string
	RawProvider string
}

type ChatProvider interface {
	Chat(request ProviderRequest) (ProviderResponse, error)
}

type DocumentRepository interface {
	List() []Document
	Create(filename string, sizeBytes int64, uploadedBy string) Document
	Get(id string) (Document, bool)
	MarkDeleting(id string) (Document, bool)
}

type CreateStoredDocumentInput struct {
	ID         string
	Filename   string
	SizeBytes  int64
	ObjectURI  string
	UploadedBy string
}

type StoredDocumentRepository interface {
	CreateStored(ctx context.Context, input CreateStoredDocumentInput) (Document, error)
}
