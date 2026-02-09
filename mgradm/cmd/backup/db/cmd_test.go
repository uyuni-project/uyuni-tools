// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestEnableParamsParsing(t *testing.T) {
	args := []string{"--force"}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *Flagpole,
		_ *cobra.Command, _ []string,
	) error {
		testutils.AssertTrue(t, "Force flag not true", flags.Force)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newDBEnableCmd(&globalFlags, tester)

	testutils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}

func TestDisableParamsParsing(t *testing.T) {
	args := []string{"--force", "--purge-volume"}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *Flagpole,
		_ *cobra.Command, _ []string,
	) error {
		testutils.AssertTrue(t, "Force flag not true", flags.Force)
		testutils.AssertTrue(t, "PurgeVolume flag not true", flags.Purge.Volume)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newDBDisableCmd(&globalFlags, tester)

	testutils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}

func TestStatusParamsParsing(t *testing.T) {
	args := []string{}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *Flagpole,
		_ *cobra.Command, _ []string,
	) error {
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newDBStatusCmd(&globalFlags, tester)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}

func TestRestoreParamsParsing(t *testing.T) {
	args := []string{"--force"}

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *Flagpole,
		_ *cobra.Command, _ []string,
	) error {
		testutils.AssertTrue(t, "Force flag not true", flags.Force)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newDBRestoreCmd(&globalFlags, tester)

	testutils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}
