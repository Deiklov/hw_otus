package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Check if there is at least one command argument
	if len(cmd) == 0 {
		return -1 // Return an error code if no command is provided
	}

	// Prepare the command
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	// Proxy stdin, stdout, and stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Set up the environment for the command
	//nolint:prealloc
	var newEnv []string
	// Append the current system environment to preserve system settings

	newEnv = append(newEnv, os.Environ()...)
	for key, val := range env {
		if val.NeedRemove {
			val.Value = ""
		}
		newEnv = append(newEnv, key+"="+val.Value)
	}

	command.Env = newEnv

	// Run the command
	if err := command.Run(); err != nil {
		// If there's an error, including a non-zero exit code, it will be of type *ExitError
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		// If we can't get the exit code, return a default error code
		return -1
	}

	// If the command completes successfully, return 0
	return 0
}
