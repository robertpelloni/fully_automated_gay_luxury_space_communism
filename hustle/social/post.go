package social

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"time"
)

type SocialPost struct {
	Platform string
	Content  string
	ScheduledAt time.Time
}

func GenerateContent(topic string) string {
	return fmt.Sprintf("Revolutionary insights on %s! #hustle #ai", topic)
}

func SchedulePost(orch *orchestrator.Orchestrator, platform, topic string) {
	content := GenerateContent(topic)
	fmt.Printf("Scheduling post for %s: %s\n", platform, content)

	orch.L1.Add(orchestrator.MemoryEntry{
		ID:        fmt.Sprintf("social-%s-%d", platform, time.Now().Unix()),
		Content:   fmt.Sprintf("Posted to %s: %s", platform, content),
		Timestamp: time.Now(),
		Tags:      []string{"social", platform},
	})
}
