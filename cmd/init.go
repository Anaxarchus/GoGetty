package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  "Initialize a new Gogetty project in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		myApp := getApp()
		if err := myApp.Init(); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Gogetty project initialized successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
