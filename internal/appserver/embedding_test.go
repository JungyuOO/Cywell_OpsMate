package appserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestHTTPEmbeddingProviderPostsChunks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret-token" {
			t.Fatalf("authorization = %q", got)
		}
		var request struct {
			Model      string `json:"model"`
			Dimensions int    `json:"dimensions"`
			Chunks     []struct {
				ChunkID string `json:"chunkId"`
				Text    string `json:"text"`
			} `json:"chunks"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Model != "embedding-model" {
			t.Fatalf("model = %q", request.Model)
		}
		if request.Dimensions != 4 {
			t.Fatalf("dimensions = %d, want 4", request.Dimensions)
		}
		if len(request.Chunks) != 1 || request.Chunks[0].Text != "check pod status" {
			t.Fatalf("chunks = %+v", request.Chunks)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"vectors": []map[string]any{{
				"chunkId":    request.Chunks[0].ChunkID,
				"model":      request.Model,
				"dimensions": request.Dimensions,
				"embedding":  []byte{1, 2, 3, 4},
			}},
		})
	}))
	defer server.Close()

	vectors, err := HTTPEmbeddingProvider{
		EndpointURL: server.URL,
		Model:       "embedding-model",
		Dimensions:  4,
		Token:       "secret-token",
	}.Embed(context.Background(), []DocumentChunk{{
		ID:   "chunk-001",
		Text: "check pod status",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(vectors) != 1 {
		t.Fatalf("vectors len = %d, want 1", len(vectors))
	}
	if !bytes.Equal(vectors[0].Embedding, []byte{1, 2, 3, 4}) {
		t.Fatalf("embedding = %v", vectors[0].Embedding)
	}
}
