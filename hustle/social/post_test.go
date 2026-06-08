package social

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/robertpelloni/hustle/orchestrator"
)

func TestTwitterProvider_DryRun(t *testing.T) {
	provider := &TwitterProvider{
		DryRun: true,
	}

	orch := orchestrator.NewOrchestrator()
	err := provider.Post(orch, "Twitter", "Test Content")
	if err != nil {
		t.Fatalf("expected no error during dry run, got: %v", err)
	}
}

func TestTwitterProvider_MissingEnv(t *testing.T) {
	provider := &TwitterProvider{
		DryRun: false,
	}

	os.Clearenv()

	orch := orchestrator.NewOrchestrator()
	err := provider.Post(orch, "Twitter", "Test Content")
	if err == nil {
		t.Fatalf("expected error due to missing env variables, got nil")
	}
	if err.Error() != "missing Twitter OAuth environment variables" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestTwitterProvider_PostSuccess(t *testing.T) {
	t.Setenv("TWITTER_CONSUMER_KEY", "test_key")
	t.Setenv("TWITTER_CONSUMER_SECRET", "test_secret")
	t.Setenv("TWITTER_ACCESS_TOKEN", "test_token")
	t.Setenv("TWITTER_ACCESS_SECRET", "test_token_secret")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var reqData twitterPostRequest
		json.Unmarshal(body, &reqData)
		if reqData.Text != "Test Content" {
			t.Errorf("expected text 'Test Content', got '%s'", reqData.Text)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	provider := &TwitterProvider{
		DryRun: false,
		APIURL: ts.URL,
	}

	orch := orchestrator.NewOrchestrator()
	err := provider.Post(orch, "Twitter", "Test Content")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}
