package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/install"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/uninstall"
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
	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigPath, "config", "c", "", "configuration file path")

	migrateCmd := migrate.NewCommand(globalFlags)
	addCommonFlags(migrateCmd)
	rootCmd.AddCommand(migrateCmd)

	installCmd := install.NewCommand(globalFlags)
	addCommonFlags(installCmd)
	rootCmd.AddCommand(installCmd)

	rootCmd.AddCommand(uninstall.NewCommand(globalFlags))

	return rootCmd
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("podman-arg", []string{}, "Extra arguments to pass to podman, separated by commas")

	cmd.Flags().String("helm-uyuni-namespace", "default", "Kubernetes namespace where to install uyuni")
	cmd.Flags().String("helm-uyuni-chart", "oci://registry.opensuse.org/uyuni/server", "URL to the uyuni helm chart")
	cmd.Flags().String("helm-uyuni-version", "", "Version of the uyuni helm chart")
	cmd.Flags().String("helm-uyuni-values", "", "Path to a values YAML file to use for Uyuni helm install")
	cmd.Flags().String("helm-certmanager-namespace", "cert-manager", "Kubernetes namespace where to install cert-manager")
	cmd.Flags().String("helm-certmanager-chart", "", "URL to the cert-manager helm chart. To be used for offline installations")
	cmd.Flags().String("helm-certmanager-version", "", "Version of the cert-manager helm chart")
	cmd.Flags().String("helm-certmanager-values", "", "Path to a values YAML file to use for cert-manager helm install")
}
