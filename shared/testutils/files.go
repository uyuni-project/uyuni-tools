// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"os"
	"testing"
)

// WriteFile writes the content in a file at the given path and fails if anything wrong happens.
func WriteFile(t *testing.T, path string, content string) {
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("failed to write test file %s: %s", path, err)
	}
}

// ReadFile returns the content of a file as a string and fails is anything wrong happens.
func ReadFile(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", path, err)
	}
	return string(content)
}

// ReadFileAsBinary returns the content of a file as a slice of int8.
func ReadFileAsBinary(t *testing.T, path string) []int8 {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", path, err)
	}

	int8Content := make([]int8, len(content))
	for i, b := range content {
		int8Content[i] = int8(b)
	}

	return int8Content
}
