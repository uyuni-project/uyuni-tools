// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

const dataDir = "data"
const outDir = "out"

const file1Content = "file1 content"

var filesData = map[string]string{
	"file1":     file1Content,
	"sub/file2": "file2 content",
}

// Prepare test files to include in the tarball.
func setup(t *testing.T) string {
	dir := t.TempDir()

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

	return dir
}

func TestWriteTarGz(t *testing.T) {
	tmpDir := setup(t)

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
	} else if string(out) != file1Content {
		t.Errorf("expected otherfile1 content %s, but got %s", file1Content, string(out))
	}
}

func TestWriteTarGzDirectory(t *testing.T) {
	tmpDir := setup(t)

	// Create the tarball
	tarballPath := path.Join(tmpDir, "test-dir.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	if err != nil {
		t.Fatalf("failed to create tarball: %s", err)
	}
	if err := tarball.AddFile(path.Join(tmpDir, dataDir), "collected"); err != nil {
		t.Fatalf("failed to add directory to tarball: %s", err)
	}
	tarball.Close()

	// Check the tarball using the tar utility
	testDir := path.Join(tmpDir, outDir)
	if out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput(); err != nil {
		t.Fatalf("failed to extract generated tarball: %s", string(out))
	}

	// Ensure all files from the source directory are present in the archive path
	for name, content := range filesData {
		outputPath := path.Join(testDir, "collected", name)
		if out, err := os.ReadFile(outputPath); err != nil {
			t.Errorf("failed to read %s: %s", outputPath, err)
		} else if string(out) != content {
			t.Errorf("expected %s content %s, but got %s", outputPath, content, string(out))
		}
	}
}

func TestWriteTarGzDirectoryDoesNotFollowSymlink(t *testing.T) {
	tmpDir := setup(t)

	dataPath := path.Join(tmpDir, dataDir)
	if err := os.Symlink(".", path.Join(dataPath, "loop")); err != nil {
		t.Fatalf("failed to create test symlink: %s", err)
	}

	tarballPath := path.Join(tmpDir, "test-symlink.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	if err != nil {
		t.Fatalf("failed to create tarball: %s", err)
	}
	if err := tarball.AddFile(dataPath, "collected"); err != nil {
		t.Fatalf("failed to add directory with symlink to tarball: %s", err)
	}
	tarball.Close()

	testDir := path.Join(tmpDir, outDir)
	if out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput(); err != nil {
		t.Fatalf("failed to extract generated tarball: %s", string(out))
	}

	loopPath := path.Join(testDir, "collected", "loop")
	info, err := os.Lstat(loopPath)
	if err != nil {
		t.Fatalf("failed to stat extracted symlink: %s", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected %s to be a symlink", loopPath)
	}

	link, err := os.Readlink(loopPath)
	if err != nil {
		t.Fatalf("failed to read extracted symlink target: %s", err)
	}
	if link != "." {
		t.Fatalf("expected symlink target '.' but got %q", link)
	}
}

func TestExtractTarGz(t *testing.T) {
	tmpDir := setup(t)

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
