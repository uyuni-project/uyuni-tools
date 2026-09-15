// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package gpg

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type apiFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags apiFlags

	gpgCmd := &cobra.Command{
		Use:   "gpg",
		Short: L("GPG key management commands"),
	}

	gpgUploadKeyCmd := &cobra.Command{
		Use:   "upload [path or URL]",
		Short: L("Upload a GPG key to the server"),
		Long:  L(`Uploads an armored GPG key to the server from a local file or a remote URL.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runGpgKeyUpload)
		},
		Args: cobra.ExactArgs(1),
	}

	gpgListKeysCmd := &cobra.Command{
		Use:   "list",
		Short: L("List all GPG keys on the server"),
		Long:  L(`Retrieves a list of registered and/or previously uploaded GPG keys.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runGpgKeyList)
		},
	}

	gpgRemoveKeyCmd := &cobra.Command{
		Use:   "remove [fingerprint]",
		Short: L("Remove a GPG key from the server"),
		Long:  L(`Removes a key identified by its fingerprint from the server's GPG keyring.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runGpgKeyRemove)
		},
		Args: cobra.ExactArgs(1),
	}

	gpgCmd.AddCommand(gpgUploadKeyCmd)
	gpgCmd.AddCommand(gpgListKeysCmd)
	gpgCmd.AddCommand(gpgRemoveKeyCmd)
	api.AddAPIFlags(gpgCmd)

	return gpgCmd
}
