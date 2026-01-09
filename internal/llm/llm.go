package llm

import (
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

const (
	OPUS_4_5 LLMModel = "claude_opus_4_5"
	GPT_5_2  LLMModel = "openai_gpt_5_2"
)

type LLMResponse struct {
	Response string
}

type LLM interface {
	GetMaxTokenCount(string) int64
	GetSelectedModel() LLMModel
	GetInputTokenCount() int
	Call() (LLMResponse, error)
}

func GetLLM(llmPlatform, input string, model LLMModel) LLM {
	switch llmPlatform {
	case string(CLAUDE):
		return NewClaudeLLM(LLMRequest{Model: model, Input: input})
	case string(OPENAI):
		break
	default:
		break
	}
	return ClaudeLLM{}
}

// TODO: add token count check
func CreateMasterPrompt(signatures map[string][]string, input string) string {
	var prompt strings.Builder
	prompt.WriteString(types.BASE_LLM_PROMPT)

	for filepath, signature := range signatures {
		prompt.WriteString(filepath + ": \n")
		for _, s := range signature {
			prompt.WriteString(s + "\n")
		}
		prompt.WriteString("\n")
	}
	return prompt.String() + "\n\n\n" + input
}
