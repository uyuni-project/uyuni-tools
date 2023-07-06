package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate"
)

// NewCommand returns a new cobra.Command implementing the root command for kinder
func NewUyuniadmCommand() *cobra.Command {
	globalFlags := &types.GlobalFlags{}
	rootCmd := &cobra.Command{
		Use:     "uyuniadm",
		Short:   "Uyuni administration tool",
		Long:    "Uyuni administration tool used to help user administer uyuni servers on k8s and podman",
		Version: "0.0.1",
	}

	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(migrate.NewCommand(globalFlags))

	return rootCmd
}
