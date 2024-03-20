// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/restart"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/start"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/stop"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/support"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/uninstall"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder.
func NewUyuniadmCommand() (*cobra.Command, error) {
	globalFlags := &types.GlobalFlags{}
	name := path.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        "Uyuni administration tool",
		Long:         "Uyuni administration tool used to help user administer uyuni servers on kubernetes and podman",
		Version:      utils.Version,
		SilenceUsage: true, // Don't show usage help on errors
	}

	usage, err := utils.GetUsageWithConfigHelpTemplate(rootCmd.UsageTemplate())
	if err != nil {
		return rootCmd, err
	}
	rootCmd.SetUsageTemplate(usage)

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
	distroCmd, err := distro.NewCommand(globalFlags)
	if err != nil {
		return rootCmd, err
	}
	rootCmd.AddCommand(distroCmd)
	rootCmd.AddCommand(completion.NewCommand(globalFlags))
	rootCmd.AddCommand(support.NewCommand(globalFlags))
	rootCmd.AddCommand(start.NewCommand(globalFlags))
	rootCmd.AddCommand(hub.NewCommand(globalFlags))
	rootCmd.AddCommand(restart.NewCommand(globalFlags))
	rootCmd.AddCommand(stop.NewCommand(globalFlags))
	rootCmd.AddCommand(inspect.NewCommand(globalFlags))
	rootCmd.AddCommand(upgrade.NewCommand(globalFlags))

	return rootCmd, err
}
