package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	newBranchFlag     string
	newCommitFlag     string
	newDirectoryFlags []string
)

var updateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update a dependency",
	Long:  "Update a dependency in the project. Optionally specify a new branch, new commit, and new directories.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: gogetty update <name> [--branch branchName] [--commit commitHash] [--directory subdirPath]...")
			return
		}
		name := args[0]
		myApp := getApp()
		if err := myApp.Update(name, newBranchFlag, newCommitFlag, newDirectoryFlags); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Dependency '%s' updated successfully\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&newBranchFlag, "branch", "", "Specify the new branch of the dependency")
	updateCmd.Flags().StringVar(&newCommitFlag, "commit", "", "Specify the new commit of the dependency")
	updateCmd.Flags().StringSliceVar(&newDirectoryFlags, "directory", nil, "Specify new subdirectories within the repository")
}
