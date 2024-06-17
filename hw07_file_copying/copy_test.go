package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	// Setup - create a temporary directory
	tempDir, err := os.MkdirTemp("", "copy_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up

	srcFilePath := filepath.Join(tempDir, "src.txt")
	content := []byte("Hello, Gopher!")
	if err := os.WriteFile(srcFilePath, content, 0o644); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}

	// Define the destination file path
	destFilePath := filepath.Join(tempDir, "dest.txt")

	t.Run("copy entire file", func(t *testing.T) {
		if err := Copy(srcFilePath, destFilePath, 0, 0); err != nil {
			t.Errorf("Copy failed: %v", err)
		}

		// Verify the content of the destination file
		destContent, err := os.ReadFile(destFilePath)
		if err != nil {
			t.Fatalf("Failed to read destination file: %v", err)
		}
		if !bytes.Equal(content, destContent) {
			t.Errorf("Content mismatch; got %s, want %s", destContent, content)
		}
	})

	// Test handling of ErrOffsetExceedsFileSize
	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy(srcFilePath, destFilePath, int64(len(content)+1), 0)
		if !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Errorf("Expected ErrOffsetExceedsFileSize; got %v", err)
		}
	})
}
