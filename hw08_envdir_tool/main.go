package main

import (
	"log"
	"os"
)

func main() {
	// Ensure there are at least 4 arguments provided
	if len(os.Args) < 5 {
		log.Fatal("Error: Insufficient arguments provided. Expected at least 4.")
	}

	// Read the directory specified in the first argument
	environmentVariables, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err)
	}

	// Execute the command with the provided arguments and environment variables
	exitCode := RunCmd(os.Args[2:], environmentVariables)

	// Exit the program with the command's exit code
	os.Exit(exitCode)
}
