package godot

import (
	"fmt"
	"path/filepath"
	"runtime"

	"gopkg.in/ini.v1"
)

type GodotScript struct {
	LineCount int
	Path      string
	Type      string
}

type GodotProject struct {
	Name                string
	GodotVersion        string
	Path                string
	CustomUserDirectory bool
	UserDirectory       string
	Features            []string
	Scripts             []GodotScript
}

func (gp *GodotProject) SetUserDirectory() {
	if gp.CustomUserDirectory {
		// Set custom directory path
		if gp.UserDirectory != "" {
			// Custom dir and name
			gp.UserDirectory = getUserDirPath(gp.UserDirectory)
		} else {
			// Custom dir only
			gp.UserDirectory = getUserDirPath(gp.Name)
		}
	} else {
		// Set default directory path
		gp.UserDirectory = getDefaultUserDirPath(gp.Name)
	}
}

func getDefaultUserDirPath(projectName string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join("%APPDATA%", "Godot", "app_userdata", projectName)
	case "darwin":
		return filepath.Join("~/Library/Application Support/Godot/app_userdata", projectName)
	default: // Linux and other OS
		return filepath.Join("~/.local/share/godot/app_userdata", projectName)
	}
}

func getUserDirPath(dirName string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join("%APPDATA%", dirName)
	case "darwin":
		return filepath.Join("~/Library/Application Support", dirName)
	default: // Linux and other OS
		return filepath.Join("~/.local/share", dirName)
	}
}

func UpdateProjectPaths(project GodotProject) error {
	for _, script := range project.Scripts {
		err := parseScriptPaths(&script, project)
		if err != nil {
			fmt.Println("Godot Script failed to parse: ", filepath.Base(script.Path))
		}
	}
	return nil
}

func GetGodotProject(dirPath string) (*GodotProject, error) {
	projectFilePath := filepath.Join(dirPath, "project.godot")
	cfg, err := ini.Load(projectFilePath)
	if err != nil {
		return nil, err
	}

	var godotProject GodotProject

	// Extract values with default handling
	if name, err := cfg.Section("application").GetKey("config/name"); err == nil {
		godotProject.Name = name.String()
	}
	if useCustomDir, err := cfg.Section("application").GetKey("config/use_custom_user_dir"); err == nil {
		godotProject.CustomUserDirectory = useCustomDir.MustBool(false)
	}
	if customDirName, err := cfg.Section("application").GetKey("config/custom_user_dir_name"); err == nil {
		godotProject.UserDirectory = customDirName.String()
	} else if godotProject.CustomUserDirectory {
		godotProject.UserDirectory = godotProject.Name
	}

	// Set user directory based on the parsed values
	godotProject.SetUserDirectory()

	// Walk the project directory to generate script lists
	gdScripts, csScripts, err := walkProjectDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	// Append the script lists to the GodotProject's Scripts field
	godotProject.Scripts = append(godotProject.Scripts, gdScripts...)
	godotProject.Scripts = append(godotProject.Scripts, csScripts...)

	return &godotProject, nil
}
