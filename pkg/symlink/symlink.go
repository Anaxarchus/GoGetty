package symlink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const symlinkWarning = `Welcome to your GoGetty-managed directory!

Please note: This directory includes symlinks managed by GoGetty and may 
be modified automatically. It's recommended to avoid saving your personal 
project files directly in this directory to prevent any accidental data 
loss.

Thanks for using GoGetty!`

func CreateSymlink(source, target string) error {
	// Resolve absolute paths
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for source '%s': %v", source, err)
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for target '%s': %v", target, err)
	}

	// Check for recursive symlink
	if absSource == absTarget {
		return fmt.Errorf("cannot create symlink: source and target are the same")
	}
	if strings.HasPrefix(absTarget, absSource+string(os.PathSeparator)) {
		return fmt.Errorf("cannot create recursive symlink")
	}

	// Create the parent directory of the target path
	parentDir := filepath.Dir(absTarget)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory '%s': %v", parentDir, err)
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
	readmePath := filepath.Join(targetBase, "WARN.md")

	// Check if README already exists
	if _, err := os.Stat(readmePath); err == nil {
		// README exists, no need to rewrite
		return nil
	}

	// Write the warning message to README
	return os.WriteFile(readmePath, []byte(symlinkWarning), 0644)
}
