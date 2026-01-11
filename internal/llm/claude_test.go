package llm

import (
	"os"
	"testing"

	"github.com/manosriram/wingman/internal/types"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.APIKey != apiKey {
		t.Errorf("NewClient() APIKey = %s, want %s", client.APIKey, apiKey)
	}
	if client.HTTPClient == nil {
		t.Error("NewClient() HTTPClient is nil")
	}
}

func TestNewClaudeLLM(t *testing.T) {
	originalKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Setenv("ANTHROPIC_API_KEY", originalKey)

	req := LLMRequest{
		Model:               "claude-sonnet-4-20250514",
		Input:               "test input",
		InputWithoutRepoMap: "test input without repo",
	}

	llm := NewClaudeLLM(req)

	if llm.SelectedModel != req.Model {
		t.Errorf("NewClaudeLLM() SelectedModel = %s, want %s", llm.SelectedModel, req.Model)
	}
	if llm.Input != req.Input {
		t.Errorf("NewClaudeLLM() Input = %s, want %s", llm.Input, req.Input)
	}
	if llm.InputWithoutRepoMap != req.InputWithoutRepoMap {
		t.Errorf("NewClaudeLLM() InputWithoutRepoMap = %s, want %s", llm.InputWithoutRepoMap, req.InputWithoutRepoMap)
	}
	if llm.Client == nil {
		t.Error("NewClaudeLLM() Client is nil")
	}
}

func TestClaudeLLM_GetMaxTokenCount(t *testing.T) {
	llm := ClaudeLLM{}

	maxTokens := llm.GetMaxTokenCount("any-model")

	if maxTokens != 200000 {
		t.Errorf("GetMaxTokenCount() = %d, want 200000", maxTokens)
	}
}

func TestClaudeLLM_GetSelectedModel(t *testing.T) {
	llm := ClaudeLLM{SelectedModel: "claude-sonnet-4-20250514"}

	model := llm.GetSelectedModel()

	if model != "claude-sonnet-4-20250514" {
		t.Errorf("GetSelectedModel() = %s, want claude-sonnet-4-20250514", model)
	}
}

func TestClaudeLLM_GetInputTokenCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "empty input",
			input: "",
			want:  1,
		},
		{
			name:  "single word",
			input: "hello",
			want:  1,
		},
		{
			name:  "multiple words",
			input: "hello world foo bar",
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := ClaudeLLM{Input: tt.input}
			got := llm.GetInputTokenCount()
			if got != tt.want {
				t.Errorf("GetInputTokenCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestResponse_GetTextResponse(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		want     string
	}{
		{
			name:     "empty content",
			response: Response{Content: []types.ContentBlock{}},
			want:     "",
		},
		{
			name: "single text block",
			response: Response{
				Content: []types.ContentBlock{
					{Type: "text", Text: "Hello world"},
				},
			},
			want: "Hello world",
		},
		{
			name: "multiple text blocks",
			response: Response{
				Content: []types.ContentBlock{
					{Type: "text", Text: "Hello "},
					{Type: "text", Text: "world"},
				},
			},
			want: "Hello world",
		},
		{
			name: "mixed content types",
			response: Response{
				Content: []types.ContentBlock{
					{Type: "text", Text: "Hello"},
					{Type: "image", Text: "ignored"},
					{Type: "text", Text: " world"},
				},
			},
			want: "Hello world",
		},
		{
			name: "non-text blocks only",
			response: Response{
				Content: []types.ContentBlock{
					{Type: "image", Text: "ignored"},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.GetTextResponse()
			if got != tt.want {
				t.Errorf("GetTextResponse() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLLMRequest_Struct(t *testing.T) {
	req := LLMRequest{
		Model:               "claude-sonnet-4-20250514",
		Input:               "test input",
		InputWithoutRepoMap: "test without repo",
	}

	if req.Model != "claude-sonnet-4-20250514" {
		t.Errorf("LLMRequest.Model = %s, want claude-sonnet-4-20250514", req.Model)
	}
	if req.Input != "test input" {
		t.Errorf("LLMRequest.Input = %s, want test input", req.Input)
	}
	if req.InputWithoutRepoMap != "test without repo" {
		t.Errorf("LLMRequest.InputWithoutRepoMap = %s, want test without repo", req.InputWithoutRepoMap)
	}
}

func TestLLMResponse_Struct(t *testing.T) {
	resp := LLMResponse{
		Response: "test response",
	}

	if resp.Response != "test response" {
		t.Errorf("LLMResponse.Response = %s, want test response", resp.Response)
	}
}
