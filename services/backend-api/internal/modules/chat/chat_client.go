package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
			Timeout: 30 * time.Second,
		},
	}
}
type PythonInferenceRequest struct {
	Messages []Message `json:"messages"`
}

type PythonInferenceResponse struct {
	Response string `json:"response"`
}
func (c *HTTPInferenceClient) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	url := fmt.Sprintf("%s/api/v1/predict", c.baseURL)
	reqBody, err := json.Marshal(PythonInferenceRequest{Messages: history})
	if err != nil {
		return "", fmt.Errorf("ai client serialization failure: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to construct outgoing ai network request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("local model engine connection timed out or failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("model generation engine reported a non-200 fatal response code: %d", resp.StatusCode)
	}

	var inferenceResp PythonInferenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&inferenceResp); err != nil {
		return "", fmt.Errorf("failed to decode message response data stream from ai model: %w", err)
	}

	return inferenceResp.Response, nil
}