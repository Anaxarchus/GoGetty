package cmd

import (
	"fmt"
	"gogetty/pkg/gogetty"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up the project dependencies",
	Long:  "Verify and clean up every dependency in the cache.json, removing any non-existent dependencies and modules.",
	Run: func(cmd *cobra.Command, args []string) {
		manager := gogetty.GogettyManager{ProjectDir: projectDir}
		if err := manager.Clean(); err != nil {
			fmt.Println("Error cleaning dependencies:", err)
		} else {
			fmt.Println("Dependencies cleaned successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVar(&projectDir, "project-dir", ".", "Specify the project directory")
}
