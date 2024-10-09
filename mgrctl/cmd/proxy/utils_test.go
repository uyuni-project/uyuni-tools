// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy_test

import (
	"path"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/proxy"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

// Test getFilename function.
func TestGetFilename(t *testing.T) {
	// Test when output is empty
	filename := proxy.GetFilename("", "testProxy.domain.com")
	test_utils.AssertEquals(t, "", "testProxy-config.tar.gz", filename)

	// Test when output is provided
	filename = proxy.GetFilename("customOutput", "testProxy.domain.com")
	test_utils.AssertEquals(t, "", "customOutput.tar.gz", filename)

	// Test when output is provided
	filename = proxy.GetFilename("/var/customOutputWitPath", "testProxy.domain.com")
	test_utils.AssertEquals(t, "", "/var/customOutputWitPath.tar.gz", filename)
}

func createTestFile(dir string, filename string, content string, t *testing.T) string {
	filepath := path.Join(dir, filename)
	test_utils.WriteFile(t, filepath, content)
	return filepath
}
