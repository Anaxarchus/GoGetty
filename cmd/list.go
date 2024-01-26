package cmd

import (
	"fmt"
	"gogetty/pkg/gogetty"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List dependencies",
	Long:  "List all dependencies in the project",
	Run: func(cmd *cobra.Command, args []string) {
		manager := gogetty.GogettyManager{ProjectDir: projectDir}
		dependencies, err := manager.List()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Dependencies:")
			for _, dep := range dependencies {
				fmt.Printf("Name: %s\nURL: %s\nBranch: %s\nCommit: %s\n\n", dep.Name, dep.URL, dep.Branch, dep.Commit)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&projectDir, "project-dir", ".", "Specify the project directory")
}
