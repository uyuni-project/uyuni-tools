// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Expected values for AssertSccFlag.
var SccFlagTestArgs = []string{
	"--scc-user", "mysccuser",
	"--scc-password", "mysccpass",
}

// Assert that all SCC flags are parsed correctly.
func AssertSccFlag(t *testing.T, cmd *cobra.Command, flags *types.SCCCredentials) {
	test_utils.AssertEquals(t, "Error parsing --scc-user", "mysccuser", flags.User)
	test_utils.AssertEquals(t, "Error parsing --scc-password", "mysccpass", flags.Password)
}
