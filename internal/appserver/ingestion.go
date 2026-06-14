package appserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type ParsedSection struct {
	Text        string
	SourceStart int
	SourceEnd   int
}

type DocumentChunk struct {
	ID          string
	DocumentID  string
	ChunkIndex  int
	Text        string
	TokenCount  int
	SourceStart int
	SourceEnd   int
}

type DocumentParser interface {
	Parse(ctx context.Context, objectURI string) ([]ParsedSection, error)
}

type TextFileParser struct{}

func (TextFileParser) Parse(ctx context.Context, objectURI string) ([]ParsedSection, error) {
	if strings.TrimSpace(objectURI) == "" {
		return nil, fmt.Errorf("object uri is required")
	}
	path := filepath.FromSlash(objectURI)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	text := strings.TrimSpace(string(content))
	if text == "" {
		return nil, fmt.Errorf("document is empty")
	}
	return []ParsedSection{{
		Text:        text,
		SourceStart: 0,
		SourceEnd:   utf8.RuneCountInString(text),
	}}, nil
}

type FixedRuneChunker struct {
	MaxRunes int
	Overlap  int
}

func (c FixedRuneChunker) Chunk(sections []ParsedSection) []DocumentChunk {
	maxRunes := c.MaxRunes
	if maxRunes <= 0 {
		maxRunes = 1200
	}
	overlap := c.Overlap
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= maxRunes {
		overlap = 0
	}

	var chunks []DocumentChunk
	for _, section := range sections {
		runes := []rune(section.Text)
		for start := 0; start < len(runes); {
			end := start + maxRunes
			if end > len(runes) {
				end = len(runes)
			}
			text := strings.TrimSpace(string(runes[start:end]))
			if text != "" {
				chunks = append(chunks, DocumentChunk{
					ChunkIndex:  len(chunks),
					Text:        text,
					TokenCount:  estimateTokenCount(text),
					SourceStart: section.SourceStart + start,
					SourceEnd:   section.SourceStart + end,
				})
			}
			if end == len(runes) {
				break
			}
			start = end - overlap
		}
	}
	return chunks
}

type IngestionService struct {
	Repository *PostgresDocumentRepository
	Parser     DocumentParser
	Chunker    FixedRuneChunker
}

func (s IngestionService) IngestDocument(ctx context.Context, documentID string) (Document, error) {
	if s.Repository == nil {
		return Document{}, fmt.Errorf("repository is required")
	}
	parser := s.Parser
	if parser == nil {
		parser = TextFileParser{}
	}

	document, err := s.Repository.BeginIngestion(ctx, documentID)
	if err != nil {
		return Document{}, err
	}
	if document.ObjectURI == "" {
		return s.fail(ctx, documentID, "object uri is required")
	}

	sections, err := parser.Parse(ctx, document.ObjectURI)
	if err != nil {
		return s.fail(ctx, documentID, err.Error())
	}
	chunks := s.Chunker.Chunk(sections)
	if len(chunks) == 0 {
		return s.fail(ctx, documentID, "no chunks produced")
	}
	if err := s.Repository.ReplaceChunks(ctx, documentID, chunks); err != nil {
		return s.fail(ctx, documentID, err.Error())
	}
	return s.Repository.CompleteIngestion(ctx, documentID)
}

func (s IngestionService) fail(ctx context.Context, documentID string, message string) (Document, error) {
	document, err := s.Repository.FailIngestion(ctx, documentID, message)
	if err != nil {
		return Document{}, err
	}
	return document, fmt.Errorf("%s", message)
}

func estimateTokenCount(text string) int {
	fields := strings.Fields(text)
	if len(fields) > 0 {
		return len(fields)
	}
	return utf8.RuneCountInString(text)
}
