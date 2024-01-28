package gitop

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type GitRepo struct {
	Path   string // Physical path of the repository
	URL    string // Repository URL from the .git/config
	Branch string // Repository Branch from HEAD
	Commit string // Repository Commit from .git/refs/heads/{Branch}
	Name   string // Name derived from URL
}

func Find(template GitRepo, repos []GitRepo) *GitRepo {
	for i, repo := range repos {
		if (template.URL == "" || repo.URL == template.URL) &&
			(template.Branch == "" || repo.Branch == template.Branch) &&
			(template.Commit == "" || repo.Commit == template.Commit) {
			return &repos[i] // Return a pointer to the actual slice element
		}
	}
	return nil // Return nil if no matching GitRepo is found
}

func getRepository(repoDir string) (GitRepo, error) {
	var repo GitRepo

	// Define common paths
	gitDir := filepath.Join(repoDir, ".git")
	configPath := filepath.Join(gitDir, "config")
	headPath := filepath.Join(gitDir, "HEAD")
	shallowPath := filepath.Join(gitDir, "shallow")

	// Check if .git directory exists to validate the Git repository
	if _, err := os.Stat(gitDir); err != nil {
		if os.IsNotExist(err) {
			return repo, fmt.Errorf("not a valid Git repository: %s", repoDir)
		}
		return repo, err
	}

	repo.Path = repoDir

	// Fetch Repository URL
	url, err := getRepoURL(configPath)
	if err != nil {
		return repo, fmt.Errorf("error fetching repository URL: %v", err)
	}
	repo.URL = strings.TrimSpace(url)
	repo.Name = strings.TrimSpace(GetNameFromURL(url))

	// Fetch Current Branch
	branch, err := getCurrentBranch(headPath)
	if err != nil {
		return repo, fmt.Errorf("error fetching current branch: %v", err)
	}
	repo.Branch = strings.TrimSpace(branch)

	// Fetch Latest Commit
	commit, err := os.ReadFile(shallowPath)
	if err != nil {
		debug.PrintStack()
		return repo, fmt.Errorf("error fetching latest commit: %v", err)
	}
	repo.Commit = strings.TrimSpace(string(commit))

	return repo, nil
}

func GetNameFromURL(gitURL string) string {
	urlParts := strings.Split(gitURL, "/")
	lastPart := urlParts[len(urlParts)-1]
	return strings.TrimSuffix(lastPart, ".git")
}

func getRepoURL(configPath string) (string, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "url =") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("repository URL not found in git config")
}

func getCurrentBranch(headPath string) (string, error) {
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	// Example content: "ref: refs/heads/main"
	if strings.HasPrefix(string(data), "ref: ") {
		parts := strings.Split(string(data), "/")
		return parts[len(parts)-1], nil
	}
	return "", fmt.Errorf("current branch not found")
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
