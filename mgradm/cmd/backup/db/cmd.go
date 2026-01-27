// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// CLI definitions

func NewDBCmd(globalFlags *types.GlobalFlags) *cobra.Command {
	dbCmd := &cobra.Command{
		Use:   "db",
		Short: L("Database backup management"),
		Long:  L("Tools for online database backup management"),
	}
	dbCmd.AddCommand(newDBEnableCmd(globalFlags, doDBEnable))
	dbCmd.AddCommand(newDBDisableCmd(globalFlags, doDBDisable))
	dbCmd.AddCommand(newDBStatusCmd(globalFlags, doDBStatus))
	dbCmd.AddCommand(newDBRestoreCmd(globalFlags, doDBRestore))
	return dbCmd
}

func newDBEnableCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[Flagpole]) *cobra.Command {
	var flags Flagpole
	cmd := &cobra.Command{
		Use:   "enable",
		Short: L("Enable continuous archiving backup"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	cmd.Flags().Bool("force", false, L("Reconfigure already configured backup"))

	return cmd
}

func newDBDisableCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[Flagpole]) *cobra.Command {
	var flags Flagpole
	cmd := &cobra.Command{
		Use:   "disable",
		Short: L("Disable continuous archiving backup"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	cmd.Flags().Bool("force", false, L("Don't ask for confirmation when purging volume"))
	cmd.Flags().Bool("purge-volume", false, L("Also remove the volume"))

	return cmd
}

func newDBStatusCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[Flagpole]) *cobra.Command {
	var flags Flagpole
	cmd := &cobra.Command{
		Use:   "status",
		Short: L("Check WAL based database backup status"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	return cmd
}

func newDBRestoreCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[Flagpole]) *cobra.Command {
	var flags Flagpole
	cmd := &cobra.Command{
		Use:   "restore",
		Short: L("Restore database backup"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}
	cmd.Flags().Bool("force", false, L("Don't ask for confirmation"))
	return cmd
}

// Actual actions

func doDBEnable(
	_ *types.GlobalFlags,
	flags *Flagpole,
	_ *cobra.Command,
	_ []string,
) error {
	return Enable(flags.Force)
}

func doDBDisable(
	_ *types.GlobalFlags,
	flags *Flagpole,
	_ *cobra.Command,
	_ []string,
) error {
	return Disable(flags)
}

func doDBStatus(
	_ *types.GlobalFlags,
	_ *Flagpole,
	_ *cobra.Command,
	_ []string,
) error {
	return Status()
}

func doDBRestore(
	_ *types.GlobalFlags,
	flags *Flagpole,
	_ *cobra.Command,
	_ []string,
) error {
	return Restore(flags.Force)
}
