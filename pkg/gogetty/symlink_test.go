package gogetty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateSymlink(t *testing.T) {
	// Setup: Create temporary directories and files to simulate the cache and project directories
	tempDir, err := os.MkdirTemp("", "gogetty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Ensure the temporary directory is removed after the test
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")
	sourceFile := filepath.Join(sourceDir, "module.txt")
	targetLink := filepath.Join(targetDir, "module-link.txt")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	if err := os.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test creating the symlink
	if err := CreateSymlink(sourceFile, targetLink); err != nil {
		t.Errorf("Failed to create symlink: %v", err)
	}

	// Verify that the symlink points to the correct source
	resolved, err := os.Readlink(targetLink)
	if err != nil {
		t.Errorf("Failed to read symlink: %v", err)
	}

	if resolved != sourceFile {
		t.Errorf("Symlink does not point to the correct source: got %v, want %v", resolved, sourceFile)
	}

	// Verify the symlink points to a file with the correct content
	linkContent, err := os.ReadFile(targetLink)
	if err != nil {
		t.Errorf("Failed to read from symlink: %v", err)
	}

	if string(linkContent) != "test content" {
		t.Errorf("Symlink does not point to a file with correct content: got %s, want %s", linkContent, "test content")
	}
}
