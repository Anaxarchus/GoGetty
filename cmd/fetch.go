package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch dependencies",
	Long:  "Fetch all dependencies defined in your project",
	Run: func(cmd *cobra.Command, args []string) {
		myApp := getApp()

		// Perform the fetch operation
		if err := myApp.Fetch(); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Dependencies fetched successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
