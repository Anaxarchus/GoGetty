package gogetty

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestInitProject(t *testing.T) {
	// Setup: Ensure no .gogetty file exists
	_ = os.Remove(".gogetty")

	// Invoke the function
	InitProject()

	// Assert: Check if .gogetty file was created
	if _, err := os.Stat(".gogetty"); os.IsNotExist(err) {
		t.Errorf(".gogetty was not created as expected")
	}

	// Cleanup: Remove the created .gogetty file
	_ = os.Remove(".gogetty")
}

func TestReadAndWriteConfig(t *testing.T) {
	// Setup: Create a new temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gogetty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Ensure the temporary directory is removed after the test
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, ".gogetty")

	// Initial config for testing WriteConfig
	initialConfig := GogettyConfig{
		Dependencies: []Dependency{{Name: "test", URL: "http://example.com", Branch: "main", Commit: "abc123"}},
		LinkSubdir:   "modules",
	}

	// Write the config
	if err := WriteConfig(configPath, initialConfig); err != nil {
		t.Errorf("Failed to write config: %v", err)
	}

	// Read the config
	readConfig, err := ReadConfig(configPath)
	if err != nil {
		t.Errorf("Failed to read config: %v", err)
	}

	// Compare the written and read config
	if !reflect.DeepEqual(initialConfig, readConfig) {
		t.Errorf("Read config does not match written config")
	}
}

func TestAddAndRemoveDependency(t *testing.T) {
	// Setup: Create a new temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gogetty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Ensure the temporary directory is removed after the test
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, ".gogetty")
	initialConfig := GogettyConfig{
		Dependencies: []Dependency{},
		LinkSubdir:   "modules",
	}
	_ = WriteConfig(configPath, initialConfig)

	// Test AddDependency
	_ = AddDependency(configPath, "http://example.com", "main", "abc123", []string{})
	updatedConfig, _ := ReadConfig(configPath)
	if len(updatedConfig.Dependencies) != 1 || updatedConfig.Dependencies[0].Name != "test" {
		t.Errorf("Failed to add dependency")
	}

	// Test RemoveDependency
	_ = RemoveDependency(configPath, "test")
	updatedConfig, _ = ReadConfig(configPath)
	if len(updatedConfig.Dependencies) != 0 {
		t.Errorf("Failed to remove dependency")
	}
}
