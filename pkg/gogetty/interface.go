package gogetty

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type GogettyAction interface {
	Init() error
	Add(url, branch, commit string) error
	Remove(name string) error
	Fetch() error
	Update(name, branch, commit string) error
	List() ([]Dependency, error)
	Clean() error
}

type GogettyManager struct {
	// Add any necessary fields here. For example:
	ProjectDir string
}

func (gm *GogettyManager) Init() error {
	// Check if .gogetty file already exists in the project directory
	projectConfigPath := filepath.Join(gm.ProjectDir, ConfigJson)
	if _, err := os.Stat(projectConfigPath); err == nil {
		return fmt.Errorf(".gogetty configuration file already exists in the project directory")
	}

	// Create the .gogetty file with default configuration
	config := GogettyConfig{
		Dependencies: []Dependency{},
		LinkSubdir:   "modules",
	}

	if err := WriteConfig(projectConfigPath, config); err != nil {
		return fmt.Errorf("failed to write .gogetty configuration: %v", err)
	}

	// Add linkSubdir to .gitignore
	gitIgnorePath := filepath.Join(gm.ProjectDir, ".gitignore")
	if err := addDirToGitIgnore(gitIgnorePath, config.LinkSubdir); err != nil {
		return fmt.Errorf("failed to update .gitignore: %v", err)
	}

	fmt.Println(".gogetty initialized successfully in the project directory")
	return nil
}

func (gm *GogettyManager) Add(url, branch, commit string, directories []string) error {
	// Check if the project's .gogetty configuration file exists
	configFilePath := filepath.Join(gm.ProjectDir, ConfigJson)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found. Please run 'gogetty init' to initialize the project")
	}

	// Add the dependency to the project's .gogetty configuration
	err := AddDependency(configFilePath, url, branch, commit, directories)
	if err != nil {
		return fmt.Errorf("failed to add dependency: %v. Please check the details and try again", err)
	}

	return nil
}

func (gm *GogettyManager) Remove(name string) error {
	absProjectDir, err := filepath.Abs(gm.ProjectDir)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path for project directory '%s': %v", gm.ProjectDir, err)
	}
	gm.ProjectDir = absProjectDir

	configFilePath := filepath.Join(gm.ProjectDir, ConfigJson)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found. Please run 'gogetty init' to initialize the project")
	}

	// Load the .gogetty configuration
	config, err := ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read .gogetty file: %v", err)
	}

	// Remove the dependency from the project's .gogetty configuration
	err = RemoveDependency(configFilePath, name)
	if err != nil {
		return fmt.Errorf("failed to remove dependency: %v", err)
	}

	// Unregister current project as a dependent
	err = RemoveDependentFromModule(name, gm.ProjectDir)
	if err != nil {
		return fmt.Errorf("failed to unregister project as a dependent of '%s': %v", name, err)
	}

	// Remove the symlink for the dependency
	symlinkPath := filepath.Join(gm.ProjectDir, config.LinkSubdir, name)
	if err := os.Remove(symlinkPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove symlink for '%s': %v", name, err)
	}

	// Check if the modules directory is empty and remove it if so
	modulesDir := filepath.Join(gm.ProjectDir, config.LinkSubdir)
	isEmpty, err := isDirEmpty(modulesDir)
	if err != nil {
		return fmt.Errorf("failed to check if modules directory is empty: %v", err)
	}
	if isEmpty {
		if err := os.Remove(modulesDir); err != nil {
			return fmt.Errorf("failed to remove empty modules directory: %v", err)
		}
	}

	return nil
}

func (gm *GogettyManager) Fetch() error {
	absProjectDir, err := filepath.Abs(gm.ProjectDir)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path for project directory '%s': %v", gm.ProjectDir, err)
	}
	gm.ProjectDir = absProjectDir

	configFilePath := filepath.Join(gm.ProjectDir, ConfigJson)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf(".gogetty configuration file not found. Please run 'gogetty init' to initialize the project")
	}

	config, err := ReadConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read .gogetty file: %v. Please ensure the file is correctly formatted", err)
	}

	// Validate or create the modules directory
	modulesDir := filepath.Join(gm.ProjectDir, config.LinkSubdir)
	if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(modulesDir, 0755); err != nil {
			return fmt.Errorf("failed to create modules directory '%s': %v", modulesDir, err)
		}
	}

	// Flag to indicate if the config needs updating
	needsConfigUpdate := false

	for i, dependency := range config.Dependencies {
		moduleName := getModuleNameFromURL(dependency.URL)

		// Update the name of the module in the dependency if necessary
		if dependency.Name != moduleName {
			config.Dependencies[i].Name = moduleName
			needsConfigUpdate = true
		}

		dependency.Name = moduleName
		WriteConfig(configFilePath, config)

		// Check if the module is already in cache
		moduleInCache, err := ModuleInCache(dependency.URL, dependency.Branch, dependency.Commit)
		if err != nil {
			return fmt.Errorf("error checking cache for module '%s': %v", moduleName, err)
		}

		if !moduleInCache {
			// Fetch the module if it's not in the cache
			err = FetchGitRepo(dependency.URL, dependency.Branch, dependency.Commit)
			if err != nil {
				return fmt.Errorf("failed to fetch dependency '%s': %v. Please check the repository URL and network connectivity", moduleName, err)
			}

			// Add the fetched module to the cache
			err = AddModuleToCache(moduleName, dependency.URL, dependency.Branch, dependency.Commit)
			if err != nil {
				return fmt.Errorf("failed to add module '%s' to cache: %v", moduleName, err)
			}
		}

		// Check if current project is already a dependent
		isDependent, err := ModuleHasDependency(moduleName, gm.ProjectDir)
		if err != nil {
			return fmt.Errorf("error checking dependent status for module '%s': %v", moduleName, err)
		}

		if !isDependent {
			// Register the current project as a dependent of the module
			err = AddDependentToModule(moduleName, gm.ProjectDir)
			if err != nil {
				return fmt.Errorf("failed to register project as a dependent of '%s': %v", moduleName, err)
			}
		}

		// Determine source path for symlink
		sourcePath := GetCachePath(filepath.Join(ModuleDir, moduleName))

		// Check if dependency specifies directories
		if len(dependency.Directories) > 0 {
			targetPath := filepath.Join(gm.ProjectDir, config.LinkSubdir)
			if err := CreateBatchSymlinks(sourcePath, targetPath, dependency.Directories); err != nil {
				return fmt.Errorf("failed to create symlink for dependency '%s': %v", moduleName, err)
			}
		} else {
			// Create a single symlink for the whole module
			targetPath := filepath.Join(gm.ProjectDir, config.LinkSubdir, moduleName)
			if err := CreateSymlink(sourcePath, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink for dependency '%s': %v", moduleName, err)
			}
		}

		fmt.Printf("Dependency '%s' processed successfully.\n", moduleName)
	}

	// Update the config file if any changes were made
	if needsConfigUpdate {
		if err := WriteConfig(configFilePath, config); err != nil {
			return fmt.Errorf("failed to update .gogetty configuration: %v", err)
		}
	}

	return nil
}

