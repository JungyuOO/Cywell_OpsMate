package appserver

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MockProvider struct{}

func (MockProvider) Chat(request ProviderRequest) (ProviderResponse, error) {
	return ProviderResponse{
		Answer:      fmt.Sprintf("Mocked CYOps response for: %s", request.Message),
		RawProvider: "mock",
	}, nil
}

type MemoryDocumentRepository struct {
	mu        sync.Mutex
	nextID    int
	documents map[string]Document
	now       func() time.Time
}

func NewMemoryDocumentRepository() *MemoryDocumentRepository {
	return &MemoryDocumentRepository{
		nextID:    1,
		documents: map[string]Document{},
		now:       time.Now,
	}
}

func (r *MemoryDocumentRepository) List() []Document {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]Document, 0, len(r.documents))
	for _, document := range r.documents {
		items = append(items, document)
	}
	return items
}

func (r *MemoryDocumentRepository) Create(filename string, sizeBytes int64, uploadedBy string) Document {
	document, _ := r.CreateStored(context.Background(), CreateStoredDocumentInput{
		Filename:   filename,
		SizeBytes:  sizeBytes,
		UploadedBy: uploadedBy,
	})
	return document
}

func (r *MemoryDocumentRepository) CreateStored(_ context.Context, input CreateStoredDocumentInput) (Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := input.ID
	if id == "" {
		id = fmt.Sprintf("doc-%03d", r.nextID)
	}
	r.nextID++
	document := Document{
		ID:              id,
		Filename:        input.Filename,
		Status:          "uploaded",
		SizeBytes:       input.SizeBytes,
		ObjectURI:       input.ObjectURI,
		EmbeddingStatus: "pending",
		UploadedBy:      input.UploadedBy,
		CreatedAt:       r.now().UTC(),
	}
	r.documents[id] = document
	return document, nil
}

func (r *MemoryDocumentRepository) Get(id string) (Document, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	document, ok := r.documents[id]
	return document, ok
}

func (r *MemoryDocumentRepository) MarkDeleting(id string) (Document, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	document, ok := r.documents[id]
	if !ok {
		return Document{}, false
	}
	document.Status = "deleting"
	r.documents[id] = document
	return document, true
}
