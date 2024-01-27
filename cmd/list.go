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
		err := myApp.List()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
