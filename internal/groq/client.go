package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hambosto/ai-generate-commit/internal/config"
)

const (
	BaseURL     = "https://api.groq.com/openai/v1/chat/completions"
	contentType = "application/json"
)

// Message represents a single message in the conversation with the AI.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest holds the request payload sent to the API for generating a completion.
type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// CompletionResponse represents the response payload from the API.
type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Client represents a GROQ API client.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new GROQ API client.
func NewClient() (*Client, error) {
	apiKey, err := config.GetConfig("GROQ_APIKEY")
	if err != nil {
		return nil, fmt.Errorf("failed to get GROQ_APIKEY: %w", err)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_APIKEY not set")
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
	}, nil
}

// GenerateCompletion sends a request to the GROQ API and returns the generated completion content.
func (c *Client) GenerateCompletion(messages []Message, model string) (string, error) {
	reqBody, err := json.Marshal(CompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var completionResp CompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return completionResp.Choices[0].Message.Content, nil
}

