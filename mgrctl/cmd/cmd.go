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
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/api"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/cp"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/exec"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/proxy"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/term"
	"github.com/uyuni-project/uyuni-tools/shared/completion"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewUyunictlCommand returns a new cobra.Command implementing the root command for mgrctl.
func NewUyunictlCommand() *cobra.Command {
	globalFlags := &types.GlobalFlags{}
	name := path.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        L("Uyuni control tool"),
		Long:         L("Tool to help managing Uyuni servers mainly through their API"),
		Version:      utils.Version,
		SilenceUsage: true, // Don't show usage help on errors
	}

	rootCmd.SetUsageTemplate(utils.GetLocalizedUsageTemplate())

	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", L("configuration file path"))
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "logLevel", "",
		L("application log level")+"(trace|debug|info|warn|error|fatal|panic)",
	)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		// do not log if running the completion cmd as the output is redirect to create a file to source
		if cmd.Name() != "completion" {
			utils.LogInit((cmd.Name() != "exec" && cmd.Name() != "term") || globalFlags.LogLevel == "trace")
			utils.SetLogLevel(globalFlags.LogLevel)
			log.Info().Msgf(L("Starting %s"), strings.Join(os.Args, " "))
		}
	}

	apiCmd := api.NewCommand(globalFlags)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(exec.NewCommand(globalFlags))
	rootCmd.AddCommand(term.NewCommand(globalFlags))
	rootCmd.AddCommand(cp.NewCommand(globalFlags))
	rootCmd.AddCommand(completion.NewCommand(globalFlags))
	rootCmd.AddCommand(proxy.NewCommand(globalFlags))

	rootCmd.AddCommand(utils.GetConfigHelpCommand())

	return rootCmd
}
