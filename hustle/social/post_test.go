package social

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/robertpelloni/hustle/orchestrator"
)

func TestLinkedInProvider_Post(t *testing.T) {
	// Helper function to mock orchestrator
	mockOrch := &orchestrator.Orchestrator{}

	t.Run("Missing Env Vars", func(t *testing.T) {
		os.Unsetenv("LINKEDIN_ACCESS_TOKEN")
		os.Unsetenv("LINKEDIN_MEMBER_ID")

		p := NewLinkedInProvider()
		err := p.Post(mockOrch, "LinkedIn", "Test content")
		if err == nil {
			t.Fatalf("expected error due to missing env vars, got nil")
		}
	})

	t.Run("Successful Post", func(t *testing.T) {
		t.Setenv("LINKEDIN_ACCESS_TOKEN", "mock-token")
		t.Setenv("LINKEDIN_MEMBER_ID", "mock-member-id")

		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer mock-token" {
				t.Errorf("Expected Bearer mock-token, got %s", authHeader)
			}
			protocolVersion := r.Header.Get("X-Restli-Protocol-Version")
			if protocolVersion != "2.0.0" {
				t.Errorf("Expected X-Restli-Protocol-Version 2.0.0, got %s", protocolVersion)
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id": "urn:li:share:12345"}`)
		}))
		defer server.Close()

		p := &LinkedInProvider{
			HTTPClient: server.Client(),
			APIURL:     server.URL,
		}

		err := p.Post(mockOrch, "LinkedIn", "Test content")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("API Error Response", func(t *testing.T) {
		t.Setenv("LINKEDIN_ACCESS_TOKEN", "mock-token")
		t.Setenv("LINKEDIN_MEMBER_ID", "mock-member-id")

		// Create a mock HTTP server returning 401 Unauthorized
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"message": "Unauthorized"}`)
		}))
		defer server.Close()

		p := &LinkedInProvider{
			HTTPClient: server.Client(),
			APIURL:     server.URL,
		}

		err := p.Post(mockOrch, "LinkedIn", "Test content")
		if err == nil {
			t.Fatalf("expected error from API response, got nil")
		}
	})
}
