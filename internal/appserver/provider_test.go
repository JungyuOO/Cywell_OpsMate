package appserver

import "testing"

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
