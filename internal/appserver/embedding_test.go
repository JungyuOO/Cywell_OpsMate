package appserver

import (
	"bytes"
	"context"
	"testing"
)

func TestDeterministicEmbeddingProviderIsStable(t *testing.T) {
	chunks := []DocumentChunk{{
		ID:   "chunk-001",
		Text: "check pod status",
	}}

	provider := DeterministicEmbeddingProvider{Dimensions: 8}
	first, err := provider.Embed(context.Background(), chunks)
	if err != nil {
		t.Fatal(err)
	}
	second, err := provider.Embed(context.Background(), chunks)
	if err != nil {
		t.Fatal(err)
	}

	if len(first) != 1 {
		t.Fatalf("vectors len = %d, want 1", len(first))
	}
	if first[0].Model != MockEmbeddingModel {
		t.Fatalf("model = %q", first[0].Model)
	}
	if first[0].Dimensions != 8 {
		t.Fatalf("dimensions = %d, want 8", first[0].Dimensions)
	}
	if !bytes.Equal(first[0].Embedding, second[0].Embedding) {
		t.Fatalf("embedding = %v, want stable %v", first[0].Embedding, second[0].Embedding)
	}
}
