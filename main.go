package main

import (
	"fmt"
	"gogetty/cmd"
	"gogetty/pkg/gogetty" // Import your gogetty package
	"os"
)

func main() {

	// Validate the environment before executing any commands
	if err := gogetty.ValidateEnvironment(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	cmd.Execute()
}
