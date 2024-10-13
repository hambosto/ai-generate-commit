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
	BaseURL     = "https://api.groq.com/openai/v1/chat/completions" // The base URL for the GROQ API
	contentType = "application/json"                                // The content type for API requests
)

// Message represents a single message in the conversation with the AI.
type Message struct {
	Role    string `json:"role"`    // The role of the sender (e.g., "user", "assistant")
	Content string `json:"content"` // The content of the message
}

// CompletionRequest holds the request payload sent to the API for generating a completion.
type CompletionRequest struct {
	Model    string    `json:"model"`    // The model to use for generating completions
	Messages []Message `json:"messages"` // The messages that make up the conversation context
}

// CompletionResponse represents the response payload from the API.
type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"` // The generated content from the AI
		} `json:"message"` // The message structure in the API response
	} `json:"choices"` // The list of choices returned by the API
}

// Client represents a GROQ API client.
type Client struct {
	httpClient *http.Client // The HTTP client used to make requests
	apiKey     string       // The API key for authenticating with the GROQ API
}

// NewClient creates a new GROQ API client.
// It retrieves the API key from the configuration and initializes the client with a timeout.
func NewClient() (*Client, error) {
	apiKey, err := config.GetConfig("GROQ_APIKEY")
	if err != nil {
		return nil, fmt.Errorf("failed to get GROQ_APIKEY: %w", err)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_APIKEY not set")
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second}, // Set a timeout for HTTP requests
		apiKey:     apiKey,                                  // Store the API key in the client
	}, nil
}

// GenerateCompletion sends a request to the GROQ API and returns the generated completion content.
// It takes a slice of messages that represents the conversation context and the model to be used.
func (c *Client) GenerateCompletion(messages []Message, model string) (string, error) {
	// Marshal the request body into JSON format
	reqBody, err := json.Marshal(CompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers for the request
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", contentType)

	// Send the request to the GROQ API
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the response status code indicates success
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body into the CompletionResponse struct
	var completionResp CompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if any completion choices were returned
	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	// Return the content of the first completion choice
	return completionResp.Choices[0].Message.Content, nil
}

