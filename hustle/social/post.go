package social

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"os"
	"time"
)

type SocialPost struct {
	Platform    string
	Content     string
	ScheduledAt time.Time
}

type Provider interface {
	Post(content string) error
}

type AuthenticatedProvider struct {
	Name string
	Auth *OAuthState
}

func (p *AuthenticatedProvider) Post(content string) error {
	fmt.Printf("[%s] Posting: %s\n", p.Name, content)

	// Simulation of using the token
	if p.Auth != nil && p.Auth.Token != nil {
		fmt.Printf("[%s] Using Bearer Token: %s...\n", p.Name, p.Auth.Token.AccessToken[:10])
	} else {
		fmt.Printf("[%s] Warning: No authentication token found. Posting in mock mode.\n", p.Name)
	}

	return nil
}

func NewTwitterProvider(auth *OAuthState) *AuthenticatedProvider {
	return &AuthenticatedProvider{Name: "Twitter", Auth: auth}
}

func NewLinkedInProvider(auth *OAuthState) *AuthenticatedProvider {
	return &AuthenticatedProvider{Name: "LinkedIn", Auth: auth}
}

type SocialModule struct {
	Orchestrator *orchestrator.Orchestrator
	Providers    map[string]*AuthenticatedProvider
}

func NewSocialModule(orch *orchestrator.Orchestrator) *SocialModule {
	m := &SocialModule{
		Orchestrator: orch,
		Providers:    make(map[string]*AuthenticatedProvider),
	}

	// Initialize Twitter
	twitterAuth := NewOAuthState("Twitter",
		os.Getenv("TWITTER_CLIENT_ID"),
		os.Getenv("TWITTER_CLIENT_SECRET"),
		[]string{"tweet.read", "tweet.write", "users.read"},
		"https://twitter.com/i/oauth2/authorize",
		"https://api.twitter.com/2/oauth2/token")
	twitterAuth.Load()
	m.Providers["Twitter"] = NewTwitterProvider(twitterAuth)

	// Initialize LinkedIn
	linkedinAuth := NewOAuthState("LinkedIn",
		os.Getenv("LINKEDIN_CLIENT_ID"),
		os.Getenv("LINKEDIN_CLIENT_SECRET"),
		[]string{"w_member_social"},
		"https://www.linkedin.com/oauth/v2/authorization",
		"https://www.linkedin.com/oauth/v2/accessToken")
	linkedinAuth.Load()
	m.Providers["LinkedIn"] = NewLinkedInProvider(linkedinAuth)

	return m
}

func GenerateContent(orch *orchestrator.Orchestrator, topic string) string {
	prompt := fmt.Sprintf("Create a short, revolutionary social media post about %s with hashtags.", topic)
	content, err := orch.LLM.Generate(prompt)
	if err != nil {
		return fmt.Sprintf("Revolutionary insights on %s! #hustle #ai", topic)
	}
	return content
}

func SchedulePost(orch *orchestrator.Orchestrator, provider Provider, platform, topic string) {
	content := GenerateContent(orch, topic)
	fmt.Printf("Scheduling post for %s: %s\n", platform, content)

	err := provider.Post(content)
	if err == nil {
		orch.L1.Add(orchestrator.MemoryEntry{
			ID:        fmt.Sprintf("social-%s-%d", platform, time.Now().Unix()),
			Content:   fmt.Sprintf("Posted to %s: %s", platform, content),
			Timestamp: time.Now(),
			Tags:      []string{"social", platform},
		})

		orch.Ledger.Add(orchestrator.Transaction{
			Amount: 0.01,
			Type:   orchestrator.Expense,
			Hustle: "SocialMedia",
			Note:   fmt.Sprintf("API post to %s", platform),
		})
	}
}
