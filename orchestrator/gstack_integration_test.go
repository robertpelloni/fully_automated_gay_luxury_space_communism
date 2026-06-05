package orchestrator

import (
	"testing"
	"os/exec"
	"os"
)

func TestGStackSubmoduleExistence(t *testing.T) {
	wd, _ := os.Getwd()
	t.Logf("Working directory: %s", wd)
	// Tests run from the package directory, so gstack is one level up
	cmd := exec.Command("ls", "-d", "../gstack")
	err := cmd.Run()
	if err != nil {
		t.Errorf("gstack submodule not found at ../gstack: %v", err)
	}
}
