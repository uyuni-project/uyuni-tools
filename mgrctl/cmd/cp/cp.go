// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cp

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	User    string
	Group   string
	Backend string
}

// NewCommand copy file to and from the containers.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	cpCmd := &cobra.Command{
		Use:   "cp [path/to/source.file] [path/to/destination.file]",
		Short: L("Copy files to and from the containers"),
		Long: L(`Takes a source and destination parameters.
	One of them can be prefixed with 'server:' to indicate the path is within the server pod.`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper, err := utils.ReadConfig(globalFlags.ConfigPath, cmd)
			if err != nil {
				return err
			}
			if err := viper.Unmarshal(&flags); err != nil {
				return utils.Errorf(err, L("failed to unmarshall configuration"))
			}
			return run(flags, cmd, args)
		},
	}

	cpCmd.Flags().String("user", "", L("User or UID to set on the destination file"))
	cpCmd.Flags().String("group", "susemanager", L("Group or GID to set on the destination file"))

	utils.AddBackendFlag(cpCmd)
	return cpCmd
}

func run(flags *flagpole, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	return cnx.Copy(args[0], args[1], flags.User, flags.Group)
}
