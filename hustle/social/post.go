package social

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"net/http"
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

type TwitterProvider struct {
	APIKey string
}

func NewTwitterProvider() *TwitterProvider {
	return &TwitterProvider{
		APIKey: os.Getenv("TWITTER_API_KEY"),
	}
}

func (p *TwitterProvider) Post(content string) error {
	if p.APIKey == "" {
		fmt.Printf("[Twitter Mock] Posting: %s\n", content)
		return nil
	}

	fmt.Printf("[Twitter API] Sending post: %s\n", content)

	payload := map[string]string{"text": content}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Real API logic would go here: client.Do(req)
	return nil
}

type LinkedInProvider struct {
	AccessToken string
}

func NewLinkedInProvider() *LinkedInProvider {
	return &LinkedInProvider{
		AccessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
	}
}

func (p *LinkedInProvider) Post(content string) error {
	if p.AccessToken == "" {
		fmt.Printf("[LinkedIn Mock] Posting: %s\n", content)
		return nil
	}

	fmt.Printf("[LinkedIn API] Sending post: %s\n", content)
	// Real API logic for LinkedIn
	return nil
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
			Amount: 0.05,
			Type:   orchestrator.Expense,
			Hustle: "SocialMedia",
			Note:   fmt.Sprintf("API post to %s", platform),
		})
	}
}
