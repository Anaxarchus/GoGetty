package cmd

import (
	"fmt"
	"gogetty/pkg/gogetty"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch dependencies",
	Long:  "Fetch all dependencies defined in your project",
	Run: func(cmd *cobra.Command, args []string) {
		manager := gogetty.GogettyManager{ProjectDir: projectDir}
		if err := manager.Fetch(); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Dependencies fetched successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().StringVar(&projectDir, "project-dir", ".", "Specify the project directory")
}
