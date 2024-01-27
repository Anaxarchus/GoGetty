package install

import (
	"bufio"
	"fmt"
	"gogetty/pkg/cache"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const bashRC = ".bashrc"
const unixBin = "${HOME}/bin"

func Install(executableName string) error {
	var destPath string

	switch runtime.GOOS {
	case "windows":
		return addToPathWindows()
	default:
		// For Unix-like systems
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Expand environment variables in the path
		unixBinPath := os.ExpandEnv(unixBin)
		destPath = filepath.Join(unixBinPath, executableName)

		// Ensure the directory exists
		if err := os.MkdirAll(unixBinPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		// Copy the file
		if err := copyFile(execPath, destPath); err != nil {
			return fmt.Errorf("failed to copy executable to bin: %w", err)
		}

		// Set executable permissions
		if err := os.Chmod(destPath, 0755); err != nil {
			return fmt.Errorf("failed to set executable permissions: %w", err)
		}

		// Add to PATH
		return addToPathUnix(unixBinPath)
	}
}

// Uninstall removes the GoGetty executable from the cache directory and removes it from the PATH.
func Uninstall(executableName string) error {
	unixBinPath := os.ExpandEnv(unixBin)
	destPath := filepath.Join(unixBinPath, executableName)
	cachePath := cache.CacheDir()

	if err := os.Remove(destPath); err != nil {
		return fmt.Errorf("failed to remove executable from cache: %w", err)
	}

	return RemoveFromPath(cachePath)
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func addToPathWindows() error {
	instructions := `
To complete the installation of GoGetty on Windows, follow these steps:

1. Create folder C:\bin if it doesn't exist:
   Open Command Prompt and run: mkdir C:\bin

2. Save the GoGetty binary (gogetty.exe) to the folder C:\bin

3. Add C:\bin to your system PATH:
   - Press the Windows key and search for 'System' (Control Panel) for Windows 8 or 10, or right-click the Computer icon on the desktop and click Properties for Windows 7.
   - Click 'Advanced system settings'.
   - Click 'Environment Variables'.
   - Under 'System Variables', find the 'PATH' variable, select it, and click 'Edit'. If there is no 'PATH' variable, click 'New'.
   - Add 'C:\bin' to the start of the variable value, followed by a semicolon (;). For example, if the value was 'C:\Windows\System32', change it to 'C:\bin;C:\Windows\System32'.
   - Click 'OK' to save your changes.

4. Restart Command Prompt to apply the new PATH settings.

5. Verify that GoGetty is on your PATH:
   In Command Prompt, run: where.exe gogetty.exe
   You should see the path 'C:\bin\gogetty.exe' if the installation was successful.
`
	fmt.Println(instructions)
	return nil
}

func addToPathUnix(dir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	bashRCPath := filepath.Join(homeDir, ".bashrc") // Or other shell config files
	zshRCPath := filepath.Join(homeDir, ".zshrc")

	// You can add checks for the existence of these files and modify as needed
	err1 := modifyPath(bashRCPath, dir, true)
	err2 := modifyPath(zshRCPath, dir, true)
	if err1 != nil && err2 != nil {
		return fmt.Errorf("failed to update PATH in shell config files: %v, %v", err1, err2)
	}

	return nil
}

// RemoveFromPath removes the specified directory from the PATH for Ubuntu and Mint.
func RemoveFromPath(dir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Choose which file to modify based on the distribution or shell preference
	filePath := filepath.Join(homeDir, bashRC) // or homeDir/profile based on the user's environment

	return modifyPath(filePath, dir, false)
}

// modifyPath adds or removes the directory from the PATH in the specified file.
func modifyPath(filePath, dir string, add bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	pathLine := fmt.Sprintf("export PATH=\"$PATH:%s\"", dir)
	var lines []string
	pathFound := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, dir) {
			pathFound = true
			if !add {
				continue // Skip the line to remove it
			}
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if add && !pathFound {
		lines = append(lines, pathLine) // Add path line
	}

	return writeLinesToFile(lines, filePath)
}

// writeLinesToFile writes the given lines to the specified file.
func writeLinesToFile(lines []string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}
