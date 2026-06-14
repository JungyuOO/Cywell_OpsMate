package appserver

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
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

type HTTPEmbeddingProvider struct {
	EndpointURL string
	Model       string
	Dimensions  int
	Client      HTTPDoer
}

func (p HTTPEmbeddingProvider) Embed(ctx context.Context, chunks []DocumentChunk) ([]EmbeddingVector, error) {
	if p.EndpointURL == "" {
		return nil, fmt.Errorf("embedding endpoint is required")
	}
	body, err := json.Marshal(map[string]any{
		"model":      p.Model,
		"dimensions": p.Dimensions,
		"chunks":     embeddingProviderChunks(chunks),
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.EndpointURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("embedding provider returned status %d", response.StatusCode)
	}

	var decoded struct {
		Vectors []struct {
			ChunkID    string `json:"chunkId"`
			Model      string `json:"model"`
			Dimensions int    `json:"dimensions"`
			Embedding  []byte `json:"embedding"`
		} `json:"vectors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	vectors := make([]EmbeddingVector, 0, len(decoded.Vectors))
	for _, vector := range decoded.Vectors {
		model := vector.Model
		if model == "" {
			model = p.Model
		}
		dimensions := vector.Dimensions
		if dimensions == 0 {
			dimensions = p.Dimensions
		}
		vectors = append(vectors, EmbeddingVector{
			ChunkID:    vector.ChunkID,
			Model:      model,
			Dimensions: dimensions,
			Embedding:  vector.Embedding,
		})
	}
	return vectors, nil
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

func embeddingProviderChunks(chunks []DocumentChunk) []map[string]string {
	payload := make([]map[string]string, 0, len(chunks))
	for _, chunk := range chunks {
		payload = append(payload, map[string]string{
			"chunkId": chunk.ID,
			"text":    chunk.Text,
		})
	}
	return payload
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
