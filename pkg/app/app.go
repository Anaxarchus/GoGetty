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

type App interface {
	Init() error
	Add(url, branch, commit string, directories []string) error
	Remove(name string) error
	Fetch() error
	Update(name, branch, commit string, directories []string) error
	List() ([]project.Dependency, error)
	Clean() error
	ListModuleNames() ([]string, error)
}

type MyApp struct {
	ProjectDir string
	Cache      []gitop.GitRepo
}

func (m *MyApp) Init() error {

	err := ValidateEnvironment()
	if err != nil {
		fmt.Printf("Error in ValidateEnvironment: %v\n", err)
		return err
	}

	err = project.Init()
	if err != nil {
		fmt.Printf("Error in project.Init: %v\n", err)
		return err
	}

	cache.AddClient(m.ProjectDir)
	return nil
}

func (m *MyApp) Add(url, branch, commit string, directories []string) error {
	err := ValidateEnvironment()
	if err != nil {
		return err
	}
	err = project.Validate("")
	if err != nil {
		return err
	}

	repo := gitop.GitRepo{
		URL:    url,
		Branch: branch,
		Commit: commit,
	}

	return project.AddDependency(repo, directories)
}

func (m *MyApp) Update(name, branch, commit string, directories []string) error {
	err := ValidateEnvironment()
	if err != nil {
		return err
	}
	err = project.Validate("")
	if err != nil {
		return err
	}

	repo := gitop.GitRepo{
		Name:   name,
		Branch: branch,
		Commit: commit,
	}

	new_dep := project.Dependency{
		Repository:  repo,
		Directories: directories,
	}

	dep, err := project.Find(name)
	if err != nil {
		return err
	}

	return project.UpdateDependency(dep, new_dep)

}

func (m *MyApp) Remove(name string) error {
	err := ValidateEnvironment()
	if err != nil {
		return err
	}
	err = project.Validate("")
	if err != nil {
		return err
	}

	return project.RemoveDependency(name)
}

func (m *MyApp) Fetch() error {
	// Validate the environment
	if err := ValidateEnvironment(); err != nil {
		return err
	}

	// Validate the project
	if err := project.Validate(""); err != nil {
		return err
	}

	// Get the project file
	proj, projErr := project.GetProjectFile(m.ProjectDir)
	if projErr != nil {
		return projErr
	}

	// Determine the target directory
	targetDir := filepath.Join(m.ProjectDir, proj.ModulesDir)

	// Check if targetDir exists and delete it if it does
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing target directory: %w", err)
		}
	}

	if len(proj.Dependencies) == 0 {
		return nil
	}

	// Perform recursive fetch
	if err := fetchRecursive(m.ProjectDir, m.Cache); err != nil {
		return err
	}

	// Write a warning file after successful fetching
	if err := symlink.WriteReadmeWithWarning(targetDir); err != nil {
		return fmt.Errorf("failed to write warning file: %w", err)
	}

	return nil
}

func (m *MyApp) List() error {
	err := ValidateEnvironment()
	if err != nil {
		return err
	}

	err = project.Validate("")
	if err != nil {
		return err
	}

	proj, err := project.GetProjectFile("")
	if err != nil {
		return err
	}

	// Check dependencies and print them
	if len(proj.Dependencies) == 0 {
		fmt.Println("No dependencies found.")
		return nil
	}

	fmt.Println("Dependencies:")
	for _, dep := range proj.Dependencies {
		printDependency(dep)
	}
	return nil
}

func (m *MyApp) Clean() error {
	err := ValidateEnvironment()
	if err != nil {
		return err
	}
	clients, err := cache.GetClients()
	if err != nil {
		return err
	}

	// Store project directories that should be deleted
	directoriesToDelete := []string{}

	// Store dependencies of projects with a .gogetty file
	dependencies := map[string]project.Dependency{}

	// Iterate over each project directory
	for _, client := range clients {
		// Check if the .gogetty file exists in the project directory
		validClient := project.Validate(client) != nil
		if validClient {
			// .gogetty file doesn't exist, mark for deletion
			directoriesToDelete = append(directoriesToDelete, client)
		} else if err == nil {
			// .gogetty file exists, read its dependencies
			project, projErr := project.GetProjectFile(client)
			if projErr == nil {
				for _, dep := range project.Dependencies {
					dependencies[dep.Repository.URL] = dep
				}
			} else {
				// Handle error reading .gogetty file if needed
				fmt.Printf("Error reading .gogetty file in %s: %v\n", client, projErr)
			}
		}
	}

	// Iterate over all modules and check if they have dependents
	for _, module := range m.Cache {
		if _, exists := dependencies[module.URL]; !exists {
			// Module has no dependents, remove it from the cache
			if err := cache.Remove(module.Path); err != nil {
				// Handle removal error if needed
				fmt.Printf("Error removing module from cache: %v\n", err)
			}
		}
	}

	// Delete project directories with no .gogetty file
	for _, dir := range directoriesToDelete {
		if err := os.RemoveAll(dir); err != nil {
			// Handle directory removal error if needed
			fmt.Printf("Error removing directory %s: %v\n", dir, err)
		}
	}

	return nil
}

func (m *MyApp) ListModuleNames() ([]string, error) {
	var moduleNames []string
	for _, repo := range m.Cache {
		if repo.Name != "" {
			moduleNames = append(moduleNames, repo.Name)
		}
	}
	return moduleNames, nil
}
