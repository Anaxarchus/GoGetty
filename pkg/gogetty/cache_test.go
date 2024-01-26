package gogetty

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCache(t *testing.T) {
	// Setup: Create a new temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gogetty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Ensure the temporary directory is removed after the test
	defer os.RemoveAll(tempDir)

	// Redirect the cache directory to the temporary directory
	originalHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomeDir)

	// Invoke the function
	if err := InitCache(); err != nil {
		t.Errorf("InitCache failed: %v", err)
	}

	// Assert: Check if cache directory and file were created
	cacheDir := filepath.Join(tempDir, ".gogetty_cache")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		t.Errorf("Cache directory was not created")
	}

	cacheFile := filepath.Join(cacheDir, "cache.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Errorf("Cache file was not created")
	}

	// Optional: Check the contents of the cache file
	file, err := os.Open(cacheFile)
	if err != nil {
		t.Fatalf("Failed to open cache file: %v", err)
	}
	defer file.Close()

	var cache Cache
	if err := json.NewDecoder(file).Decode(&cache); err != nil {
		t.Errorf("Cache file contains invalid JSON")
	}
}
