package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List dependencies",
	Long:  "List all dependencies in the project",
	Run: func(cmd *cobra.Command, args []string) {
		myApp := getApp()
		dependencies, err := myApp.List()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Dependencies:")
			for _, dep := range dependencies {
				fmt.Printf("Name: %s\nBranch: %s\nCommit: %s\n\n", dep.Repository.Name, dep.Repository.Branch, dep.Repository.Commit)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
