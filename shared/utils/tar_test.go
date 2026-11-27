// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
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
		err := os.Mkdir(subDir, 0700)
		testutils.AssertNoError(t, "failed to create "+dirPath+" directory", err)
	}

	// Add some content to the data directory
	for name, content := range filesData {
		filePath := path.Dir(name)
		if filePath != "." {
			absDir := path.Join(dir, dataDir, filePath)
			err := os.MkdirAll(absDir, 0700)
			testutils.AssertNoError(t, "failed to create subdirectory "+absDir+" for test", err)
		}
		testutils.WriteFile(t, path.Join(dir, dataDir, name), content)
	}

	return dir
}

func TestWriteTarGz(t *testing.T) {
	tmpDir := setup(t)

	// Create the tarball
	tarballPath := path.Join(tmpDir, "test.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	testutils.AssertNoError(t, "failed to create tarball", err)
	err = tarball.AddFile(path.Join(tmpDir, dataDir, "file1"), "otherfile1")
	testutils.AssertNoError(t, "failed to add file1 to tarball", err)
	err = tarball.AddFile(path.Join(tmpDir, dataDir, "sub/file2"), "sub/file2")
	testutils.AssertNoError(t, "failed to add sub/file2 to tarball", err)
	tarball.Close()

	// Check the tarball using the tar utility
	testDir := path.Join(tmpDir, outDir)
	out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput()
	testutils.AssertNoError(t, "failed to extract generated tarball: "+string(out), err)

	// Ensure we have all expected files
	for _, file := range []string{"otherfile1", "sub/file2"} {
		testutils.AssertTrue(t, "Missing "+file+" in archive", FileExists(path.Join(testDir, file)))
	}

	// Check the content of a file
	testutils.AssertEquals(
		t, "expected otherfile1 content", file1Content, testutils.ReadFile(t, path.Join(testDir, "otherfile1")),
	)
}

func TestWriteTarGzDirectory(t *testing.T) {
	tmpDir := setup(t)

	// Create the tarball
	tarballPath := path.Join(tmpDir, "test-dir.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	testutils.AssertNoError(t, "failed to create tarball", err)
	err = tarball.AddFile(path.Join(tmpDir, dataDir), "collected")
	testutils.AssertNoError(t, "failed to add directory to tarball", err)
	tarball.Close()

	// Check the tarball using the tar utility
	testDir := path.Join(tmpDir, outDir)
	out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput()
	testutils.AssertNoError(t, "failed to extract generated tarball: "+string(out), err)

	// Ensure all files from the source directory are present in the archive path
	for name, content := range filesData {
		outputPath := path.Join(testDir, "collected", name)
		testutils.AssertEquals(t, "expected content for "+outputPath, content, testutils.ReadFile(t, outputPath))
	}
}

func TestWriteTarGzDirectoryDoesNotFollowSymlink(t *testing.T) {
	tmpDir := setup(t)

	dataPath := path.Join(tmpDir, dataDir)
	err := os.Symlink(".", path.Join(dataPath, "loop"))
	testutils.AssertNoError(t, "failed to create test symlink", err)

	tarballPath := path.Join(tmpDir, "test-symlink.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	testutils.AssertNoError(t, "failed to create tarball", err)
	err = tarball.AddFile(dataPath, "collected")
	testutils.AssertNoError(t, "failed to add directory with symlink to tarball", err)
	tarball.Close()

	testDir := path.Join(tmpDir, outDir)
	out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput()
	testutils.AssertNoError(t, "failed to extract generated tarball: "+string(out), err)

	loopPath := path.Join(testDir, "collected", "loop")
	info, err := os.Lstat(loopPath)
	testutils.AssertNoError(t, "failed to stat extracted symlink", err)
	testutils.AssertTrue(t, "expected symlink: "+loopPath, info.Mode()&os.ModeSymlink != 0)

	link, err := os.Readlink(loopPath)
	testutils.AssertNoError(t, "failed to read extracted symlink target", err)
	testutils.AssertEquals(t, "expected symlink target '.'", ".", link)
}

func TestExtractTarGz(t *testing.T) {
	tmpDir := setup(t)

	// Create an archive using the tar tool
	tarballPath := path.Join(tmpDir, "test.tar.gz")
	dataPath := path.Join(tmpDir, dataDir)
	out, err := exec.Command("tar", "czf", tarballPath, "-C", dataPath, ".").CombinedOutput()
	testutils.AssertNoError(t, "failed to create test tar.gz: "+string(out), err)

	// Extract the tarball
	testDir := path.Join(tmpDir, outDir)
	err = ExtractTarGz(tarballPath, testDir)
	testutils.AssertNoError(t, "Failed to extract tar.gz", err)

	// Check the extracted content
	for name, content := range filesData {
		testutils.AssertEquals(t, "expected content for "+name, content, testutils.ReadFile(t, path.Join(testDir, name)))
	}
}
