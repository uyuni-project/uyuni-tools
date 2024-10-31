// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"path"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

// Test getFilename function.
func TestGetFilename(t *testing.T) {
	// Test when output is empty
	filename := GetFilename("", "testProxy.domain.com")
	testutils.AssertEquals(t, "", "testProxy-config.tar.gz", filename)

	// Test when output is provided
	filename = GetFilename("customOutput", "testProxy.domain.com")
	testutils.AssertEquals(t, "", "customOutput.tar.gz", filename)

	// Test when output is provided
	filename = GetFilename("/var/customOutputWitPath", "testProxy.domain.com")
	testutils.AssertEquals(t, "", "/var/customOutputWitPath.tar.gz", filename)
}

func createTestFile(dir string, filename string, content string, t *testing.T) string {
	filepath := path.Join(dir, filename)
	testutils.WriteFile(t, filepath, content)
	return filepath
}
