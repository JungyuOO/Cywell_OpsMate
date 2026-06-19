package appserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMinimalProviderContextOmitsDocumentText(t *testing.T) {
	context := MinimalProviderContext([]Citation{
		{DocumentID: "doc-001", ChunkID: "chunk-001", Title: "runbook.pdf"},
	})

	if len(context) != 1 {
		t.Fatalf("context len = %d, want 1", len(context))
	}
	if context[0].Text != "" {
		t.Fatalf("context text = %q, want empty until retrieval policy is implemented", context[0].Text)
	}
	if context[0].Source != "doc-001/chunk-001" {
		t.Fatalf("context source = %q, want doc-001/chunk-001", context[0].Source)
	}
}

func TestLightspeedProviderSkeletonDoesNotCallExternalAPI(t *testing.T) {
	provider := LightspeedProvider{}

	response, err := provider.Chat(ProviderRequest{Message: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if response.RawProvider != "lightspeed" {
		t.Fatalf("provider = %q, want lightspeed", response.RawProvider)
	}
}

func TestLightspeedProviderPostsToConfiguredEndpoint(t *testing.T) {
	var body string
	var authorization string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		authorization = r.Header.Get("Authorization")
		buffer, _ := io.ReadAll(r.Body)
		body = string(buffer)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"answer":"provider answer"}`))
	}))
	defer server.Close()

	provider := LightspeedProvider{Config: LightspeedProviderConfig{EndpointURL: server.URL, Token: "secret-token"}}

	response, err := provider.Chat(ProviderRequest{Message: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if response.Answer != "provider answer" {
		t.Fatalf("answer = %q, want provider answer", response.Answer)
	}
	if !strings.Contains(body, `"message":"hello"`) {
		t.Fatalf("request body = %q, want message", body)
	}
	if strings.Contains(body, `"model"`) {
		t.Fatalf("request body = %q, did not expect model; model is selected by OLSConfig", body)
	}
	if strings.Contains(body, `"provider"`) {
		t.Fatalf("request body = %q, did not expect provider; provider is selected by OLSConfig", body)
	}
	if authorization != "Bearer secret-token" {
		t.Fatalf("authorization = %q, want bearer token", authorization)
	}
}

func TestLightspeedProviderReadsOpenAIStyleResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"openai style answer"}}]}`))
	}))
	defer server.Close()

	provider := LightspeedProvider{Config: LightspeedProviderConfig{EndpointURL: server.URL}}

	response, err := provider.Chat(ProviderRequest{Message: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if response.Answer != "openai style answer" {
		t.Fatalf("answer = %q, want openai style answer", response.Answer)
	}
}

func TestLightspeedProviderReadsResponseField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":"internal llm answer"}`))
	}))
	defer server.Close()

	provider := LightspeedProvider{Config: LightspeedProviderConfig{EndpointURL: server.URL}}

	response, err := provider.Chat(ProviderRequest{Message: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if response.Answer != "internal llm answer" {
		t.Fatalf("answer = %q, want internal llm answer", response.Answer)
	}
}
