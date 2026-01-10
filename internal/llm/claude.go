package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/manosriram/wingman/internal/types"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

type Response struct {
	ID           string               `json:"id"`
	Type         string               `json:"type"`
	Role         string               `json:"role"`
	Content      []types.ContentBlock `json:"content"`
	Model        string               `json:"model"`
	StopReason   string               `json:"stop_reason"`
	StopSequence string               `json:"stop_sequence"`
	Usage        types.Usage          `json:"usage"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) SendMessage(req types.Request) (*Response, error) {
	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", types.APIBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", types.APIVersion)

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
		var errResp types.ErrorResponse
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

// TODO: read stream response
func (c *Client) SendPrompt(prompt string) (string, error) {
	req := types.Request{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		Messages: []types.Message{
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
	SelectedModel       string
	Input               string
	InputWithoutRepoMap string
	Client              *Client
}

type LLMRequest struct {
	Model               string
	Input               string
	InputWithoutRepoMap string
}

func NewClaudeLLM(req LLMRequest) ClaudeLLM {
	anthropicApiKey := os.Getenv("ANTHROPIC_API_KEY")
	c := NewClient(anthropicApiKey)

	return ClaudeLLM{
		SelectedModel:       req.Model,
		Input:               req.Input,
		InputWithoutRepoMap: req.InputWithoutRepoMap,
		Client:              c,
	}
}

func (c ClaudeLLM) GetMaxTokenCount(model string) int64 {
	return 200000
}

func (c ClaudeLLM) GetSelectedModel() string {
	return c.SelectedModel
}

func (c ClaudeLLM) GetInputTokenCount() int {
	return len(strings.Split(c.Input, " "))
}

func (c ClaudeLLM) WriteToHistory(request string, response LLMResponse) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(wd+"/.wingman.history.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content := fmt.Sprintf("%s\n%s\n\n\n", request, response.Response)
	if _, err = f.WriteString(content); err != nil {
		panic(err)
	}

	return nil
}

func (c ClaudeLLM) Call() (LLMResponse, error) {
	wd, err := os.Getwd()
	if err != nil {
		return LLMResponse{}, err
	}

	f, err := os.OpenFile(wd+"/wingman.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)

	}
	defer f.Close()

	response, err := c.Client.SendPrompt(c.Input)
	if _, err = f.WriteString(response); err != nil {
		panic(err)
	}

	return LLMResponse{
		Response: response,
	}, nil
}
