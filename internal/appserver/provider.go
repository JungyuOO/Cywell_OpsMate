package appserver

type LightspeedProviderConfig struct {
	EndpointURL       string
	CredentialsSecret string
	DefaultProvider   string
	DefaultModel      string
}

type LightspeedProvider struct {
	Config LightspeedProviderConfig
}

func (p LightspeedProvider) Chat(request ProviderRequest) (ProviderResponse, error) {
	return ProviderResponse{
		Answer:      "Lightspeed provider client is not wired yet.",
		RawProvider: "lightspeed",
	}, nil
}

func MinimalProviderContext(citations []Citation) []ProviderContext {
	context := make([]ProviderContext, 0, len(citations))
	for _, citation := range citations {
		context = append(context, ProviderContext{
			Type:   "rag_chunk",
			Source: citation.DocumentID + "/" + citation.ChunkID,
		})
	}
	return context
}
