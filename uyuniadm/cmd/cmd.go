package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/install"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/uninstall"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder
func NewUyuniadmCommand() *cobra.Command {
	utils.LogInit("uyuniadm")
	globalFlags := &types.GlobalFlags{}
	rootCmd := &cobra.Command{
		Use:     "uyuniadm",
		Short:   "Uyuni administration tool",
		Long:    "Uyuni administration tool used to help user administer uyuni servers on k8s and podman",
		Version: "0.0.1",
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		utils.SetLogLevel(globalFlags.LogLevel)
		log.Info().Msgf("Executing command: %s", cmd.Name())
	}

	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", "configuration file path")
	rootCmd.PersistentFlags().StringVar(&globalFlags.LogLevel, "logLevel", "", "application log level (trace|debug|info|warn|error|fatal|panic)")

	migrateCmd := migrate.NewCommand(globalFlags)
	rootCmd.AddCommand(migrateCmd)

	installCmd := install.NewCommand(globalFlags)
	rootCmd.AddCommand(installCmd)

	rootCmd.AddCommand(uninstall.NewCommand(globalFlags))

	return rootCmd
}
