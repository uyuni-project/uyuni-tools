// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/cache"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/install"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/logs"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/restart"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/start"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/status"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/stop"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/support"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/uninstall"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd/upgrade"
	"github.com/uyuni-project/uyuni-tools/shared/completion"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder.
func NewUyuniproxyCommand() (*cobra.Command, error) {
	globalFlags := &types.GlobalFlags{}
	name := path.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        L("Uyuni proxy administration tool"),
		Long:         L("Tool to help administering Uyuni proxies in containers"),
		Version:      utils.Version,
		SilenceUsage: true, // Don't show usage help on errors
	}

	rootCmd.AddGroup(&cobra.Group{
		ID:    "deploy",
		Title: L("Server Deployment:"),
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "management",
		Title: L("Server Management:"),
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "tool",
		Title: L("Administrator tools:"),
	})

	rootCmd.SetUsageTemplate(utils.GetLocalizedUsageTemplate())

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		utils.LogInit(true)
		utils.SetLogLevel(globalFlags.LogLevel)

		// do not log if running the completion cmd as the output is redirected to create a file to source
		if cmd.Name() != "completion" && cmd.Name() != "__complete" {
			log.Info().Msgf(L("Welcome to %s"), name)
			log.Info().Msgf(L("Executing command: %s"), cmd.Name())
		}
	}

	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", L("configuration file path"))
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "logLevel", "", L("application log level")+"(trace|debug|info|warn|error|fatal|panic)")

	installCmd := install.NewCommand(globalFlags)
	rootCmd.AddCommand(installCmd)
	uninstallCmd, err := uninstall.NewCommand(globalFlags)
	if err != nil {
		return rootCmd, err
	}
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(completion.NewCommand(globalFlags))
	rootCmd.AddCommand(cache.NewCommand(globalFlags))
	rootCmd.AddCommand(status.NewCommand(globalFlags))
	rootCmd.AddCommand(start.NewCommand(globalFlags))
	rootCmd.AddCommand(stop.NewCommand(globalFlags))
	rootCmd.AddCommand(restart.NewCommand(globalFlags))
	rootCmd.AddCommand(upgrade.NewCommand(globalFlags))
	rootCmd.AddCommand(logs.NewCommand(globalFlags))

	if supportCommand := support.NewCommand(globalFlags); supportCommand != nil {
		rootCmd.AddCommand(supportCommand)
	}

	rootCmd.AddCommand(utils.GetConfigHelpCommand())

	return rootCmd, nil
}
