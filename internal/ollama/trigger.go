package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Model   string
	HTTP    *http.Client
}

func NewClient(baseURL, model string) *Client {
	return &Client{
		BaseURL: baseURL,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 60 * time.Second, // gerar resposta pode demorar
		},
	}
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Generate manda o prompt pronto (string) e devolve a resposta do modelo.
// Não sabe nada sobre tema, template, ou de onde o prompt veio.
func (c *Client) Generate(prompt string) (string, error) {
	reqBody := generateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("montar requisição: %w", err)
	}

	resp, err := c.HTTP.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("chamar ollama: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ler resposta: %w", err)
	}

	var out generateResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("decodificar resposta: %w", err)
	}

	return out.Response, nil
}
