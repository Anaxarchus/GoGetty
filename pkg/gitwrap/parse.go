package gitwrap

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type GitConfig struct {
}

func ParseRepository(repoDir string) (GitRepo, error) {
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
	url, err := parseURL(configPath)
	if err != nil {
		return repo, fmt.Errorf("error fetching repository URL: %v", err)
	}
	repo.URL = strings.TrimSpace(url)
	repo.Name = strings.TrimSpace(parseName(url))

	// Fetch Current Branch
	branch, err := parseBranch(headPath)
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

func parseURL(configPath string) (string, error) {
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

func parseName(gitURL string) string {
	urlParts := strings.Split(gitURL, "/")
	lastPart := urlParts[len(urlParts)-1]
	return strings.TrimSuffix(lastPart, ".git")
}

func parseBranch(headPath string) (string, error) {
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(string(data), "ref: ") {
		parts := strings.Split(string(data), "/")
		return parts[len(parts)-1], nil
	}
	return "", fmt.Errorf("current branch not found")
}
