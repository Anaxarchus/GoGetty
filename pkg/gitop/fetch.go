package gitop

import (
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func Fetch(cacheDir, gitURL, branch, commit string) (*GitRepo, error) {
	// Derive the name from the gitURL
	parsedURL, err := url.Parse(gitURL)
	if err != nil {
		return &GitRepo{}, err
	}

	// Extract the repository name from the URL
	name := path.Base(parsedURL.Path)

	// Strip file extensions, if any
	name = strings.TrimSuffix(name, ".git")
	name = strings.TrimSuffix(name, ".bundle") // Add more extensions if needed

	if name == "." || name == "/" {
		return &GitRepo{}, fmt.Errorf("invalid gitURL: %s", gitURL)
	}

	// Append the name to the cacheDir
	fullCacheDir := filepath.Join(cacheDir, name)

	// Use ConstructGitCommand to build the command
	gitCommand := constructFetchCommand(gitURL, branch, commit, fullCacheDir)

	// Split the command for execution
	cmdParts := strings.Fields(gitCommand)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = cacheDir

	// Run the command and handle errors
	if err := cmd.Run(); err != nil {
		return &GitRepo{}, err
	}

	// Create and populate a GitRepo object
	repo := GitRepo{
		Path:   fullCacheDir,
		URL:    gitURL,
		Branch: branch,
		Commit: commit,
		Name:   name,
	}

	fmt.Printf("repo fetched: %v\n", repo)

	return &repo, nil
}

func constructFetchCommand(gitURL, branch, commit, cacheDir string) string {
	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, gitURL)

	// Construct checkout part if commit is specified
	if commit != "" {
		checkoutPart := []string{"&&", "git", "checkout", commit}
		args = append(args, checkoutPart...)
	}

	// Append the cacheDir as the target directory
	args = append(args, cacheDir)

	return "git " + strings.Join(args, " ")
}
