package symlink

import (
	"fmt"
	"os"
	"path/filepath"
)

const symlinkWarning = `Welcome to your GoGetty-managed directory!

Please note: This directory includes symlinks managed by GoGetty and may be modified automatically. It's recommended to avoid saving your personal project files directly in this directory to prevent any accidental data loss. Thanks for using GoGetty!`

func CreateSymlink(source, target string) error {
	// Resolve absolute paths
	absSource, err := filepath.Abs(source)
	if err != nil {
		return err
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	// Remove existing symlink/file if it exists
	if _, err := os.Lstat(absTarget); err == nil {
		if err := os.Remove(absTarget); err != nil {
			return err
		}
	}

	return os.Symlink(absSource, absTarget)
}

func CreateSymlinkBundle(sourceBase, targetBase string, subDirs []string) error {
	// Determine the name of the new directory to be created at the target location
	newDirName := filepath.Base(sourceBase)
	newDirPath := filepath.Join(targetBase, newDirName)

	// Create the new directory
	if err := os.MkdirAll(newDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %v", newDirPath, err)
	}

	// Create symlinks for each subdirectory within the new directory
	for _, subDir := range subDirs {
		sourcePath := filepath.Join(sourceBase, subDir)
		targetPath := filepath.Join(newDirPath, subDir) // target path within the new directory

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
