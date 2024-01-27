package main

import (
	"fmt"
	"gogetty/cmd"
	"gogetty/pkg/app"
	"os"
)

func main() {

	// Validate the environment before executing any commands
	if err := app.ValidateEnvironment(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	cmd.Execute()
}
