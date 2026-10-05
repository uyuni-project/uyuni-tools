// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package disable

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	tester := func(_ *types.GlobalFlags, _ *disableFlags, _ *cobra.Command, _ []string) error {
		return nil
	}
	cmd := newCmd(&types.GlobalFlags{}, tester)
	testutils.AssertHasAllFlags(t, cmd, []string{})
	cmd.SetArgs([]string{})
	testutils.AssertNoError(t, "command failed: ", cmd.Execute())
}

func TestRejectsArguments(t *testing.T) {
	cmd := newCmd(&types.GlobalFlags{}, func(_ *types.GlobalFlags, _ *disableFlags, _ *cobra.Command, _ []string) error {
		return nil
	})
	cmd.SetArgs([]string{"unexpected"})
	testutils.AssertError(t, "accepts 0 arg(s)", cmd.Execute())
}
