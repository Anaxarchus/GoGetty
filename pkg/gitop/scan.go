package gitop

import (
	"fmt"
	"os"
	"path/filepath"
)

func Scan(directory string) ([]GitRepo, error) {
	var repos []GitRepo

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error encountered while walking through directory: %v\n", err)
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path) // Get the path of the repository

			repo, err := getRepository(repoPath)
			if err != nil {
				fmt.Printf("Error processing repository at %s: %v\n", repoPath, err)
				return filepath.SkipDir // Skip this directory but continue walking
			}

			repos = append(repos, repo)
			return filepath.SkipDir // Skip the .git directory to avoid redundant processing
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error during directory scan: %v\n", err)
		return nil, err
	}

	return repos, nil
}
