package cmd

import (
	"fmt"
	"gogetty/pkg/gogetty"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  "Initialize a new Gogetty project in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		manager := gogetty.GogettyManager{ProjectDir: projectDir}
		if err := manager.Init(); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Gogetty project initialized successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&projectDir, "project-dir", ".", "Specify the project directory")
}
