package main

import (
	"fmt"
	"github.com/robertpelloni/hustle/hustle/social"
	"github.com/robertpelloni/hustle/orchestrator"
)

func main() {
	fmt.Println("=== Social Media Hustle Module ===")
	orch := orchestrator.NewOrchestrator()

	social.SchedulePost(orch, "Twitter", "AI Agent Orchestration")
	fmt.Println("Module execution complete.")
}
