// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestWriteTarGz(t *testing.T) {
	// Create the tarball
	tarballPath := path.Join(testDir, "test.tar.gz")
	tarball, err := NewTarGz(tarballPath)
	if err != nil {
		t.Fatalf("failed to create tarball: %s", err)
	}
	if err := tarball.AddFile(path.Join(testDir, dataDir, "tar_file1"), "tar_otherfile1"); err != nil {
		t.Fatalf("failed to add tar_file1 to tarball: %s", err)
	}
	if err := tarball.AddFile(path.Join(testDir, dataDir, "sub/tar_file2"), "sub/tar_file2"); err != nil {
		t.Fatalf("failed to add sub/tar_file2 to tarball: %s", err)
	}
	tarball.Close()

	// Check the tarball using the tar utility
	testDir := path.Join(testDir, outDir)
	if out, err := exec.Command("tar", "xzf", tarballPath, "-C", testDir).CombinedOutput(); err != nil {
		t.Fatalf("failed to extract generated tarball: %s", string(out))
	}

	// Ensure we have all expected files
	for _, file := range []string{"tar_otherfile1", "sub/tar_file2"} {
		if !FileExists(path.Join(testDir, file)) {
			t.Errorf("Missing %s in archive", file)
		}
	}

	// Check the content of a file
	if out, err := os.ReadFile(path.Join(testDir, "tar_otherfile1")); err != nil {
		t.Errorf("failed to read tar_otherfile1: %s", err)
	} else if string(out) != tarFile1Content {
		t.Errorf("expected tar_otherfile1 content %s, but got %s", tarFile1Content, string(out))
	}
}

func TestExtractTarGz(t *testing.T) {
	// Create an archive using the tar tool
	tarballPath := path.Join(testDir, "test.tar.gz")
	dataPath := path.Join(testDir, dataDir)
	if out, err := exec.Command("tar", "czf", tarballPath, "-C", dataPath, ".").CombinedOutput(); err != nil {
		t.Fatalf("failed to create test tar.gz: %s", string(out))
	}

	// Extract the tarball
	testDir := path.Join(testDir, outDir)
	if err := ExtractTarGz(tarballPath, testDir); err != nil {
		t.Errorf("Failed to extract tar.gz: %s", err)
	}

	// Check the extracted content
	for name, content := range tarFilesData {
		if out, err := os.ReadFile(path.Join(testDir, name)); err != nil {
			t.Errorf("failed to read %s: %s", name, err)
		} else if string(out) != content {
			t.Errorf("expected %s content %s, but got %s", name, content, string(out))
		}
	}
}
