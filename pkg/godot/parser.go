package godot

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// WalkProjectDirectory walks the entire project directory and generates script lists.
func walkProjectDirectory(dirPath string) ([]GodotScript, []GodotScript, error) {
	gdScripts := []GodotScript{}
	csScripts := []GodotScript{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".gd" {
			gdScripts = append(gdScripts, GodotScript{
				Path: path,
				Type: "gd",
			})
		} else if filepath.Ext(path) == ".cs" {
			csScripts = append(csScripts, GodotScript{
				Path: path,
				Type: "cs",
			})
		}
		return nil
	})

	return gdScripts, csScripts, err
}

func parseScriptPaths(script *GodotScript, project GodotProject) error {
	// Open the script file for reading
	file, err := os.Open(script.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a scanner to read the script line by line
	scanner := bufio.NewScanner(file)
	var modifiedLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Replace "res://" paths with projectDir
		line = strings.ReplaceAll(line, "res://", project.Path)

		// Replace "user://" paths with the UserDirectory
		line = strings.ReplaceAll(line, "user://", project.UserDirectory)

		// Append the modified line to the result
		modifiedLines = append(modifiedLines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Reopen the script file for writing
	file, err = os.Create(script.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the modified lines back to the script file
	for _, line := range modifiedLines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
