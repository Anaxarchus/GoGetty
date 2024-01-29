package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message + " [y/N]: ")

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return false
	}
	response = strings.TrimSpace(response)
	return strings.ToLower(response) == "y"
}

func Option(message string, options ...string) int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("\t%s [%d]\n", option, i)
	}

	for {
		fmt.Print("Enter your choice: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice >= len(options) {
			fmt.Println("Invalid selection. Please enter a valid index.")
			continue
		}

		return choice
	}
}
