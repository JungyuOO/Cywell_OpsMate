package appserver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestTextFileParserParsesPlainText(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runbook.md")
	if err := os.WriteFile(path, []byte("# Runbook\n\nCheck pod status."), 0o600); err != nil {
		t.Fatal(err)
	}

	sections, err := TextFileParser{}.Parse(context.Background(), filepath.ToSlash(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 {
		t.Fatalf("sections len = %d, want 1", len(sections))
	}
	if sections[0].Text != "# Runbook\n\nCheck pod status." {
		t.Fatalf("text = %q", sections[0].Text)
	}
}

func TestFixedRuneChunkerChunksDeterministically(t *testing.T) {
	chunks := FixedRuneChunker{MaxRunes: 5, Overlap: 1}.Chunk([]ParsedSection{{
		Text:        "abcdefghijkl",
		SourceStart: 10,
		SourceEnd:   22,
	}})

	if len(chunks) != 3 {
		t.Fatalf("chunks len = %d, want 3", len(chunks))
	}
	if chunks[0].Text != "abcde" || chunks[0].SourceStart != 10 || chunks[0].SourceEnd != 15 {
		t.Fatalf("chunk 0 = %+v", chunks[0])
	}
	if chunks[1].Text != "efghi" || chunks[1].SourceStart != 14 || chunks[1].SourceEnd != 19 {
		t.Fatalf("chunk 1 = %+v", chunks[1])
	}
	if chunks[2].Text != "ijkl" || chunks[2].SourceStart != 18 || chunks[2].SourceEnd != 22 {
		t.Fatalf("chunk 2 = %+v", chunks[2])
	}
}
