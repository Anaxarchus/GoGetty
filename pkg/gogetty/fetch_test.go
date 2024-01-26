package gogetty

import (
	"os/exec"
	"testing"
)

func TestGitInstallation(t *testing.T) {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Git is not installed or not found in PATH")
	}
}

// Test
func TestConstructGitCommand(t *testing.T) {
	gitURL := "https://example.com/repo.git"
	branch := "test-branch"
	commit := "123abc"

	expected := "git clone --depth 1 --branch test-branch https://example.com/repo.git"
	result := ConstructGitCommand(gitURL, branch, commit)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Add more test cases as needed
}
