// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/api"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/cp"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/exec"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/org"
	"github.com/uyuni-project/uyuni-tools/shared/completion"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder.
func NewUyunictlCommand() (*cobra.Command, error) {
	globalFlags := &types.GlobalFlags{}
	name := path.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        "Uyuni control tool",
		Long:         "Uyuni control tool used to help user managing Uyuni and SUSE Manager Servers mainly through its API",
		Version:      utils.Version,
		SilenceUsage: true, // Don't show usage help on errors
	}

	usage, err := utils.GetUsageWithConfigHelpTemplate(rootCmd.UsageTemplate())
	if err != nil {
		return rootCmd, err
	}
	rootCmd.SetUsageTemplate(usage)

	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", "configuration file path")
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "logLevel", "", "application log level (trace|debug|info|warn|error|fatal|panic)")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		utils.LogInit(cmd.Name() != "exec")
		utils.SetLogLevel(globalFlags.LogLevel)

		// do not log if running the completion cmd as the output is redirect to create a file to source
		if cmd.Name() != "completion" {
			log.Info().Msgf("Welcome to %s", name)
			log.Info().Msgf("Executing command: %s", cmd.Name())
		}
	}

	apiCmd, err := api.NewCommand(globalFlags)
	if err != nil {
		//FIXME this should return err, but with it the code stop compiling
		return rootCmd, nil
	}
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(exec.NewCommand(globalFlags))
	rootCmd.AddCommand(cp.NewCommand(globalFlags))
	rootCmd.AddCommand(completion.NewCommand(globalFlags))
	orgCmd, err := org.NewCommand(globalFlags)
	if err != nil {
		return rootCmd, err
	}
	rootCmd.AddCommand(orgCmd)

	return rootCmd, nil
}
