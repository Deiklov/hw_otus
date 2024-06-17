package main

import "testing"

func TestRunCmd(t *testing.T) {
	// Define a test case where we expect the command to succeed
	cmd := []string{"echo", "Hello, world!"}
	env := Environment{} // Assuming Environment is a map or struct you've defined elsewhere

	// Run the command
	returnCode := RunCmd(cmd, env)

	// Check if the return code indicates success
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	// You can add more test cases here, for example, testing with different commands,
	// or testing how it handles commands that should fail.
}
