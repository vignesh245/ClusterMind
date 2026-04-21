package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vignesh245/ClusterMind/internal/ai"
)

// OllamaProvider implements the ai.Provider interface for local Ollama.
type OllamaProvider struct {
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}
	return &OllamaProvider{
		BaseURL: baseURL,
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.BaseURL+"/api/version", nil)
	if err != nil {
		return err
	}
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama health check failed with status: %d", resp.StatusCode)
	}
	return nil
}

type ollamaRequest struct {
	Model   string       `json:"model"`
	Prompt  string       `json:"prompt"`
	System  string       `json:"system"`
	Format  string       `json:"format"`
	Stream  bool         `json:"stream"`
	Options ollamaOptions `json:"options"`
}

type ollamaOptions struct {
	Temperature float32 `json:"temperature"`
	NumPredict  int     `json:"num_predict"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (p *OllamaProvider) Complete(ctx context.Context, req ai.CompletionRequest) (ai.CompletionResponse, error) {
	payload := ollamaRequest{
		Model:  p.Model,
		System: req.SystemPrompt,
		Prompt: req.UserPrompt,
		Format: "json",
		Stream: false,
		Options: ollamaOptions{
			Temperature: req.Temperature,
			NumPredict:  req.MaxTokens,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return ai.CompletionResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return ai.CompletionResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.HTTPClient.Do(httpReq)
	if err != nil {
		return ai.CompletionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return ai.CompletionResponse{}, fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return ai.CompletionResponse{}, err
	}

	return ai.CompletionResponse{
		Content:      ollamaResp.Response,
		FinishReason: "stop",
	}, nil
}
