package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "envtest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test.

	// Set up test cases.
	testCases := []struct {
		filename    string
		content     string
		expectedEnv Environment
	}{
		{"VAR1", "value1\n", Environment{"VAR1": EnvValue{"value1", false}}},
		{"VAR2", "", Environment{"VAR2": EnvValue{"", true}}},
		{"VAR_EMPTY_LINE", "\n", Environment{"VAR_EMPTY_LINE": EnvValue{"", false}}},
		{"ignore=me", "should be ignored", nil},
	}

	// Create files in the temp directory according to test cases.
	for _, tc := range testCases {
		filePath := filepath.Join(tempDir, tc.filename)
		if err := os.WriteFile(filePath, []byte(tc.content), 0o666); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}

	env, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("ReadDir returned an error: %v", err)
	}

	// Verify the results.
	for _, tc := range testCases {
		if tc.expectedEnv == nil {
			continue // Skip files that should be ignored.
		}
		if val, ok := env[tc.filename]; !ok || !reflect.DeepEqual(val, tc.expectedEnv[tc.filename]) {
			t.Errorf("For file %s, expected %v, got %v", tc.filename, tc.expectedEnv[tc.filename], val)
		}
	}
}
