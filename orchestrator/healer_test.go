package orchestrator

import (
	"testing"
)

func TestHealerLoop(t *testing.T) {
	orch := NewOrchestrator()
	healer := NewHealer(orch)

	issue := "Test failure simulation"
	healer.Loop(issue)

	// In the mock, retryCount > 1 triggers verification success.
	// RetryCount starts at 0, Fix increments it.
	// Attempt 1: Fix (RetryCount=1), Verify fails (RetryCount !> 1).
	// Attempt 2: Fix (RetryCount=2), Verify passes (RetryCount > 1).
	if healer.RetryCount != 2 {
		t.Errorf("Expected 2 retries to succeed, got %d", healer.RetryCount)
	}

	if len(orch.L2.Entries) != 1 {
		t.Errorf("Expected resolution to be logged in L2 memory, got %d entries", len(orch.L2.Entries))
	}
}
