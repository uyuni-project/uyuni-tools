// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

// APIFlagsTestArgs is the slice of parameters to use with AssertAPIFlags.
var APIFlagsTestArgs = []string{
	"--api-server", "mysrv",
	"--api-user", "apiuser",
	"--api-password", "api-pass",
	"--api-cacert", "path/to/ca.crt",
	"--api-insecure",
}

// AssertAPIFlags checks that all API parameters are parsed correctly.
func AssertAPIFlags(t *testing.T, cmd *cobra.Command, flags *api.ConnectionDetails) {
	test_utils.AssertEquals(t, "Error parsing --api-server", "mysrv", flags.Server)
	test_utils.AssertEquals(t, "Error parsing --api-user", "apiuser", flags.User)
	test_utils.AssertEquals(t, "Error parsing --api-password", "api-pass", flags.Password)
	test_utils.AssertEquals(t, "Error parsing --api-cacert", "path/to/ca.crt", flags.CApath)
	test_utils.AssertTrue(t, "Error parsing --api-insecure", flags.Insecure)
}
