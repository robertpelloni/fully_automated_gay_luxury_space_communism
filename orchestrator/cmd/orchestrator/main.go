package main

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"os"
	"strings"
)

func main() {
	fmt.Println("=== AI Hustle Machine Orchestrator ===")

	// Source version from VERSION.md
	versionData, err := os.ReadFile("VERSION.md")
	version := "unknown"
	if err == nil {
		version = strings.TrimSpace(string(versionData))
	}

	orch := orchestrator.NewOrchestrator()

	// Real-time status reporting
	orchestrator.WriteStatusReport(version, "Running", "Orchestrator started successfully")

	fmt.Printf("Orchestrator (v%s) initialized with L1: %d, L2: %d, L3: %d entries.\n",
		version, len(orch.L1.Entries), len(orch.L2.Entries), len(orch.L3.Entries))
}
