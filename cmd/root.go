package cmd

import (
	"bufio"
	"fmt"
	"gogetty/pkg/gitwrap"
	"os"
	"strings"

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
- Update a dependency: gogetty update <moduleName> [--branch <branchName>] [--commit <commitHash>] [--directory <commaSeperatedDirectories>]
- Remove a dependency: gogetty remove <moduleName>
- Clean up dependencies, and remove unused modules from the cache: gogetty clean
- Fetch all dependencies, downloading missing modules to the cache, and creating symbolic links: gogetty fetch
- List dependencies or modules: gogetty list [modules/dependencies]

- Run git commands on your modules: gogetty [moduleName] <gitCommand> <gitParameters>`,
}

func promptUserForAction() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("There's a conflict with the command. Do you want to execute the command as GoGetty (1) or as Git (2)? [1/2]: ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

func Execute() {
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cmdInput := args[0]
			myApp := getApp()
			find := gitwrap.GitRepo{
				Name: cmdInput,
			}
			module := gitwrap.FindInList(find, myApp.Cache)

			// Check if the command is a valid Cobra command
			_, _, err := cmd.Find(args)

			if module != nil {
				// Case 1: Valid module, always forward
				if err != nil {
					return gitwrap.Forward(module.Path, cmdInput, args[1:]...)
				}

				// Case 2: Valid module and valid command, ask the user
				userChoice, promptErr := promptUserForAction()
				if promptErr != nil {
					return promptErr
				}
				if strings.ToUpper(userChoice) == "2" {
					return gitwrap.Forward(module.Path, cmdInput, args[1:]...)
				}
			}
		}

		_, err := cmd.ExecuteC()
		return err
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
