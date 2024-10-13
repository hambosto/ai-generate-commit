package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hambosto/ai-generate-commit/internal/config"
)

const BaseURL = "https://api.groq.com/openai/v1/chat/completions"

// Message represents a single message in the conversation with the AI.
type Message struct {
	Role    string `json:"role"`    // Role could be "user" or "assistant"
	Content string `json:"content"` // The content of the message
}

// CompletionRequest holds the request payload sent to the API for generating a completion.
type CompletionRequest struct {
	Model    string    `json:"model"`    // The model ID to use, e.g., "gpt-3.5-turbo"
	Messages []Message `json:"messages"` // A list of messages in the conversation
}

// CompletionResponse represents the response payload from the API.
type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"` // The content of the generated completion
		} `json:"message"`
	} `json:"choices"` // A list of possible completions (choices)
}

// GenerateCompletion takes a list of messages and a model ID, sends the data to the API, and returns the generated completion content.
func GenerateCompletion(messages []Message, model string) (string, error) {
	// Fetch the API key from the config
	apiKey := config.GetConfig("GROQ_APIKEY")
	if len(apiKey) == 0 {
		return "", fmt.Errorf("GROQ_APIKEY not set")
	}

	// Prepare the request body with the model and messages
	reqBody, err := json.Marshal(CompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers, including the Authorization header with the API key
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into CompletionResponse
	var completionResp CompletionResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if any completions were returned
	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	// Return the first completion message's content
	return completionResp.Choices[0].Message.Content, nil
}
