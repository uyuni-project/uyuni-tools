// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// SCCFlagTestArgs is the expected values for AssertSccFlag.
var SCCFlagTestArgs = []string{
	"--scc-user", "mysccuser",
	"--scc-password", "mysccpass",
}

// AssertSCCFlag checks that all SCC flags are parsed correctly.
func AssertSCCFlag(t *testing.T, flags *types.SCCCredentials) {
	testutils.AssertEquals(t, "Error parsing --scc-user", "mysccuser", flags.User)
	testutils.AssertEquals(t, "Error parsing --scc-password", "mysccpass", flags.Password)
}
