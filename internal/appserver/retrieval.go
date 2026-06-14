package appserver

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
)

type RetrievalRequest struct {
	Message     string
	DocumentIDs []string
	Limit       int
}

type RetrievalResult struct {
	Context   []ProviderContext
	Citations []Citation
}

type Retriever interface {
	Retrieve(ctx context.Context, request RetrievalRequest) (RetrievalResult, error)
}

type PostgresRetriever struct {
	Repository *PostgresDocumentRepository
	Embedder   EmbeddingProvider
}

func (r PostgresRetriever) Retrieve(ctx context.Context, request RetrievalRequest) (RetrievalResult, error) {
	if r.Repository == nil {
		return RetrievalResult{}, nil
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 3
	}

	type candidate struct {
		document Document
		chunk    DocumentChunk
		score    float64
	}
	var candidates []candidate
	embedder := r.Embedder
	if embedder == nil {
		embedder = DeterministicEmbeddingProvider{}
	}
	queryVectors, err := embedder.Embed(ctx, []DocumentChunk{{ID: "query", Text: request.Message}})
	if err != nil {
		return RetrievalResult{}, err
	}
	if len(queryVectors) != 1 {
		return RetrievalResult{}, fmt.Errorf("query embedding returned %d vectors, want 1", len(queryVectors))
	}

	for _, documentID := range request.DocumentIDs {
		document, err := r.Repository.GetContext(ctx, documentID)
		if err != nil || document.Status != "ready" || document.EmbeddingStatus != "ready" {
			continue
		}
		chunks, err := r.Repository.ListChunksContext(ctx, documentID)
		if err != nil {
			return RetrievalResult{}, err
		}
		embeddings, err := r.Repository.ListEmbeddingsContext(ctx, documentID)
		if err != nil {
			return RetrievalResult{}, err
		}
		byChunkID := map[string]EmbeddingVector{}
		for _, embedding := range embeddings {
			byChunkID[embedding.ChunkID] = embedding
		}
		for _, chunk := range chunks {
			embedding, ok := byChunkID[chunk.ID]
			if !ok {
				continue
			}
			candidates = append(candidates, candidate{
				document: document,
				chunk:    chunk,
				score:    cosineBytes(queryVectors[0].Embedding, embedding.Embedding),
			})
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].chunk.ChunkIndex < candidates[j].chunk.ChunkIndex
		}
		return candidates[i].score > candidates[j].score
	})

	if len(candidates) > limit {
		candidates = candidates[:limit]
	}

	var result RetrievalResult
	for index, candidate := range candidates {
		result.Context = append(result.Context, ProviderContext{
			Type:   "rag_chunk",
			Text:   trimContextText(candidate.chunk.Text),
			Source: candidate.document.ID + "/" + candidate.chunk.ID,
		})
		result.Citations = append(result.Citations, Citation{
			DocumentID: candidate.document.ID,
			Title:      candidate.document.Filename,
			ChunkID:    candidate.chunk.ID,
			Rank:       index + 1,
			Score:      fmt.Sprintf("%.4f", candidate.score),
		})
	}
	return result, nil
}

func cosineBytes(left []byte, right []byte) float64 {
	limit := len(left)
	if len(right) < limit {
		limit = len(right)
	}
	if limit == 0 {
		return 0
	}

	var dot float64
	var leftNorm float64
	var rightNorm float64
	for i := 0; i < limit; i++ {
		l := float64(left[i])
		r := float64(right[i])
		dot += l * r
		leftNorm += l * l
		rightNorm += r * r
	}
	if leftNorm == 0 || rightNorm == 0 {
		return 0
	}
	return dot / (math.Sqrt(leftNorm) * math.Sqrt(rightNorm))
}

func trimContextText(text string) string {
	text = strings.TrimSpace(text)
	const maxRunes = 1200
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return string(runes[:maxRunes])
}
