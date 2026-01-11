package llm

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manosriram/wingman/internal/types"
)

type LLMFamily string
type LLMModel string

const (
	CLAUDE LLMFamily = "claude"
	OPENAI LLMFamily = "openai"
	GEMINI LLMFamily = "gemini"
)

type LLMResponse struct {
	Response string
}

type LLM interface {
	GetMaxTokenCount(string) int64
	GetSelectedModel() string
	GetInputTokenCount() int
	Call(string) (*LLMResponse, error)
	WriteToHistory(request string, response *LLMResponse) error
}

func NewLLM(model string) (LLM, error) {
	if model == "" {
		return nil, errors.New("model cannot be empty")
	}

	if strings.HasPrefix(model, "claude") {
		if os.Getenv("ANTHROPIC_API_KEY") == "" {
			return nil, errors.New("env ANTHROPIC_API_KEY not set")
		}
		return NewClaudeLLM(LLMRequest{Model: model}), nil
	} else if strings.HasPrefix(model, "gpt") {
		return nil, errors.New("OpenAI models not yet implemented")
	}

	return nil, errors.New("unsupported model: " + model)
}

// TODO: add token count check
func CreateMasterPrompt(signatures map[string][]string, addedFiles map[string]string, input string) string {
	var prompt strings.Builder
	prompt.WriteString(types.BASE_LLM_PROMPT)

	for filepath, signature := range signatures {
		prompt.WriteString(filepath + ": \n")
		for _, s := range signature {
			prompt.WriteString(s + "\n")
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("\n")
	for path, content := range addedFiles {
		fmt.Fprintf(&prompt, "%s : %s", path, content)
		prompt.WriteString("\n")
	}
	prompt.WriteString("\n")

	return prompt.String() + "\n\n\n" + "Now answer the below question keeping in mind the above context\n\n" + input

}
