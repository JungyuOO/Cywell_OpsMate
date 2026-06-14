package appserver

import (
	"context"
	"crypto/sha256"
	"fmt"
)

const MockEmbeddingModel = "cyops-mock-embedding-v1"

type EmbeddingVector struct {
	ChunkID    string
	Model      string
	Dimensions int
	Embedding  []byte
}

type EmbeddingProvider interface {
	Embed(ctx context.Context, chunks []DocumentChunk) ([]EmbeddingVector, error)
}

type DeterministicEmbeddingProvider struct {
	Model      string
	Dimensions int
}

func (p DeterministicEmbeddingProvider) Embed(ctx context.Context, chunks []DocumentChunk) ([]EmbeddingVector, error) {
	model := p.Model
	if model == "" {
		model = MockEmbeddingModel
	}
	dimensions := p.Dimensions
	if dimensions <= 0 {
		dimensions = 32
	}

	vectors := make([]EmbeddingVector, 0, len(chunks))
	for _, chunk := range chunks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if chunk.ID == "" {
			return nil, fmt.Errorf("chunk id is required")
		}
		sum := sha256.Sum256([]byte(chunk.Text))
		embedding := make([]byte, dimensions)
		for i := range embedding {
			embedding[i] = sum[i%len(sum)]
		}
		vectors = append(vectors, EmbeddingVector{
			ChunkID:    chunk.ID,
			Model:      model,
			Dimensions: dimensions,
			Embedding:  embedding,
		})
	}
	return vectors, nil
}

type EmbeddingService struct {
	Repository *PostgresDocumentRepository
	Provider   EmbeddingProvider
}

func (s EmbeddingService) EmbedDocument(ctx context.Context, documentID string) (Document, error) {
	if s.Repository == nil {
		return Document{}, fmt.Errorf("repository is required")
	}
	provider := s.Provider
	if provider == nil {
		provider = DeterministicEmbeddingProvider{}
	}

	if _, err := s.Repository.BeginEmbedding(ctx, documentID); err != nil {
		return Document{}, err
	}
	chunks, err := s.Repository.ListChunksContext(ctx, documentID)
	if err != nil {
		return s.fail(ctx, documentID, err.Error())
	}
	if len(chunks) == 0 {
		return s.fail(ctx, documentID, "no chunks available")
	}
	vectors, err := provider.Embed(ctx, chunks)
	if err != nil {
		return s.fail(ctx, documentID, err.Error())
	}
	if err := s.Repository.ReplaceEmbeddings(ctx, vectors); err != nil {
		return s.fail(ctx, documentID, err.Error())
	}
	return s.Repository.CompleteEmbedding(ctx, documentID)
}

func (s EmbeddingService) fail(ctx context.Context, documentID string, message string) (Document, error) {
	document, err := s.Repository.FailEmbedding(ctx, documentID, message)
	if err != nil {
		return Document{}, err
	}
	return document, fmt.Errorf("%s", message)
}
