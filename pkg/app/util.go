package app

import (
	"fmt"
	"gogetty/pkg/cache"
	"gogetty/pkg/gitwrap"
	"gogetty/pkg/godot"
	"gogetty/pkg/input"
	"gogetty/pkg/project"
	"gogetty/pkg/symlink"
	"os"
	"path/filepath"
)

func fetchRecursive(projectDir string, modules []gitwrap.GitRepo) error {
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
		repo := gitwrap.FindInList(dep.Repository, modules)
		if repo == nil {
			repo, err = gitwrap.Fetch(cache.ModuleDir(), dep.Repository.URL, dep.Repository.Branch, dep.Repository.Commit)
			if err != nil {
				if !os.IsNotExist(err) {
					allErrors = append(allErrors, err)
				}
				continue
			}

			godotProject, err := godot.GetGodotProject(repo.Path)
			if err == nil {
				fmt.Println("A Godot project file was found.")
				fmt.Println("Note: assume-unchanged will update the repository's index to assume `godot.project` unchanged")
				fmt.Println("meaning that until restored, git will stop tracking changes to that file. This makes it safe")
				fmt.Println("to commit your changes without losing your project file, but you'll need to remember to reverse")
				fmt.Println("the operation in the future: `gogetty <moduleName> update-index --no-assume-unchanged godot.project`")

				userChoice := input.Option("What do you want to do?", "nothing", "delete", "assume-unchanged and delete")
				switch userChoice {
				case 1: // delete
					err = godot.RemoveProjectFile(godotProject.Path)
					if err != nil {
						fmt.Printf("Error while removing Godot project file: %v\n", err)
					} else {
						fmt.Println("Godot project file removed successfully.")
					}
				case 2: // assume-unchanged and delete
					err = gitwrap.AssumeUnchanged(godotProject.Path, "godot.project")
					if err != nil {
						fmt.Printf("Error while setting Godot project file as assume-unchanged: %v\n", err)
						break
					}
					err = godot.RemoveProjectFile(godotProject.Path)
					if err != nil {
						fmt.Printf("Error while removing Godot project file: %v\n", err)
					} else {
						fmt.Println("Godot project file removed successfully.")
					}
				default: // nothing or invalid choice
					fmt.Println("No action taken.")
				}

				if err = godot.UpdateProjectPaths(*godotProject); err != nil {
					fmt.Printf("Error while updating godot project paths: %v\n", err)
				}
			} else {
				fmt.Printf("Error while getting Godot project: %v\n", err)
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
