// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
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
func AssertAPIFlags(t *testing.T, flags *api.ConnectionDetails) {
	testutils.AssertEquals(t, "Error parsing --api-server", "mysrv", flags.Server)
	testutils.AssertEquals(t, "Error parsing --api-user", "apiuser", flags.User)
	testutils.AssertEquals(t, "Error parsing --api-password", "api-pass", flags.Password)
	testutils.AssertEquals(t, "Error parsing --api-cacert", "path/to/ca.crt", flags.CApath)
	testutils.AssertTrue(t, "Error parsing --api-insecure", flags.Insecure)
}
