package orchestrator

import (
	"fmt"
	"os/exec"
)

type RollbackHandler struct {
	Orchestrator *Orchestrator
}

func NewRollbackHandler(o *Orchestrator) *RollbackHandler {
	return &RollbackHandler{Orchestrator: o}
}

func (r *RollbackHandler) Execute() error {
	fmt.Println("Executing automated rollback to previous stable state...")

	// 1. Abort any ongoing merge
	fmt.Println("Aborting any ongoing merges...")
	exec.Command("git", "merge", "--abort").Run()

	// 2. Reset to HEAD (clean up partial changes)
	fmt.Println("Resetting to HEAD...")
	out, err := exec.Command("git", "reset", "--hard", "HEAD").CombinedOutput()
	if err != nil {
		return fmt.Errorf("git reset failed: %v, output: %s", err, string(out))
	}

	// 3. Clean untracked files
	fmt.Println("Cleaning untracked files...")
	exec.Command("git", "clean", "-fd").Run()

	return nil
}
