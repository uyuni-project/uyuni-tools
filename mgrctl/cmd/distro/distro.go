// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type distroUploadFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name                  string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags distroUploadFlags

	distroCmd := &cobra.Command{
		Use:   "distro",
		Short: L("Distro management commands"),
	}

	distroUploadCmd := &cobra.Command{
		Use:   "upload [path or URL]",
		Short: L("Upload a distro ISO to the server"),
		Long: L(`Uploads a distro ISO to the server from a local file or a remote URL.

The distro is registered automatically if the ISO contains a .treeinfo file.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, runDistroUpload)
		},
		Args: cobra.ExactArgs(1),
	}
	distroUploadCmd.Flags().String("name", "",
		L("Distribution name, defaults to the ISO file name without .iso"))

	distroCmd.AddCommand(distroUploadCmd)
	api.AddAPIFlags(distroCmd)

	return distroCmd
}
