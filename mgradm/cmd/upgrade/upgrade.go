package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	upgradeCmd := &cobra.Command{
		Use:   "upgrade server",
		Short: "upgrade local server",
		Long:  "Upgrade local server",
	}

	upgradeCmd.AddCommand(podman.NewCommand(globalFlags))

	if kubernetesCmd := kubernetes.NewCommand(globalFlags); kubernetesCmd != nil {
		upgradeCmd.AddCommand(kubernetesCmd)
	}

	return upgradeCmd
}
