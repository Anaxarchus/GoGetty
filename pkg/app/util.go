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
	err := project.Validate(projectDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			return nil
		}
	}

	proj, err := project.GetProjectFile(projectDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			return nil
		}
	}

	var allErrors []error

	for _, dep := range proj.Dependencies {
		repo := gitop.Find(dep.Repository, modules)
		if repo == nil {
			repo, err = gitop.Fetch(cache.ModuleDir(), dep.Repository.URL, dep.Repository.Branch, dep.Repository.Commit)
			if err != nil {
				if !os.IsNotExist(err) {
					allErrors = append(allErrors, err)
				}
				continue
			}

			if err = fetchRecursive(repo.Path, modules); err != nil {
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
		return os.MkdirAll(dirName, 0755)
	}
	return nil
}

func printDependency(dep project.Dependency) {
	indent := "    "
	if dep.Repository.Name != "" {
		fmt.Println(indent+"Name:", dep.Repository.Name)
	}
	if dep.Repository.URL != "" {
		fmt.Println(indent+"Url:", dep.Repository.URL)
	}
	if dep.Repository.Branch != "" {
		fmt.Println(indent+"Branch:", dep.Repository.Branch)
	}
	if dep.Repository.Commit != "" {
		fmt.Println(indent+"Commit:", dep.Repository.Commit)
	}
	if len(dep.Directories) > 0 {
		fmt.Println(indent + "Directories:")
		for _, dir := range dep.Directories {
			fmt.Println(indent + indent + dir)
		}
	}
	fmt.Println()
}
