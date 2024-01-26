package gogetty

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// createUnixLink handles symlink creation on Unix-like systems.
func createUnixLink(source, target string) error {
	return os.Symlink(source, target)
}

// createWindowsLink handles symlink creation on Windows.
func createWindowsLink(source, target string) error {
	cmd := exec.Command("cmd", "/c", "mklink", "/J", target, source)
	return cmd.Run()
}

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

	if runtime.GOOS == "windows" {
		return createWindowsLink(absSource, absTarget)
	} else {
		return createUnixLink(absSource, absTarget)
	}
}

func CreateBatchSymlinks(sourceBase, targetBase string, subDirs []string) error {
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

func RemoveBatchSymlinks(targetBase string) error {
	// Determine the name of the directory to be removed
	dirToRemove := filepath.Base(targetBase)

	// Construct the full path of the directory
	fullPath := filepath.Join(filepath.Dir(targetBase), dirToRemove)

	// Remove the entire directory
	if err := os.RemoveAll(fullPath); err != nil {
		return fmt.Errorf("failed to remove directory '%s': %v", fullPath, err)
	}

	return nil
}
