// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
)

// will hold the path to the testing directory for the utils package
var testDir string

const tmpDirPattern = "uyuni-tools-utils-test-"
const dataDir = "data"
const outDir = "out"

const tarFile1Content = "file1 content"

var tarFilesData = map[string]string{
	"tar_file1":     tarFile1Content,
	"sub/tar_file2": "file2 content",
}

func TestMain(m *testing.M) {
	teardown := setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

// Prepare test files and directories
func setup() func() {
	tmpDir, err := os.MkdirTemp("", tmpDirPattern)
	if err != nil {
		log.Fatalf("failed to create temporary directory for test: %s", err)
	}
	testDir = tmpDir

	// Create sub directories for the data and the test
	for _, dirPath := range []string{dataDir, outDir} {
		subDir := path.Join(testDir, dirPath)
		if err := os.Mkdir(subDir, 0700); err != nil {
			log.Fatalf("failed to create %s directory: %s", dirPath, err)
		}
	}

	if err := tarTestSetup(testDir); err != nil {
		log.Fatalf("failed to setup tar test: %s", err)
	}
	// more setups if needed be

	// return the teardown func
	return func() {
		if err := os.RemoveAll(testDir); err != nil {
			log.Fatal(err)
		}
	}
}

func tarTestSetup(dir string) error {
	// Add some content to the data directory
	for name, content := range tarFilesData {
		filePath := path.Dir(name)
		if filePath != "." {
			absDir := path.Join(dir, dataDir, filePath)
			if err := os.MkdirAll(absDir, 0700); err != nil {
				return fmt.Errorf("failed to create subdirectory %s for tar test: %s", absDir, err)
			}
		}
		if err := os.WriteFile(path.Join(dir, dataDir, name), []byte(content), 0700); err != nil {
			return fmt.Errorf("failed to write tar test data file %s: %s", name, err)
		}
	}
	return nil
}
