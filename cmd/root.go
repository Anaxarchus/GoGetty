package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Declare flags at the package level
var (
	branchFlag     string
	commitFlag     string
	directoryFlags []string
)

var rootCmd = &cobra.Command{
	Use:   "gogetty",
	Short: "GoGetty is a versatile dependency manager for projects using Git repositories.",
	Long: `GoGetty simplifies managing dependencies directly from Git repositories. 
Versioning is supported through the use of commit hashes.

Examples of using GoGetty:

- Initialize a new project: gogetty init
- Add a dependency to your project: gogetty add <git-repo-url> [--branch <branchName>] [--commit <commitHash>] [--directory <commaSeperatedDirectories>]
- Update a dependency: gogetty update <dependencyName> [--branch <branchName>] [--commit <commitHash>] [--directory <commaSeperatedDirectories>]
- Remove a dependency: gogetty remove <dependencyName>
- Clean up dependencies, and remove unused modules from the cache: gogetty clean
- Fetch all dependencies, downloading missing modules to the cache, and creating symbolic links: gogetty fetch`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
