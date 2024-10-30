// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// SccFlagTestArgs is the expected values for AssertSccFlag.
var SccFlagTestArgs = []string{
	"--scc-user", "mysccuser",
	"--scc-password", "mysccpass",
}

// AssertSccFlag checks that all SCC flags are parsed correctly.
func AssertSccFlag(t *testing.T, cmd *cobra.Command, flags *types.SCCCredentials) {
	testutils.AssertEquals(t, "Error parsing --scc-user", "mysccuser", flags.User)
	testutils.AssertEquals(t, "Error parsing --scc-password", "mysccpass", flags.Password)
}
