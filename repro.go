package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	tempDir, _ := os.MkdirTemp("", "sync-repro-*")
	defer os.RemoveAll(tempDir)

	runCmd := func(dir string, name string, args ...string) (string, error) {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	remoteDir := filepath.Join(tempDir, "remote")
	os.Mkdir(remoteDir, 0755)
	runCmd(remoteDir, "git", "init", "--bare")

	localDir := filepath.Join(tempDir, "local")
	runCmd(tempDir, "git", "clone", remoteDir, "local")

	runCmd(localDir, "git", "config", "user.email", "test@example.com")
	runCmd(localDir, "git", "config", "user.name", "Test User")

	os.WriteFile(filepath.Join(localDir, "VERSION.md"), []byte("1.0.0"), 0644)
	os.WriteFile(filepath.Join(localDir, "README.md"), []byte("Initial content"), 0644)

	syncContent, _ := os.ReadFile("sync.sh")
	os.WriteFile(filepath.Join(localDir, "sync.sh"), syncContent, 0755)

	runCmd(localDir, "git", "add", ".")
	runCmd(localDir, "git", "commit", "-m", "Initial commit")
	runCmd(localDir, "git", "push", "origin", "master") // Test master

	runCmd(localDir, "git", "checkout", "-b", "feat/test-branch")
	os.WriteFile(filepath.Join(localDir, "feature.txt"), []byte("Feature content"), 0644)
	runCmd(localDir, "git", "add", ".")
	runCmd(localDir, "git", "commit", "-m", "Add feature")

	runCmd(localDir, "git", "checkout", "master")
	os.WriteFile(filepath.Join(localDir, "README.md"), []byte("Updated content on main"), 0644)
	runCmd(localDir, "git", "add", ".")
	runCmd(localDir, "git", "commit", "-m", "Update main")
	runCmd(localDir, "git", "push", "origin", "master")

	runCmd(localDir, "git", "checkout", "feat/test-branch")
	runCmd(localDir, "git", "checkout", "-b", "feat/ready-test")

	os.WriteFile(filepath.Join(localDir, "uncommitted.txt"), []byte("Stash me"), 0644)

	fmt.Println("Running sync.sh...")
	cmd := exec.Command("./sync.sh")
	cmd.Dir = localDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if _, err := os.Stat(filepath.Join(localDir, "uncommitted.txt")); os.IsNotExist(err) {
		fmt.Println("FAILURE: uncommitted.txt is missing!")
	} else {
		fmt.Println("SUCCESS: uncommitted.txt preserved.")
	}
}
