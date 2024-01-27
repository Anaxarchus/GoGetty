package cmd

import (
	"fmt"
	"gogetty/pkg/app"

	"github.com/spf13/cobra"
)

var myApp *app.MyApp

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a dependency",
	Long:  "Remove a dependency from the project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := myApp.Remove(name); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Dependency '%s' removed successfully\n", name)
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		moduleNames, err := myApp.ListModuleNames()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return moduleNames, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	myApp = getApp()
	rootCmd.AddCommand(removeCmd)
}
