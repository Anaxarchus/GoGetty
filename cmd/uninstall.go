package cmd

import (
	"fmt"
	"gogetty/install"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls GoGetty from your system, removing it from your system PATH",
	Long: `Uninstalls the GoGetty executable from the cache directory and 
removes it from the PATH environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := install.Uninstall("gogetty"); err != nil {
			fmt.Printf("Error uninstalling GoGetty: %v\n", err)
		} else {
			fmt.Println("GoGetty uninstalled successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
