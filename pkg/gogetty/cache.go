package gogetty

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Cache struct {
	Modules map[string]ModuleInfo `json:"modules"`
}

type ModuleInfo struct {
	Path       string   `json:"path"`
	URL        string   `json:"url"`
	Branch     string   `json:"branch"`
	Commit     string   `json:"commit"`
	Dependents []string `json:"dependents"`
}

const CacheDir = ".gogetty"
const CacheJson = "cache.json"
const ModuleDir = "modules"

func GetCachePath(subpath ...string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory.", err)
	}

	return filepath.Join(append([]string{homeDir, CacheDir}, subpath...)...)
}

func InitCache() error {
	modulesDir := GetCachePath(ModuleDir)
	cacheFile := GetCachePath(CacheJson)

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(modulesDir, 0755); err != nil {
		return err
	}

	// Create or open the cache file
	file, err := os.OpenFile(cacheFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize cache structure
	cache := Cache{Modules: make(map[string]ModuleInfo)}
	json.NewEncoder(file).Encode(cache)

	return nil
}

func UpdateCacheWithModule(gitURL, branch, commit string) error {
	cacheFilePath := GetCachePath(CacheJson)

	// Load existing cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return err
	}

	// Update cache with new module information
	moduleName := filepath.Base(gitURL) // Or derive the name in another way
	moduleInfo := ModuleInfo{
		Path:   filepath.Join(GetCachePath(ModuleDir), moduleName),
		Branch: branch,
		Commit: commit,
	}
	cache.Modules[moduleName] = moduleInfo

	// Write updated cache back to file
	return writeCache(cacheFilePath, cache)
}

func loadCache(path string) (Cache, error) {
	var cache Cache
	file, err := os.Open(path)
	if err != nil {
		return cache, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cache)
	return cache, err
}

func writeCache(path string, cache Cache) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cache)
}

func ValidateModule(url, branch, commit string) (bool, error) {
	moduleName := filepath.Base(url) // Derive the module name from the URL

	// Check if the module exists on the filesystem
	exists, err := ModuleExists(moduleName)
	if err != nil {
		return false, fmt.Errorf("error checking if module exists: %v", err)
	}
	if !exists {
		return false, nil // The module does not exist on the filesystem
	}

	// Check if the module is in the cache with the specified details
	inCache, err := ModuleInCache(url, branch, commit)
	if err != nil {
		return false, fmt.Errorf("error checking if module is in cache: %v", err)
	}

	return inCache, nil // The module is valid if it's in the cache with the specified details
}

func ModuleExists(moduleName string) (bool, error) {
	cacheFilePath := GetCachePath(CacheJson)

	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return false, err
	}

	moduleInfo, exists := cache.Modules[moduleName]
	if !exists {
		return false, nil
	}

	_, err = os.Stat(moduleInfo.Path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func ModuleInCache(url, branch, commit string) (bool, error) {
	cacheFilePath := GetCachePath(CacheJson)

	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return false, err
	}

	for _, module := range cache.Modules {
		if module.URL == url && (branch == "" || module.Branch == branch) && (commit == "" || module.Commit == commit) {
			return true, nil
		}
	}

	return false, nil
}

func ModuleHasDependency(moduleName string, dependencyPath string) (bool, error) {
	cacheFilePath := GetCachePath(CacheJson)

	// Load the cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return false, err
	}

	// Check if the specified module exists and iterate through its dependents
	if moduleInfo, exists := cache.Modules[moduleName]; exists {
		for _, dependent := range moduleInfo.Dependents {
			if dependent == dependencyPath {
				return true, nil
			}
		}
	}

	return false, nil
}

func AddModuleToCache(moduleName, url, branch, commit string) error {
	cacheFilePath := GetCachePath(CacheJson)
	modulePath := GetCachePath(filepath.Join(ModuleDir, moduleName))

	// Load existing cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return err
	}

	// Update or add the module information
	moduleInfo, exists := cache.Modules[moduleName]
	if exists {
		// Update existing module info
		moduleInfo.URL = url
		moduleInfo.Branch = branch
		moduleInfo.Commit = commit
	} else {
		// Add new module info
		moduleInfo = ModuleInfo{
			Path:   modulePath,
			URL:    url,
			Branch: branch,
			Commit: commit,
		}
		cache.Modules[moduleName] = moduleInfo
	}

	// Write updated cache back to file
	return writeCache(cacheFilePath, cache)
}

func DeleteModuleFromCache(moduleName string) error {
	cacheFilePath := GetCachePath(CacheJson)
	modulePath := GetCachePath(filepath.Join(ModuleDir, moduleName))

	// Load existing cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return err
	}

	// Check if the module exists in cache
	if _, exists := cache.Modules[moduleName]; !exists {
		return fmt.Errorf("module '%s' not found in cache", moduleName)
	}

	// Delete the module directory
	err = os.RemoveAll(modulePath)
	if err != nil {
		return fmt.Errorf("failed to delete module '%s': %v", moduleName, err)
	}

	// Update cache by removing the module entry
	delete(cache.Modules, moduleName)

	// Write updated cache back to file
	return writeCache(cacheFilePath, cache)
}

func AddDependentToModule(moduleName, dependent string) error {
	cacheFilePath := GetCachePath(CacheJson)

	// Load existing cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return err
	}

	// Add dependent to the module
	module, exists := cache.Modules[moduleName]
	if !exists {
		return fmt.Errorf("module '%s' not found in cache", moduleName)
	}
	if !contains(module.Dependents, dependent) {
		module.Dependents = append(module.Dependents, dependent)
		cache.Modules[moduleName] = module
	}

	// Write updated cache back to file
	return writeCache(cacheFilePath, cache)
}

func RemoveDependentFromModule(moduleName, dependent string) error {
	cacheFilePath := GetCachePath(CacheJson)

	// Load existing cache
	cache, err := loadCache(cacheFilePath)
	if err != nil {
		return err
	}

	// Remove dependent from the module
	module, exists := cache.Modules[moduleName]
	if !exists {
		return fmt.Errorf("module '%s' not found in cache", moduleName)
	}
	module.Dependents = removeElement(module.Dependents, dependent)
	cache.Modules[moduleName] = module

	// Write updated cache back to file
	return writeCache(cacheFilePath, cache)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeElement(slice []string, item string) []string {
	result := []string{}
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
