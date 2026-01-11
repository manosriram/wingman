package shell

import (
	"os"
	"testing"

	"github.com/manosriram/wingman/internal/llm"
)

func TestProgramFlags_Struct(t *testing.T) {
	model := "test-model"
	flags := ProgramFlags{
		Model: &model,
	}

	if *flags.Model != "test-model" {
		t.Errorf("ProgramFlags.Model = %s, want test-model", *flags.Model)
	}
}

func TestShell_Struct(t *testing.T) {
	model := "test-model"
	shell := Shell{
		ShellDir: "/test/dir",
		Flags: ProgramFlags{
			Model: &model,
		},
		Repository: nil,
		LLM:        nil,
	}

	if shell.ShellDir != "/test/dir" {
		t.Errorf("Shell.ShellDir = %s, want /test/dir", shell.ShellDir)
	}
	if *shell.Flags.Model != "test-model" {
		t.Errorf("Shell.Flags.Model = %s, want test-model", *shell.Flags.Model)
	}
}

func TestCmdChannel_Struct(t *testing.T) {
	tests := []struct {
		name     string
		channel  CmdChannel
		wantResp string
		wantErr  bool
	}{
		{
			name: "success response",
			channel: CmdChannel{
				Response: "test response",
				Error:    nil,
			},
			wantResp: "test response",
			wantErr:  false,
		},
		{
			name: "error response",
			channel: CmdChannel{
				Response: "",
				Error:    os.ErrNotExist,
			},
			wantResp: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.channel.Response != tt.wantResp {
				t.Errorf("CmdChannel.Response = %s, want %s", tt.channel.Response, tt.wantResp)
			}
			if (tt.channel.Error != nil) != tt.wantErr {
				t.Errorf("CmdChannel.Error = %v, wantErr %v", tt.channel.Error, tt.wantErr)
			}
		})
	}
}

type MockLLM struct {
	CallResponse    string
	CallError       error
	SelectedModel   string
	MaxTokenCount   int64
	InputTokenCount int
	WriteHistoryErr error
}

func (m *MockLLM) GetMaxTokenCount(model string) int64 {
	return m.MaxTokenCount
}

func (m *MockLLM) GetSelectedModel() string {
	return m.SelectedModel
}

func (m *MockLLM) GetInputTokenCount() int {
	return m.InputTokenCount
}

func (m *MockLLM) Call(prompt string) (*llm.LLMResponse, error) {
	if m.CallError != nil {
		return nil, m.CallError
	}
	return &llm.LLMResponse{Response: m.CallResponse}, nil
}

func (m *MockLLM) WriteToHistory(request string, response *llm.LLMResponse) error {
	return m.WriteHistoryErr
}

func setupTestDir(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "shell-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	goModContent := `module testmodule

go 1.21
`
	if err := os.WriteFile(tmpDir+"/go.mod", []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	goFileContent := `package main

func main() {
	println("hello")
}
`
	if err := os.WriteFile(tmpDir+"/main.go", []byte(goFileContent), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}

	return tmpDir
}

func cleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	os.RemoveAll(dir)
}

func TestMockLLM_Call(t *testing.T) {
	mock := &MockLLM{
		CallResponse:  "test response",
		SelectedModel: "test-model",
	}

	resp, err := mock.Call("test prompt")

	if err != nil {
		t.Errorf("MockLLM.Call() unexpected error: %v", err)
	}
	if resp.Response != "test response" {
		t.Errorf("MockLLM.Call() response = %s, want test response", resp.Response)
	}
}

func TestMockLLM_CallError(t *testing.T) {
	mock := &MockLLM{
		CallError: os.ErrNotExist,
	}

	resp, err := mock.Call("test prompt")

	if err == nil {
		t.Error("MockLLM.Call() expected error")
	}
	if resp != nil {
		t.Error("MockLLM.Call() expected nil response on error")
	}
}

func TestMockLLM_GetSelectedModel(t *testing.T) {
	mock := &MockLLM{
		SelectedModel: "claude-sonnet-4-20250514",
	}

	model := mock.GetSelectedModel()

	if model != "claude-sonnet-4-20250514" {
		t.Errorf("MockLLM.GetSelectedModel() = %s, want claude-sonnet-4-20250514", model)
	}
}

func TestMockLLM_GetMaxTokenCount(t *testing.T) {
	mock := &MockLLM{
		MaxTokenCount: 100000,
	}

	count := mock.GetMaxTokenCount("any")

	if count != 100000 {
		t.Errorf("MockLLM.GetMaxTokenCount() = %d, want 100000", count)
	}
}

func TestMockLLM_GetInputTokenCount(t *testing.T) {
	mock := &MockLLM{
		InputTokenCount: 500,
	}

	count := mock.GetInputTokenCount()

	if count != 500 {
		t.Errorf("MockLLM.GetInputTokenCount() = %d, want 500", count)
	}
}

func TestMockLLM_WriteToHistory(t *testing.T) {
	mock := &MockLLM{}

	err := mock.WriteToHistory("request", &llm.LLMResponse{Response: "response"})

	if err != nil {
		t.Errorf("MockLLM.WriteToHistory() unexpected error: %v", err)
	}
}

func TestMockLLM_WriteToHistoryError(t *testing.T) {
	mock := &MockLLM{
		WriteHistoryErr: os.ErrPermission,
	}

	err := mock.WriteToHistory("request", &llm.LLMResponse{Response: "response"})

	if err == nil {
		t.Error("MockLLM.WriteToHistory() expected error")
	}
}
