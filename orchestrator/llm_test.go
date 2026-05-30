package orchestrator

import (
	"testing"
)

func TestLLMInterface(t *testing.T) {
	mock := &MockLLM{Response: "Custom response"}
	res, err := mock.Generate("Tell me a joke")
	if err != nil {
		t.Fatalf("LLM generation failed: %v", err)
	}
	if res != "Custom response" {
		t.Errorf("Expected 'Custom response', got %s", res)
	}
}

func TestOrchestratorLLMIntegration(t *testing.T) {
	orch := NewOrchestrator()
	if orch.LLM == nil {
		t.Fatal("Orchestrator LLM should not be nil")
	}

	res, _ := orch.LLM.Generate("Test")
	if res == "" {
		t.Error("Orchestrator LLM generated empty response")
	}
}
