package appserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LightspeedProviderConfig struct {
	EndpointURL       string
	CredentialsSecret string
	DefaultProvider   string
	DefaultModel      string
}

type HTTPDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

type LightspeedProvider struct {
	Config LightspeedProviderConfig
	Client HTTPDoer
}

func (p LightspeedProvider) Chat(request ProviderRequest) (ProviderResponse, error) {
	if p.Config.EndpointURL == "" {
		return ProviderResponse{
			Answer:      "Lightspeed provider client is not wired yet.",
			RawProvider: "lightspeed",
		}, nil
	}

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}

	body, err := json.Marshal(map[string]any{
		"message":        request.Message,
		"context":        request.Context,
		"clusterContext": request.ClusterContext,
		"provider":       p.Config.DefaultProvider,
		"model":          p.Config.DefaultModel,
	})
	if err != nil {
		return ProviderResponse{}, err
	}
	httpRequest, err := http.NewRequest(http.MethodPost, p.Config.EndpointURL, bytes.NewReader(body))
	if err != nil {
		return ProviderResponse{}, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return ProviderResponse{}, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		return ProviderResponse{}, fmt.Errorf("lightspeed provider returned status %d", httpResponse.StatusCode)
	}

	var response struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return ProviderResponse{}, err
	}
	return ProviderResponse{
		Answer:      response.Answer,
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
