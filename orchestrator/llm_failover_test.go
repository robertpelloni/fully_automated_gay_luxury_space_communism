package orchestrator

import (
	"testing"
)

func TestLLMFailover(t *testing.T) {
	primary := &MockLLM{Fail: true}
	secondary := &MockLLM{Response: "Secondary Success"}

	waterfall := &WaterfallLLM{
		Providers: []LLMProvider{primary, secondary},
	}

	res, err := waterfall.Generate("Test prompt")
	if err != nil {
		t.Fatalf("Failover failed: %v", err)
	}

	if res != "Secondary Success" {
		t.Errorf("Expected fallback response 'Secondary Success', got %s", res)
	}
}

func TestLLMAllFail(t *testing.T) {
	primary := &MockLLM{Fail: true}
	secondary := &MockLLM{Fail: true}

	waterfall := &WaterfallLLM{
		Providers: []LLMProvider{primary, secondary},
	}

	_, err := waterfall.Generate("Test prompt")
	if err == nil {
		t.Fatal("Expected error when all providers fail, got nil")
	}
}
