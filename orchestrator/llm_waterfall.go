package orchestrator

import (
	"fmt"
)

type WaterfallLLM struct {
	Providers []LLMProvider
}

func (w *WaterfallLLM) Generate(prompt string) (string, error) {
	var lastErr error
	for _, p := range w.Providers {
		res, err := p.Generate(prompt)
		if err == nil {
			return res, nil
		}
		lastErr = err
		fmt.Printf("LLM Provider failed, falling back: %v\n", err)
	}
	return "", fmt.Errorf("all LLM providers failed: %v", lastErr)
}
