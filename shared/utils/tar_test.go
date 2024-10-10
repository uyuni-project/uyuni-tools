// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

const dataDir = "data"
const outDir = "out"

const file1_content = "file1 content"

var filesData = map[string]string{
	"file1":     file1_content,
	"sub/file2": "file2 content",
}

// Prepare test files to include in the tarball.
func setup(t *testing.T) (string, func(t *testing.T)) {
	dir, clean := test_utils.CreateTmpFolder(t)

	// Create sub directories for the data and the test
	for _, dirPath := range []string{dataDir, outDir} {
		subDir := path.Join(dir, dirPath)
		if err := os.Mkdir(subDir, 0700); err != nil {
			t.Fatalf("failed to create %s directory: %s", dirPath, err)
		}
	}

	// Add some content to the data directory
	for name, content := range filesData {
		filePath := path.Dir(name)
		if filePath != "." {
			absDir := path.Join(dir, dataDir, filePath)
			if err := os.MkdirAll(absDir, 0700); err != nil {
				t.Fatalf("failed to create subdirectory %s for test: %s", absDir, err)
			}
		}
		if err := os.WriteFile(path.Join(dir, dataDir, name), []byte(content), 0700); err != nil {
			t.Fatalf("failed to write test data file %s: %s", name, err)
		}
	}

	// Returns the teardown function.
	return dir, func(t *testing.T) {
		clean()
	}
}

func TestWriteTarGz(t *testing.T) {
	tmpDir, teardown := setup(t)
	defer teardown(t)

	// Create the tarball
	tarballPath := path.Join(tmpDir, "test.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	if err != nil {
		t.Fatalf("failed to create tarball: %s", err)
	}
	if err := tarball.AddFile(path.Join(tmpDir, dataDir, "file1"), "otherfile1"); err != nil {
		t.Fatalf("failed to add file1 to tarball: %s", err)
	}
	if err := tarball.AddFile(path.Join(tmpDir, dataDir, "sub/file2"), "sub/file2"); err != nil {
		t.Fatalf("failed to add sub/file2 to tarball: %s", err)
	}
	tarball.Close()

	// Check the tarball using the tar utility
	testDir := path.Join(tmpDir, outDir)
	if out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput(); err != nil {
		t.Fatalf("failed to extract generated tarball: %s", string(out))
	}

	// Ensure we have all expected files
	for _, file := range []string{"otherfile1", "sub/file2"} {
		if !FileExists(path.Join(testDir, file)) {
			t.Errorf("Missing %s in archive", file)
		}
	}

	// Check the content of a file
	if out, err := os.ReadFile(path.Join(testDir, "otherfile1")); err != nil {
		t.Errorf("failed to read otherfile1: %s", err)
	} else if string(out) != file1_content {
		t.Errorf("expected otherfile1 content %s, but got %s", file1_content, string(out))
	}
}

func TestExtractTarGz(t *testing.T) {
	tmpDir, teardown := setup(t)
	defer teardown(t)

	// Create an archive using the tar tool
	tarballPath := path.Join(tmpDir, "test.tar.gz")
	dataPath := path.Join(tmpDir, dataDir)
	if out, err := exec.Command("tar", "czf", tarballPath, "-C", dataPath, ".").CombinedOutput(); err != nil {
		t.Fatalf("failed to create test tar.gz: %s", string(out))
	}

	// Extract the tarball
	testDir := path.Join(tmpDir, outDir)
	if err := ExtractTarGz(tarballPath, testDir); err != nil {
		t.Errorf("Failed to extract tar.gz: %s", err)
	}

	// Check the extracted content
	for name, content := range filesData {
		if out, err := os.ReadFile(path.Join(testDir, name)); err != nil {
			t.Errorf("failed to read %s: %s", name, err)
		} else if string(out) != content {
			t.Errorf("expected %s content %s, but got %s", name, content, string(out))
		}
	}
}
