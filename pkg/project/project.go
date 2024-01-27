package project

import (
	"encoding/json"
	"fmt"
	"gogetty/pkg/gitop"
	"os"
	"path/filepath"
)

type Project struct {
	Dependencies []Dependency `json:"modules"`
	ModulesDir   string       `json:"modulesDirectory"`
}

type Dependency struct {
	Repository  gitop.GitRepo `json:"repository"`
	Directories []string      `json:"directories"`
}

const ProjectJson = ".gogetty"

func Init() error {
	if _, err := os.Stat(ProjectJson); err == nil {
		return fmt.Errorf(".gogetty already exists")
	}

	config := Project{
		Dependencies: []Dependency{},
		ModulesDir:   "modules",
	}
	file, err := os.Create(ProjectJson)
	if err != nil {
		return fmt.Errorf("error creating .gogetty: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error writing to .gogetty: %v", err)
	}

	return nil // No error occurred
}

func GetProjectFile(projectDir string) (Project, error) {
	return readProject(projectDir)
}

func UpdateDependency(old_dependency Dependency, new_dependency Dependency) error {
	project, err := readProject("")
	if err != nil {
		return err
	}

	updated := false
	for i, dep := range project.Dependencies {
		if dep.Repository.URL == old_dependency.Repository.URL {
			// Update the existing dependency with the new information
			project.Dependencies[i] = new_dependency
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("Dependency not found for URL: %s", old_dependency.Repository.URL)
	}

	return writeProject("", project)
}

func AddDependency(repo gitop.GitRepo, directories []string) error {
	project, err := readProject("")
	if err != nil {
		return err
	}

	newDependency := Dependency{
		Repository:  repo,
		Directories: directories,
	}
	fmt.Println("repo: ", repo)

	updated := false
	for i, dep := range project.Dependencies {
		if dep.Repository.URL == newDependency.Repository.URL {
			project.Dependencies[i] = newDependency
			updated = true
			break
		}
	}

	if !updated {
		project.Dependencies = append(project.Dependencies, newDependency)
	}

	return writeProject("", project)
}

func RemoveDependency(repoName string) error {
	project, err := readProject("")
	if err != nil {
		return err
	}

	for i, d := range project.Dependencies {
		if d.Repository.Name == repoName {
			project.Dependencies = append(project.Dependencies[:i], project.Dependencies[i+1:]...)
			return writeProject("", project)
		}
	}

	return fmt.Errorf("dependency not found: %s", repoName)
}

func Find(name string) (Dependency, error) {
	project, err := readProject("")
	if err != nil {
		return Dependency{}, err
	}

	for _, dep := range project.Dependencies {
		if dep.Repository.Name == name {
			return dep, nil
		}
	}

	return Dependency{}, fmt.Errorf("dependency not found for module: %s", name)
}

func readProject(projectDir string) (Project, error) {
	var project Project
	path := filepath.Join(projectDir, ProjectJson)
	file, err := os.Open(path)
	if err != nil {
		return project, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&project); err != nil {
		return project, err
	}

	return project, nil
}

func writeProject(projectDir string, project Project) error {
	path := filepath.Join(projectDir, ProjectJson)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(project)
}
