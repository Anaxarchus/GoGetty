package symlink

import (
	"fmt"
	"os"
	"path/filepath"
)

const symlinkWarning = `Welcome to your GoGetty-managed directory!

Please note: This directory includes symlinks managed by GoGetty and may 
be modified automatically. It's recommended to avoid saving your personal 
project files directly in this directory to prevent any accidental data 
loss.

Thanks for using GoGetty!`

func CreateSymlink(source, target string) error {
	// Create the parent directory of the target path
	parentDir := filepath.Dir(target)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory '%s': %v", parentDir, err)
	}

	// Remove existing symlink/file if it exists
	if _, err := os.Lstat(target); err == nil {
		if err := os.Remove(target); err != nil {
			return err
		}
	}

	return os.Symlink(source, target)
}

func CreateSymlinkBundle(sourceBase, targetBase string, subDirs []string) error {

	// Create symlinks for each subdirectory within the new directory
	for _, subDir := range subDirs {
		sourcePath := filepath.Join(sourceBase, subDir)
		targetPath := filepath.Join(targetBase, subDir)

		// Create the symlink
		if err := CreateSymlink(sourcePath, targetPath); err != nil {
			return fmt.Errorf("failed to create symlink from '%s' to '%s': %v", sourcePath, targetPath, err)
		}
	}
	return nil
}

func WriteReadmeWithWarning(targetBase string) error {
	readmePath := filepath.Join(targetBase, "WARNING.md")

	// Check if README already exists
	if _, err := os.Stat(readmePath); err == nil {
		// README exists, no need to rewrite
		return nil
	}

	// Write the warning message to README
	return os.WriteFile(readmePath, []byte(symlinkWarning), 0644)
}
