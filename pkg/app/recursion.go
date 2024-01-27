package app

import (
	"fmt"
	"gogetty/pkg/cache"
	"gogetty/pkg/gitop"
	"gogetty/pkg/project"
	"gogetty/pkg/symlink"
	"os"
	"path/filepath"
)

func fetchRecursive(projectDir string, modules []gitop.GitRepo) error {
	err := project.Validate("")
	if err != nil {
		return err
	}

	proj, err := project.GetProjectFile("")
	if err != nil {
		return err
	}

	var allErrors []error

	for _, dep := range proj.Dependencies {
		repo := gitop.Find(dep.Repository, modules)
		if repo == nil {
			repo, fetchErr := gitop.Fetch(cache.ModuleDir(), dep.Repository.URL, dep.Repository.Branch, dep.Repository.Commit)
			if fetchErr != nil {
				allErrors = append(allErrors, fetchErr)
				continue
			}

			if err := fetchRecursive(repo.Path, modules); err != nil {
				allErrors = append(allErrors, err)
			}
		}

		targetDir := filepath.Join(projectDir, proj.ModulesDir, repo.Name)
		if err := ensureDir(targetDir); err != nil {
			allErrors = append(allErrors, fmt.Errorf("error creating directory %s: %v", targetDir, err))
			continue
		}
		if len(dep.Directories) > 0 {
			symlink.CreateSymlinkBundle(repo.Path, targetDir, dep.Directories)
		} else {
			symlink.CreateSymlink(repo.Path, targetDir)
		}
		new_dep := project.Dependency{
			Repository:  *repo,
			Directories: dep.Directories,
		}
		project.UpdateDependency(dep, new_dep)
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", allErrors)
	}

	return nil
}

// ensureDir checks if a directory exists, and if not, creates it
func ensureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0755) // or an appropriate permission as required
	}
	return nil
}