package appserver

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
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

type RetrievalObservation struct {
	Mode          string
	Duration      time.Duration
	Slow          bool
	FailureReason string
	ResultCount   int
}

type RetrievalObserver interface {
	ObserveRetrieval(observation RetrievalObservation)
}

type retrievalCandidate struct {
	document Document
	chunk    DocumentChunk
	score    float64
}

type Retriever interface {
	Retrieve(ctx context.Context, request RetrievalRequest) (RetrievalResult, error)
}

type PostgresRetriever struct {
	Repository *PostgresDocumentRepository
	Embedder   EmbeddingProvider
	Mode       string
	Observer   RetrievalObserver
	SlowAfter  time.Duration
}

func (r PostgresRetriever) Retrieve(ctx context.Context, request RetrievalRequest) (RetrievalResult, error) {
	start := time.Now()
	mode := retrievalModeOrDefault(r.Mode)
	var result RetrievalResult
	var failureReason string
	defer func() {
		if r.Observer != nil {
			duration := time.Since(start)
			r.Observer.ObserveRetrieval(RetrievalObservation{
				Mode:          mode,
				Duration:      duration,
				Slow:          r.SlowAfter > 0 && duration > r.SlowAfter,
				FailureReason: failureReason,
				ResultCount:   len(result.Citations),
			})
		}
	}()

	if r.Repository == nil {
		return result, nil
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 3
	}

	embedder := r.Embedder
	if embedder == nil {
		embedder = DeterministicEmbeddingProvider{}
	}
	queryVectors, err := embedder.Embed(ctx, []DocumentChunk{{ID: "query", Text: request.Message}})
	if err != nil {
		failureReason = "query_embedding_failed"
		return RetrievalResult{}, err
	}
	if len(queryVectors) != 1 {
		failureReason = "query_embedding_count_mismatch"
		return RetrievalResult{}, fmt.Errorf("query embedding returned %d vectors, want 1", len(queryVectors))
	}

	if mode == "pgvector" {
		result, err = r.retrievePGVector(ctx, request, queryVectors[0], limit)
		if err != nil {
			failureReason = "pgvector_query_failed"
			return RetrievalResult{}, err
		}
		return result, nil
	}

	result, err = r.retrieveBYTEA(ctx, request, queryVectors[0], limit)
	if err != nil {
		failureReason = "bytea_query_failed"
		return RetrievalResult{}, err
	}
	return result, nil
}

func (r PostgresRetriever) retrieveBYTEA(ctx context.Context, request RetrievalRequest, queryVector EmbeddingVector, limit int) (RetrievalResult, error) {
	var candidates []retrievalCandidate
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
			candidates = append(candidates, retrievalCandidate{
				document: document,
				chunk:    chunk,
				score:    cosineBytes(queryVector.Embedding, embedding.Embedding),
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

func (r PostgresRetriever) retrievePGVector(ctx context.Context, request RetrievalRequest, queryVector EmbeddingVector, limit int) (RetrievalResult, error) {
	rows, err := r.Repository.ListRankedChunksPGVector(ctx, request.DocumentIDs, queryVector.Embedding, limit)
	if err != nil {
		return RetrievalResult{}, err
	}

	var result RetrievalResult
	for index, row := range rows {
		result.Context = append(result.Context, ProviderContext{
			Type:   "rag_chunk",
			Text:   trimContextText(row.Chunk.Text),
			Source: row.Document.ID + "/" + row.Chunk.ID,
		})
		result.Citations = append(result.Citations, Citation{
			DocumentID: row.Document.ID,
			Title:      row.Document.Filename,
			ChunkID:    row.Chunk.ID,
			Rank:       index + 1,
			Score:      fmt.Sprintf("%.4f", row.Score),
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
