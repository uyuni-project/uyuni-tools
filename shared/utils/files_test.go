// SPDX-FileCopyrightText: 2025 Massimo Ambrosi
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTempDirRemovesDirectoryByDefault(t *testing.T) {
	// Ensure preserve flag is false (default state)
	SetShouldPreserveTmpDir(false)

	tempDir, cleaner, err := TempDir()
	if err != nil {
		t.Fatalf("TempDir() failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(tempDir); err != nil {
		t.Fatalf("TempDir() did not create directory: %v", err)
	}

	// Create a test file inside to verify removal works
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Call the cleaner function
	cleaner()

	// Verify directory was removed
	if _, err := os.Stat(tempDir); err == nil {
		t.Errorf("TempDir cleaner did not remove directory: %s still exists", tempDir)
	} else if !os.IsNotExist(err) {
		t.Errorf("Unexpected error checking directory: %v", err)
	}
}

func TestTempDirPreservesDirectoryWhenFlagIsSet(t *testing.T) {
	// Set preserve flag to true
	SetShouldPreserveTmpDir(true)
	defer SetShouldPreserveTmpDir(false) // Reset after test

	tempDir, cleaner, err := TempDir()
	if err != nil {
		t.Fatalf("TempDir() failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(tempDir); err != nil {
		t.Fatalf("TempDir() did not create directory: %v", err)
	}

	// Create a test file inside to verify it's preserved
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("test content")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Call the cleaner function
	cleaner()

	// Verify directory still exists
	if _, err := os.Stat(tempDir); err != nil {
		t.Errorf("TempDir cleaner removed directory when it should be preserved: %v", err)
	}

	// Verify the test file still exists and has correct content
	if _, err := os.Stat(testFile); err != nil {
		t.Errorf("Test file was removed when directory should be preserved: %v", err)
	} else {
		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Errorf("Failed to read preserved test file: %v", err)
		} else if string(content) != string(testContent) {
			t.Errorf("Preserved file content mismatch: got %s, want %s", content, testContent)
		}
	}

	// Clean up manually for this test
	os.RemoveAll(tempDir)
}

func TestTempDirFlagStateIsCapturedAtCreation(t *testing.T) {
	// Test that the flag value is captured when TempDir() is called,
	// not when cleaner() is called

	// Start with preserve = false
	SetShouldPreserveTmpDir(false)

	tempDir1, cleaner1, err := TempDir()
	if err != nil {
		t.Fatalf("TempDir() failed: %v", err)
	}

	// Change flag to true
	SetShouldPreserveTmpDir(true)
	defer SetShouldPreserveTmpDir(false)

	// Create another temp dir with flag = true
	tempDir2, cleaner2, err := TempDir()
	if err != nil {
		t.Fatalf("TempDir() failed: %v", err)
	}

	// Call cleaners
	cleaner1() // Should remove (created when flag was false)
	cleaner2() // Should preserve (created when flag was true)

	// Verify first directory was removed
	if _, err := os.Stat(tempDir1); err == nil {
		t.Errorf("First temp dir should have been removed: %s still exists", tempDir1)
	}

	// Verify second directory was preserved
	if _, err := os.Stat(tempDir2); err != nil {
		t.Errorf("Second temp dir should have been preserved: %v", err)
	} else {
		// Clean up manually
		os.RemoveAll(tempDir2)
	}
}
