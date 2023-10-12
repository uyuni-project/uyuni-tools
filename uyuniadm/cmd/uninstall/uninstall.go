package uninstall

import (
	"github.com/rs/zerolog/log"
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

			// TODO Change to subcommands like other uyuniadm commands
			cnx := utils.NewConnection("")
			command, err := cnx.GetCommand()
			if err != nil {
				log.Fatal().Err(err)
			}
			switch command {
			case "podman":
				uninstallForPodman(dryRun, purge)
			case "kubectl":
				uninstallForKubernetes(dryRun)
			}
		},
	}
	uninstallCmd.Flags().BoolP("dry-run", "n", false, "Only show what would be done")
	uninstallCmd.Flags().Bool("purge-volumes", false, "Also remove the volume (podman only)")

	return uninstallCmd
}
