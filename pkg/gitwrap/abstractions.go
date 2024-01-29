package gitwrap

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitRepo struct {
	Path   string // Physical path of the repository
	URL    string // Repository URL from the .git/config
	Branch string // Repository Branch from HEAD
	Commit string // Repository Commit from .git/refs/heads/{Branch}
	Name   string // Name derived from URL
}

// AssumeUnchanged marks a file as "assume-unchanged" in a given directory
func AssumeUnchanged(dirPath, fileName string) error {
	// Construct the full path to the file
	fullPath := filepath.Join(dirPath, fileName)

	// Construct and execute the git update-index --assume-unchanged command
	cmd := exec.Command("git", "update-index", "--assume-unchanged", fullPath)

	// Set the working directory
	cmd.Dir = dirPath

	// Execute the command
	_, err := cmd.CombinedOutput()
	return err
}

// Forward executes a given git command with parameters in the specified path
func Forward(path, gitCommand string, params ...string) error {
	// Create the full command slice
	commandArgs := append([]string{gitCommand}, params...)
	cmd := exec.Command("git", commandArgs...)

	// Set the working directory if a path is provided
	if path != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		cmd.Dir = absPath
	}

	// Set the output to be written to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	return err
}

// Ignore appends multiple ignoreStrings to the .gitignore file in the specified GitRepo.
func Ignore(repoDir string, ignoreStrings ...string) error {
	// Construct the path to the .gitignore file in the repository
	gitignorePath := filepath.Join(repoDir, ".gitignore")

	// Open existing .gitignore file or create a new one
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening .gitignore file: %w", err)
	}
	defer file.Close()

	// Read existing .gitignore content to check for duplicates
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading .gitignore file: %w", err)
	}

	for _, ignoreString := range ignoreStrings {
		// Check if ignoreString already exists in file
		if strings.Contains(string(content), ignoreString) {
			fmt.Printf("Ignore string '%s' already exists in .gitignore\n", ignoreString)
			continue
		}

		// Append the ignoreString
		if _, err = file.WriteString("\n" + ignoreString); err != nil {
			return fmt.Errorf("error writing to .gitignore file: %w", err)
		}
	}

	return nil
}

// RemoveIgnore removes specified ignoreStrings from the .gitignore file in the specified GitRepo.
func RemoveIgnore(repoDir string, ignoreStrings ...string) error {
	// Construct the path to the .gitignore file in the repository
	gitignorePath := filepath.Join(repoDir, ".gitignore")

	// Open the .gitignore file for reading
	file, err := os.OpenFile(gitignorePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening .gitignore file: %w", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var updatedLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line matches any of the ignoreStrings
		shouldRemove := false
		for _, ignoreString := range ignoreStrings {
			if line == ignoreString {
				shouldRemove = true
				break
			}
		}

		// If the line should not be removed, keep it in the updatedLines slice
		if !shouldRemove {
			updatedLines = append(updatedLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .gitignore file: %w", err)
	}

	// Close the file and reopen it for writing, truncating its content
	file.Close()

	// Reopen the .gitignore file for writing, truncating its content
	file, err = os.Create(gitignorePath)
	if err != nil {
		return fmt.Errorf("error creating .gitignore file: %w", err)
	}
	defer file.Close()

	// Write the updated lines back to the .gitignore file
	for _, line := range updatedLines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to .gitignore file: %w", err)
		}
	}

	return nil
}
