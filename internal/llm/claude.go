package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	APIBaseURL = "https://api.anthropic.com/v1/messages"
	APIVersion = "2023-06-01"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type ErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type Response struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
}

type Request struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	TopK        int       `json:"top_k,omitempty"`
	System      string    `json:"system,omitempty"`
}

func (c *Client) SendMessage(req Request) (*Response, error) {
	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", APIBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", APIVersion)

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s - %s", errResp.Error.Type, errResp.Error.Message)
	}

	// Parse successful response
	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &apiResp, nil
}

func (c *Client) SendPrompt(prompt string) (string, error) {
	req := Request{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	resp, err := c.SendMessage(req)
	if err != nil {
		return "", err
	}

	return resp.GetTextResponse(), nil
}

func (r *Response) GetTextResponse() string {
	if len(r.Content) == 0 {
		return ""
	}
	var text strings.Builder
	for _, block := range r.Content {
		if block.Type == "text" {
			text.WriteString(block.Text)
		}
	}
	return text.String()
}

type ClaudeLLM struct {
	SelectedModel LLMModel
	Input         string
	Client        *Client
}

func NewClaudeLLM(model LLMModel, input string) ClaudeLLM {
	c := NewClient("")
	return ClaudeLLM{
		SelectedModel: model,
		Input:         input,
		Client:        c,
	}
}

func (c ClaudeLLM) GetMaxTokenCount(model string) int64 {
	return 200000
}

func (c ClaudeLLM) GetSelectedModel() LLMModel {
	return c.SelectedModel
}

func (c ClaudeLLM) GetInputTokenCount() int {
	return len(strings.Split(c.Input, " "))
}

func (c ClaudeLLM) Call() (LLMResponse, error) {
	response, err := c.Client.SendPrompt(c.Input)
	if err != nil {
		return LLMResponse{}, err
	}
	return LLMResponse{
		Response: response,
	}, nil

}
