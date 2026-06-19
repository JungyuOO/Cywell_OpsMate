package appserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LightspeedProviderConfig struct {
	EndpointURL       string
	CredentialsSecret string
	Token             string
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
		"query":          request.Message,
		"context":        request.Context,
		"clusterContext": request.ClusterContext,
	})
	if err != nil {
		return ProviderResponse{}, err
	}
	httpRequest, err := http.NewRequest(http.MethodPost, p.Config.EndpointURL, bytes.NewReader(body))
	if err != nil {
		return ProviderResponse{}, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	if p.Config.Token != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+p.Config.Token)
	}

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return ProviderResponse{}, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		return ProviderResponse{}, fmt.Errorf("lightspeed provider returned status %d", httpResponse.StatusCode)
	}

	answer, err := extractProviderAnswer(httpResponse.Body)
	if err != nil {
		return ProviderResponse{}, err
	}
	return ProviderResponse{
		Answer:      answer,
		RawProvider: "lightspeed",
	}, nil
}

func extractProviderAnswer(body io.Reader) (string, error) {
	var payload any
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return "", err
	}
	answer := findProviderAnswer(payload)
	if answer == "" {
		return "", fmt.Errorf("lightspeed provider response did not include answer text")
	}
	return answer, nil
}

func findProviderAnswer(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case []any:
		for _, item := range typed {
			if answer := findProviderAnswer(item); answer != "" {
				return answer
			}
		}
	case map[string]any:
		for _, key := range []string{"answer", "response", "output", "content", "text", "generated_text"} {
			if answer := findProviderAnswer(typed[key]); answer != "" {
				return answer
			}
		}
		for _, key := range []string{"message", "data", "result", "prediction"} {
			if answer := findProviderAnswer(typed[key]); answer != "" {
				return answer
			}
		}
		for _, key := range []string{"choices", "outputs", "predictions"} {
			if answer := findProviderAnswer(typed[key]); answer != "" {
				return answer
			}
		}
	}
	return ""
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
