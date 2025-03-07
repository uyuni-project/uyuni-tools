// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package backup

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/create"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/restore"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newCreateCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[shared.Flagpole]) *cobra.Command {
	var flags shared.Flagpole

	createCmd := &cobra.Command{
		Use:   "create output-directory",
		Args:  cobra.ExactArgs(1),
		Short: L("Create backup"),
		Long:  L("Create backup of the already configured Uyuni system"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	createCmd.Flags().StringSlice("skipvolumes", []string{}, L("Skip backup of selected volumes"))
	createCmd.Flags().StringSlice("extravolumes", []string{}, L("Backup additional volumes to the build-in ones"))
	createCmd.Flags().Bool("skipdatabase", false, L("Do not backup database volume, allow online backup."))
	createCmd.Flags().Bool("skipimages", false, L("Do not backup container images"))
	createCmd.Flags().Bool("skipconfig", false, L("Do not backup podman configuration. On restore defaults will be used"))
	createCmd.Flags().Bool("norestart", false, L("Do not restart services after backup is done"))
	createCmd.Flags().Bool("dryrun", false, L("Print expected actions, but no action is done"))

	return createCmd
}

func newRestoreCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[shared.Flagpole]) *cobra.Command {
	var flags shared.Flagpole

	restoreCmd := &cobra.Command{
		Use:   "restore directory",
		Args:  cobra.ExactArgs(1),
		Short: L("Restore backup from the directory"),
		Long:  L("Restore backup of the previously configured Uyuni system from a specified directory"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	restoreCmd.Flags().StringSlice("skipvolumes", []string{}, L("Skip restore of selected volumes"))
	restoreCmd.Flags().Bool("skipdatabase", false, L("Do not restore database volume"))
	restoreCmd.Flags().Bool("skipimages", false, L("Skip restore of container images"))
	restoreCmd.Flags().Bool("skipconfig", false, L("Do not restore podman configuration. Defaults will be used"))
	restoreCmd.Flags().Bool("norestart", false, L("Do not restart service after restore is done"))
	restoreCmd.Flags().Bool("dryRun", false, L("Print expected actions, but no action is done"))
	restoreCmd.Flags().Bool("force", false, L("Force overwrite of existing items"))
	restoreCmd.Flags().Bool("continue", false, L("Skip existing items and restore the rest"))
	restoreCmd.Flags().Bool("skipverify", false, L("Skip verification of the backup files"))

	return restoreCmd
}

// NewCommand command for distribution management.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	backupCmd := &cobra.Command{
		Use:     "backup",
		GroupID: "tool",
		Short:   L("Backup solution"),
		Long:    L("Tools for local backup management"),
	}
	backupCmd.AddCommand(newCreateCmd(globalFlags, doBackup))
	backupCmd.AddCommand(newRestoreCmd(globalFlags, doRestore))
	return backupCmd
}

// Backup helper to catch errors with unified error message.
func doBackup(
	global *types.GlobalFlags,
	flags *shared.Flagpole,
	cmd *cobra.Command,
	args []string,
) error {
	outputDirectory := args[0]
	err := create.Create(global, flags, cmd, args)
	if err != nil {
		var backupError *shared.BackupError
		ok := errors.As(err, &backupError)
		if ok {
			log.Error().Msgf("%s", backupError.Err.Error())
			if backupError.Abort && backupError.DataRemains {
				return fmt.Errorf(L("Backup aborted, partially backed up files remains in '%s'"), outputDirectory)
			}
			if !backupError.Abort {
				// nolint:lll
				return errors.New(L("Important data were backed up successfully, but errors were present. Restore will use default values where needed"))
			}
		}
		return err
	}
	return nil
}

// Backup helper to catch errors with unified error message.
func doRestore(
	global *types.GlobalFlags,
	flags *shared.Flagpole,
	cmd *cobra.Command,
	args []string,
) error {
	err := restore.Restore(global, flags, cmd, args)
	if err != nil {
		var backupError *shared.BackupError
		ok := errors.As(err, &backupError)
		if ok {
			if backupError.Abort && backupError.DataRemains {
				log.Error().Msgf("%s", backupError.Err.Error())
				return errors.New(L("Restore aborted with partially restored files. Resolve the error and try again"))
			}
			if !backupError.Abort {
				return errors.New(L("Important data were restored successfully, but with warnings"))
			}
		}
		return err
	}
	return nil
}
