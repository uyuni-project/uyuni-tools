// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"path"
	"strings"

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

// NewUyuniproxyCommand returns a new cobra.Command implementing the root command for mgrpxy.
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

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		// do not log if running the completion cmd as the output is redirected to create a file to source
		if cmd.Name() != "completion" && cmd.Name() != "__complete" {
			utils.LogInit(true)
			utils.SetLogLevel(globalFlags.LogLevel)
			utils.SetShouldPreserveTmpDir(globalFlags.KeepTempDir)
			log.Info().Msgf(L("Starting %s"), strings.Join(os.Args, " "))
			log.Info().Msgf(L("Use of this software implies acceptance of the End User License Agreement."))
		}
	}

	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", L("configuration file path"))
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.KeepTempDir, "keepTemp", "", false,
		L("keep temporary directories for debugging purpose"))
	if err := rootCmd.PersistentFlags().MarkHidden("keepTemp"); err != nil {
		log.Warn().Err(err).Msg("Failed to hide keepTemp flag")
	}
	utils.AddLogLevelFlags(rootCmd, &globalFlags.LogLevel)

	installCmd := install.NewCommand(globalFlags)
	rootCmd.AddCommand(installCmd)
	uninstallCmd := uninstall.NewCommand(globalFlags)
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
