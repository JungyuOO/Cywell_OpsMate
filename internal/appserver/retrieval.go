package appserver

import (
	"context"
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
}

func (r PostgresRetriever) Retrieve(ctx context.Context, request RetrievalRequest) (RetrievalResult, error) {
	if r.Repository == nil {
		return RetrievalResult{}, nil
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 3
	}

	var result RetrievalResult
	for _, documentID := range request.DocumentIDs {
		if len(result.Citations) >= limit {
			break
		}
		document, err := r.Repository.GetContext(ctx, documentID)
		if err != nil || document.Status != "ready" {
			continue
		}
		chunks, err := r.Repository.ListChunksContext(ctx, documentID)
		if err != nil {
			return RetrievalResult{}, err
		}
		for _, chunk := range chunks {
			if len(result.Citations) >= limit {
				break
			}
			result.Context = append(result.Context, ProviderContext{
				Type:   "rag_chunk",
				Text:   trimContextText(chunk.Text),
				Source: document.ID + "/" + chunk.ID,
			})
			result.Citations = append(result.Citations, Citation{
				DocumentID: document.ID,
				Title:      document.Filename,
				ChunkID:    chunk.ID,
			})
		}
	}
	return result, nil
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
