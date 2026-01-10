package llm

import (
	"fmt"
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
	Call() (LLMResponse, error)
	WriteToHistory(request string, response LLMResponse) error
}

func GetLLM(input, model string) LLM {
	if strings.HasPrefix(string(model), "claude") {
		return NewClaudeLLM(LLMRequest{Model: model, Input: input})
	} else if strings.HasPrefix(string(model), "gpt") {
		// openai
	}

	return ClaudeLLM{}
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
