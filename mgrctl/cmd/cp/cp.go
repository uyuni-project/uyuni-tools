// SPDX-FileCopyrightText: 2024 SUSE LLC
// SPDX-FileCopyrightText: 2026 YatesMold
//
// SPDX-License-Identifier: Apache-2.0

package cp

import (
	"fmt"
	"strings"

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
		Use:   "cp [path/to/source1.file ...] [path/to/destination]",
		Short: L("Copy files to and from the containers"),
		Long: L(`Takes at least one source and destination parameters.
	Several source files can be passed, the last parameter will be the destination.
	Either the source parameters or the destination can be prefixed with 'server:' to indicate the path is within the server pod.`),
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper, err := utils.ReadConfig(cmd, utils.GlobalConfigFilename, globalFlags.ConfigPath)
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

func run(flags *flagpole, _ *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	dst := args[len(args)-1]
	srcs := args[:len(args)-1]

	// Validate that source paths don't mix local and server paths.
	hasServerSrc := strings.HasPrefix(srcs[0], "server:")
	for _, src := range srcs[1:] {
		if strings.HasPrefix(src, "server:") != hasServerSrc {
			return fmt.Errorf(L("cannot mix local and server sources"))
		}
	}

	for _, src := range srcs {
		if err := cnx.Copy(src, dst, flags.User, flags.Group); err != nil {
			return err
		}
	}
	return nil
}
