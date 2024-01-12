// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/completion"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/distro"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/uninstall"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder
func NewUyuniadmCommand() *cobra.Command {
	globalFlags := &types.GlobalFlags{}
	name := path.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        "Uyuni administration tool",
		Long:         "Uyuni administration tool used to help user administer uyuni servers on kubernetes and podman",
		Version:      utils.Version,
		SilenceUsage: true, // Don't show usage help on errors
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		utils.LogInit(true)
		utils.SetLogLevel(globalFlags.LogLevel)

		// do not log if running the completion cmd as the output is redirected to create a file to source
		if cmd.Name() != "completion" {
			log.Info().Msgf("Welcome to %s", name)
			log.Info().Msgf("Executing command: %s", cmd.Name())
		}
	}

	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", "configuration file path")
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "logLevel", "", "application log level (trace|debug|info|warn|error|fatal|panic)")

	migrateCmd := migrate.NewCommand(globalFlags)
	rootCmd.AddCommand(migrateCmd)

	installCmd := install.NewCommand(globalFlags)
	rootCmd.AddCommand(installCmd)

	rootCmd.AddCommand(uninstall.NewCommand(globalFlags))
	rootCmd.AddCommand(distro.NewCommand(globalFlags))
	rootCmd.AddCommand(completion.NewCommand(globalFlags))

	return rootCmd
}
