package cmd

import (
	"fmt"
	"gogetty/pkg/gogetty"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a dependency",
	Long:  "Remove a dependency from the project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: gogetty remove <name>")
			return
		}
		name := args[0]
		manager := gogetty.GogettyManager{ProjectDir: projectDir}
		if err := manager.Remove(name); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Dependency '%s' removed successfully\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVar(&projectDir, "project-dir", ".", "Specify the project directory")
}
