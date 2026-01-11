package llm

import (
	"os"
	"testing"
)

func TestNewLLM_ClaudeModel(t *testing.T) {
	originalKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Setenv("ANTHROPIC_API_KEY", originalKey)

	tests := []struct {
		name      string
		model     string
		wantError bool
	}{
		{
			name:      "claude sonnet model",
			model:     "claude-sonnet-4-20250514",
			wantError: false,
		},
		{
			name:      "claude opus model",
			model:     "claude-opus-4-5-20251101",
			wantError: false,
		},
		{
			name:      "claude prefix model",
			model:     "claude-any-version",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm, err := NewLLM(tt.model)

			if tt.wantError && err == nil {
				t.Errorf("NewLLM(%s) expected error, got nil", tt.model)
			}
			if !tt.wantError && err != nil {
				t.Errorf("NewLLM(%s) unexpected error: %v", tt.model, err)
			}
			if !tt.wantError && llm == nil {
				t.Errorf("NewLLM(%s) returned nil LLM", tt.model)
			}
		})
	}
}

func TestNewLLM_MissingAPIKey(t *testing.T) {
	originalKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Setenv("ANTHROPIC_API_KEY", originalKey)

	llm, err := NewLLM("claude-sonnet-4-20250514")

	if err == nil {
		t.Error("NewLLM() expected error when ANTHROPIC_API_KEY is not set")
	}
	if llm != nil {
		t.Error("NewLLM() expected nil LLM when API key is missing")
	}
	if err != nil && err.Error() != "env ANTHROPIC_API_KEY not set" {
		t.Errorf("NewLLM() wrong error message: %v", err)
	}
}

func TestNewLLM_UnsupportedModel(t *testing.T) {
	llm, err := NewLLM("unsupported-model")

	if llm != nil {
		t.Error("NewLLM() should return nil for unsupported model")
	}
	if err == nil {
		t.Error("NewLLM() should return error for unsupported model")
	}
}

func TestNewLLM_EmptyModel(t *testing.T) {
	llm, err := NewLLM("")

	if llm != nil {
		t.Error("NewLLM() should return nil for empty model")
	}
	if err == nil {
		t.Error("NewLLM() should return error for empty model")
	}
}

func TestNewLLM_GPTModel(t *testing.T) {
	llm, err := NewLLM("gpt-4")

	if llm != nil {
		t.Error("NewLLM() should return nil for unimplemented GPT model")
	}
	if err == nil {
		t.Error("NewLLM() should return error for unimplemented GPT model")
	}
}

func TestCreateMasterPrompt_EmptyInputs(t *testing.T) {
	signatures := make(map[string][]string)
	addedFiles := make(map[string]string)
	input := ""

	prompt := CreateMasterPrompt(signatures, addedFiles, input)

	if prompt == "" {
		t.Error("CreateMasterPrompt() returned empty string")
	}
}

func TestCreateMasterPrompt_WithSignatures(t *testing.T) {
	signatures := map[string][]string{
		"/path/to/file.go": {"func Foo()", "func Bar(x int)"},
	}
	addedFiles := make(map[string]string)
	input := "What does Foo do?"

	prompt := CreateMasterPrompt(signatures, addedFiles, input)

	if prompt == "" {
		t.Error("CreateMasterPrompt() returned empty string")
	}

	if !containsString(prompt, "/path/to/file.go") {
		t.Error("CreateMasterPrompt() missing file path")
	}
	if !containsString(prompt, "func Foo()") {
		t.Error("CreateMasterPrompt() missing signature")
	}
	if !containsString(prompt, "What does Foo do?") {
		t.Error("CreateMasterPrompt() missing input")
	}
}

func TestCreateMasterPrompt_WithAddedFiles(t *testing.T) {
	signatures := make(map[string][]string)
	addedFiles := map[string]string{
		"/path/to/file.go": "package main\n\nfunc main() {}",
	}
	input := "Explain this code"

	prompt := CreateMasterPrompt(signatures, addedFiles, input)

	if !containsString(prompt, "/path/to/file.go") {
		t.Error("CreateMasterPrompt() missing added file path")
	}
	if !containsString(prompt, "package main") {
		t.Error("CreateMasterPrompt() missing added file content")
	}
}

func TestCreateMasterPrompt_WithBothSignaturesAndFiles(t *testing.T) {
	signatures := map[string][]string{
		"/path/to/a.go": {"func A()"},
		"/path/to/b.go": {"func B()", "func C()"},
	}
	addedFiles := map[string]string{
		"/path/to/a.go": "package a\n\nfunc A() {}",
	}
	input := "How do A and B relate?"

	prompt := CreateMasterPrompt(signatures, addedFiles, input)

	if !containsString(prompt, "func A()") {
		t.Error("CreateMasterPrompt() missing signature A")
	}
	if !containsString(prompt, "func B()") {
		t.Error("CreateMasterPrompt() missing signature B")
	}
	if !containsString(prompt, "package a") {
		t.Error("CreateMasterPrompt() missing file content")
	}
	if !containsString(prompt, "How do A and B relate?") {
		t.Error("CreateMasterPrompt() missing input")
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
