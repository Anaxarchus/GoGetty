package gogetty

import (
	"os/exec"
	"strings"
)

func FetchGitRepo(gitURL, branch, commit string) error {

	modulesDir := GetCachePath(ModuleDir)

	// Use ConstructGitCommand to build the command
	gitCommand := ConstructGitCommand(gitURL, branch, commit)

	// Split the command for execution
	cmdParts := strings.Fields(gitCommand)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = modulesDir
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ConstructGitCommand(gitURL, branch, commit string) string {
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

	return "git " + strings.Join(args, " ")
}
