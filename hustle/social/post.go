package social

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
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

type TwitterProvider struct{}
func (p *TwitterProvider) Post(content string) error {
	fmt.Printf("[Twitter] Posting: %s\n", content)
	return nil
}

type LinkedInProvider struct{}
func (p *LinkedInProvider) Post(content string) error {
	fmt.Printf("[LinkedIn] Posting: %s\n", content)
	return nil
}

func GenerateContent(topic string) string {
	return fmt.Sprintf("Revolutionary insights on %s! #hustle #ai", topic)
}

func SchedulePost(orch *orchestrator.Orchestrator, provider Provider, platform, topic string) {
	content := GenerateContent(topic)
	fmt.Printf("Scheduling post for %s: %s\n", platform, content)

	err := provider.Post(content)
	if err == nil {
		orch.L1.Add(orchestrator.MemoryEntry{
			ID:        fmt.Sprintf("social-%s-%d", platform, time.Now().Unix()),
			Content:   fmt.Sprintf("Posted to %s: %s", platform, content),
			Timestamp: time.Now(),
			Tags:      []string{"social", platform},
		})

		// Log minimal expense for API usage
		orch.Ledger.Add(orchestrator.Transaction{
			Amount: 0.01,
			Type:   orchestrator.Expense,
			Hustle: "SocialMedia",
			Note:   fmt.Sprintf("API post to %s", platform),
		})
	}
}
