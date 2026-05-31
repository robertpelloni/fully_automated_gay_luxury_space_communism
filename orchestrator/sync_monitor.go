package orchestrator

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type SyncHealth struct {
	LastSyncTime time.Time
	Status       string
	IssueCount   int
	GitState     string
}

func GetSyncHealth() SyncHealth {
	state := "Clean"
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		state = "Error fetching status"
	} else if len(out) > 0 {
		state = "Uncommitted Changes"
	}

	return SyncHealth{
		LastSyncTime: time.Now(),
		Status:       "Active",
		IssueCount:   0,
		GitState:     state,
	}
}

func MonitorSync() {
	health := GetSyncHealth()
	fmt.Printf("Sync Status: %s, Git State: %s, Last Check: %v\n", health.Status, health.GitState, health.LastSyncTime)
}

func CheckSubmodules() string {
	out, err := exec.Command("git", "submodule", "status").Output()
	if err != nil {
		return "Error checking submodules"
	}
	if len(out) == 0 {
		return "No submodules found"
	}
	return strings.TrimSpace(string(out))
}
