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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateCompletion(messages []Message, model string) (string, error) {
	apiKey := config.GetConfig("GROQ_APIKEY")
	if len(apiKey) == 0 {
		return "", fmt.Errorf("GROQ_APIKEY not set")
	}

	reqBody, err := json.Marshal(CompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", nil
	}

	req, err := http.NewRequest("POST", BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", nil
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	var completionResp CompletionResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return "", err
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return completionResp.Choices[0].Message.Content, nil
}
