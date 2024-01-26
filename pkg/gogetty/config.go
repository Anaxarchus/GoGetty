package gogetty

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type GogettyConfig struct {
	Dependencies []Dependency `json:"dependencies"`
	LinkSubdir   string       `json:"linkSubdir"` // New field for symlink subdirectory
}

type Dependency struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Branch      string   `json:"branch"`
	Commit      string   `json:"commit"`
	Directories []string `json:"directories"`
}

const ConfigJson = ".gogetty"

func InitProject() {
	if _, err := os.Stat(ConfigJson); err == nil {
		fmt.Println(".gogetty already exists")
		return
	}

	config := GogettyConfig{
		Dependencies: []Dependency{},
		LinkSubdir:   "modules", // Set default value
	}
	file, err := os.Create(ConfigJson)
	if err != nil {
		fmt.Println("Error creating .gogetty:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		fmt.Println("Error writing to .gogetty:", err)
	}
}

func ReadConfig(filePath string) (GogettyConfig, error) {
	var config GogettyConfig

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func WriteConfig(filePath string, config GogettyConfig) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func AddDependency(filePath, url, branch, commit string, directories []string) error {
	config, err := ReadConfig(filePath)
	if err != nil {
		return err
	}

	newDependency := Dependency{
		URL:         url,
		Branch:      branch,
		Commit:      commit,
		Directories: directories, // Set the directories
	}

	// Check if dependency already exists and update it
	updated := false
	for i, dep := range config.Dependencies {
		if dep.URL == url {
			config.Dependencies[i] = newDependency
			updated = true
			break
		}
	}

	if !updated {
		config.Dependencies = append(config.Dependencies, newDependency)
	}

	return WriteConfig(filePath, config)
}

func RemoveDependency(filePath, name string) error {
	config, err := ReadConfig(filePath)
	if err != nil {
		return err
	}

	for i, d := range config.Dependencies {
		if d.Name == name {
			config.Dependencies = append(config.Dependencies[:i], config.Dependencies[i+1:]...)
			return WriteConfig(filePath, config)
		}
	}

	return fmt.Errorf("dependency not found: %s", name)
}

func addDirToGitIgnore(gitIgnorePath, dir string) error {
	// Check if .gitignore file exists
	file, err := os.OpenFile(gitIgnorePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read existing content to avoid duplicating the entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == dir {
			return nil // Directory already in .gitignore
		}
	}

	// Add directory to .gitignore
	_, err = file.WriteString("\n" + dir + "\n")
	return err
}
