package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ListType string

const (
	ModulesType      ListType = "modules"
	DependenciesType ListType = "dependencies"
)

var validListTypes = []ListType{ModulesType, DependenciesType}

func isValidListType(listType ListType) bool {
	for _, v := range validListTypes {
		if v == listType {
			return true
		}
	}
	return false
}

var listCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List modules/dependencies",
	Long:  "List all modules in the cache or dependencies in the project",
	Run: func(cmd *cobra.Command, args []string) {
		myApp := getApp()
		listType := ListType(args[0])
		if !isValidListType(listType) {
			fmt.Println("Invalid list type. Valid options are 'modules' and 'dependencies'.")
			return
		}
		if listType == ModulesType {
			for _, mod := range myApp.GetModuleList() {
				fmt.Println(mod)
			}
		} else {
			err := myApp.List()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
