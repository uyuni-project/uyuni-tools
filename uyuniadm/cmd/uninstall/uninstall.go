package uninstall

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall a server",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			purge, _ := cmd.Flags().GetBool("purge-volumes")

			command := utils.GetCommand()
			switch command {
			case "podman":
				uninstallForPodman(globalFlags, dryRun, purge)
			case "kubectl":
				uninstallForKubernetes(globalFlags, dryRun)
			}
		},
	}
	uninstallCmd.Flags().BoolP("dry-run", "n", false, "Only show what would be done")
	uninstallCmd.Flags().Bool("purge-volumes", false, "Also remove the volume (podman only)")

	return uninstallCmd
}
