package orchestrator

import (
	"fmt"
)

type LLMProvider interface {
	Generate(prompt string) (string, error)
}

type MockLLM struct {
	Response string
}

func (m *MockLLM) Generate(prompt string) (string, error) {
	fmt.Printf("Mock LLM generating for prompt: %s\n", prompt)
	if m.Response != "" {
		return m.Response, nil
	}
	return "Generated content based on " + prompt, nil
}
