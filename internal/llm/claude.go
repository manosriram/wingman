package llm

import "strings"

type ClaudeLLM struct {
	SelectedModel LLMModel
	Input         string
}

func NewClaudeLLM(model LLMModel, input string) ClaudeLLM {
	return ClaudeLLM{
		SelectedModel: model,
		Input:         input,
	}
}

func (c *ClaudeLLM) GetMaxTokenCount() int64 {
	return 200000
}

func (c *ClaudeLLM) GetSelectedModel() LLMModel {
	return c.SelectedModel
}

func (c *ClaudeLLM) GetInputTokenCount() int {
	return len(strings.Split(c.Input, " "))
}

func (c *ClaudeLLM) Call() LLMResponse {
	return LLMResponse{}
}