func (gm *GogettyManager) Update(name, newBranch, newCommit string, newDirectories []string) error {
	configFilePath := filepath.Join(gm.ProjectDir, ConfigJson)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf(".gogetty configuration file not found. Please run 'gogetty init' to initialize the project")
	}

	// Load existing cache
	cache, err := loadCache(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read cache file: %v", err)
	}

	moduleInfo, exists := cache.Modules[name]
	if !exists {
		return fmt.Errorf("module '%s' not found in cache", name)
	}

	// Update or add the module
	newURL := moduleInfo.URL // Assuming the URL remains the same
	if err := gm.Add(newURL, newBranch, newCommit, newDirectories); err != nil {
		return fmt.Errorf("failed to add updated module: %v", err)
	}

	// Clean up the old module's batch directory if it existed
	oldBatchDir := filepath.Join(gm.ProjectDir, name)
	if err := RemoveBatchSymlinks(oldBatchDir); err != nil {
		return fmt.Errorf("failed to remove old batch directory: %v", err)
	}

	// Handle new batch creation or single module update
	newPath := GetCachePath(filepath.Join(ModuleDir, filepath.Base(newURL)))
	if len(newDirectories) > 0 {
		// Create symlinks for the new batch of directories
		if err := CreateBatchSymlinks(newPath, gm.ProjectDir, newDirectories); err != nil {
			return fmt.Errorf("failed to create new batch symlinks: %v", err)
		}
	} else {
		// Update symlink to point to the new module
		targetPath := filepath.Join(gm.ProjectDir, name)
		if err := CreateSymlink(newPath, targetPath); err != nil {
			return fmt.Errorf("failed to create new symlink: %v", err)
		}
	}

	return nil
}

func (gm *GogettyManager) List() ([]Dependency, error) {
	projectConfigPath := filepath.Join(gm.ProjectDir, ConfigJson)

	// Check if .gogetty file exists
	if _, err := os.Stat(projectConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".gogetty configuration file not found in the project directory")
	}

	// Read the .gogetty configuration file
	config, err := ReadConfig(projectConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .gogetty file: %v", err)
	}

	// Return the list of dependencies
	return config.Dependencies, nil
}

func getModuleNameFromURL(url string) string {
	baseName := filepath.Base(url)
	return strings.TrimSuffix(baseName, ".git") // Remove .git if present
}

func isDirEmpty(dirname string) (bool, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Try to read at least one entry
	if err == io.EOF {
		return true, nil // EOF means the directory is empty
	}
	return false, err // Any other error (including non-EOF) means it's not empty or an error occurred
}

func (gm *GogettyManager) Clean() error {
	cacheFilePath := GetCachePath(CacheJson)

	// Load existing cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return fmt.Errorf("failed to read cache file: %v", err)
	}

	// Track if the cache is updated
	cacheUpdated := false

	// Iterate over all modules in the cache
	for moduleName, moduleInfo := range cache.Modules {
		// Check if the module still exists
		if _, err := os.Stat(moduleInfo.Path); os.IsNotExist(err) {
			delete(cache.Modules, moduleName)
			cacheUpdated = true
			continue
		}

		// Iterate over dependents and verify their existence
		for i := len(moduleInfo.Dependents) - 1; i >= 0; i-- {
			dependentProjectPath := filepath.Join(moduleInfo.Dependents[i], ConfigJson)
			if _, err := os.Stat(dependentProjectPath); os.IsNotExist(err) {
				moduleInfo.Dependents = append(moduleInfo.Dependents[:i], moduleInfo.Dependents[i+1:]...)
				cacheUpdated = true
			}
		}

		// If no dependents left, remove the module
		if len(moduleInfo.Dependents) == 0 {
			if err := DeleteModuleFromCache(moduleName); err != nil {
				return fmt.Errorf("failed to delete module '%s': %v", moduleName, err)
			}
			delete(cache.Modules, moduleName)
			cacheUpdated = true
		}
	}

	// Update the cache file if changes were made
	if cacheUpdated {
		if err := writeCache(cacheFilePath, cache); err != nil {
			return fmt.Errorf("failed to update cache file: %v", err)
		}
	}

	return nil
}
