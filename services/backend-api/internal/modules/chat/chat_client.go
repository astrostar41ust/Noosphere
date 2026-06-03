package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strings"
)

type AIClient interface {
	GenerateResponse(ctx context.Context, history []Message) (string, error)
}

type HTTPInferenceClient struct {
	baseURL    string
	httpClient *http.Client
}
func NewHTTPInferenceClient(baseURL string) *HTTPInferenceClient {
	return &HTTPInferenceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Minute,
		},
	}
}
type PythonInferenceRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type PythonInferenceResponse struct {
	Response string `json:"response"`
}

func (c *HTTPInferenceClient) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	url := fmt.Sprintf("%s/api/v1/chat", c.baseURL)

	var sb strings.Builder
	for _, msg := range history {
		if strings.ToLower(msg.Role) == "user" {
			sb.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		} else {
			sb.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		}
	}
	sb.WriteString("Assistant: ")

	payload := PythonInferenceRequest{
		Prompt:      sb.String(),
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("ai client serialization failure: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to construct outgoing request object: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("local model connection timed out: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("model generation engine reported status code: %d", resp.StatusCode)
	}

	var inferenceResp PythonInferenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&inferenceResp); err != nil {
		return "", fmt.Errorf("failed to decode text stream from ai model: %w", err)
	}

	return inferenceResp.Response, nil
}