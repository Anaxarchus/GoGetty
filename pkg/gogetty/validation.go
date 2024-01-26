package gogetty

import (
	"fmt"
	"os"
	"os/exec"
)

func ValidateEnvironment() error {

	// Check if git is installed
	if err := checkGitInstalled(); err != nil {
		fmt.Println("Error: Git is required but not found:", err)
		fmt.Println("Please install Git to use this application. Visit https://git-scm.com for installation instructions.")
		os.Exit(1)
	}

	// Check if the cache is initialized
	if err := checkCacheInitialized(); err != nil {
		err = InitCache()
		if err != nil {
			fmt.Println("Error: Cache failed to initialize:", err)
		}
	}

	return nil
}

func checkGitInstalled() error {
	_, err := exec.LookPath("git")
	return err
}

func checkCacheInitialized() error {
	cachePath := GetCachePath(CacheJson)
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return fmt.Errorf("cache is not initialized")
	}
	return nil
}
