package project

import (
	"os"
	"path/filepath"
)

func Validate(projectDir string) error {
	path := filepath.Join(projectDir, ProjectJson)
	_, err := os.Stat(path)
	if err != nil {
		return Init()
	}
	return nil
}
