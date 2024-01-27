package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a dependency",
	Long: `Add a new dependency to the project. Optionally specify a branch, commit, 
and specific directories within the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: gogetty add <url> [--branch branchName] [--commit commitHash] [--directory subdirPath]...")
			return
		}
		url := args[0]

		myApp := getApp()

		if err := myApp.Add(url, branchFlag, commitFlag, directoryFlags); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Dependency added successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&branchFlag, "branch", "", "Specify the branch of the repository")
	addCmd.Flags().StringVar(&commitFlag, "commit", "", "Specify the commit hash of the repository")
	addCmd.Flags().StringSliceVar(&directoryFlags, "directory", nil, "Specify subdirectories within the repository")
}
