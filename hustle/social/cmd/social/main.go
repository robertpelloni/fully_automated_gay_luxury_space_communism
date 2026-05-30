package main

import (
	"fmt"
	"github.com/robertpelloni/hustle/hustle/social"
	"github.com/robertpelloni/hustle/orchestrator"
)

func main() {
	fmt.Println("=== Social Media Hustle Module ===")
	orch := orchestrator.NewOrchestrator()

	provider := &social.TwitterProvider{}
	social.SchedulePost(orch, provider, "Twitter", "AI Agent Orchestration")

	fmt.Printf("Module execution complete. Current Profit: $%.2f\n", orch.Ledger.Profit())
}
