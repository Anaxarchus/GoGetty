package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a dependency",
	Long:  "Remove a dependency from the project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		myApp := getApp()
		if err := myApp.Remove(name); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Dependency '%s' removed successfully\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
