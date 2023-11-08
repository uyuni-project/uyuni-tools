package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyunictl/cmd/cp"
	"github.com/uyuni-project/uyuni-tools/uyunictl/cmd/distro"
	"github.com/uyuni-project/uyuni-tools/uyunictl/cmd/exec"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder
func NewUyunictlCommand() *cobra.Command {
	globalFlags := &types.GlobalFlags{}
	rootCmd := &cobra.Command{
		Use:     "uyunictl",
		Short:   "Uyuni control tool",
		Long:    "Uyuni control tool used to help user managing Uyuni and SUSE Manager Servers mainly through its API",
		Version: "0.1.0",
	}

	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", "configuration file path")
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "logLevel", "", "application log level (trace|debug|info|warn|error|fatal|panic)")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		utils.LogInit("uyunictl", cmd.Name() != "exec")
		utils.SetLogLevel(globalFlags.LogLevel)
		log.Info().Msgf("Executing command: %s", cmd.Name())
	}

	rootCmd.AddCommand(exec.NewCommand(globalFlags))
	rootCmd.AddCommand(cp.NewCommand(globalFlags))
	rootCmd.AddCommand(distro.NewCommand(globalFlags))

	return rootCmd
}
