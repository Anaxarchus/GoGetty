package cmd

import (
	"fmt"
	"gogetty/install"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs GoGetty to your system, adding it to your system PATH",
	Long: `Installs the GoGetty executable to the cache directory and 
adds it to the PATH environment variable, making it 
accessible from anywhere in your command line interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := install.Install("gogetty"); err != nil {
			fmt.Printf("Error installing GoGetty: %v\n", err)
		} else {
			fmt.Println("GoGetty installed successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
