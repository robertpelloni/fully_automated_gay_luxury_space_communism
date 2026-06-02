package social

import (
	"context"
	"github.com/robertpelloni/hustle/orchestrator"
	"testing"
)

func TestSocialOAuthScaffold(t *testing.T) {
	auth := NewOAuthState("Test", "client-id", "client-secret", []string{"scope"}, "https://auth.com", "https://token.com")
	if auth.Provider != "Test" {
		t.Errorf("Expected provider Test, got %s", auth.Provider)
	}

	// Test token retrieval failure (expected)
	_, err := auth.GetToken(context.Background())
	if err == nil {
		t.Errorf("Expected error for missing token")
	}
}

func TestSocialModule(t *testing.T) {
	orch := orchestrator.NewOrchestrator()
	mod := NewSocialModule(orch)

	if _, ok := mod.Providers["Twitter"]; !ok {
		t.Errorf("Twitter provider not initialized")
	}

	if _, ok := mod.Providers["LinkedIn"]; !ok {
		t.Errorf("LinkedIn provider not initialized")
	}
}

func TestSchedulePost(t *testing.T) {
	orch := orchestrator.NewOrchestrator()
	auth := &OAuthState{Provider: "Test"}
	provider := NewTwitterProvider(auth)

	SchedulePost(orch, provider, "Twitter", "AI Agent Testing")

	if len(orch.L1.Entries) == 0 {
		t.Errorf("Post entry not saved to memory")
	}
}
