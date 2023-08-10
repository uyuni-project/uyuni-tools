package migrate

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Image    string
	ImageTag string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	migrateCmd := &cobra.Command{
		Use:   "migrate [source server FQDN]",
		Short: "migrate a remote server to containers",
		Long: `Migrate a remote server to containers

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * the SSH configuration should stored in $HOME/.ssh, including the keys
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

NOTE: for now installing on a remote cluster or podman is not supported yet!
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			command := utils.GetCommand()
			switch command {
			case "podman":
				migrateToPodman(globalFlags, flags, cmd, args)
			case "kubectl":
				migrateToKubernetes(globalFlags, flags, cmd, args)
			}
		},
	}

	// TODO We probably want to move these default values to a config file
	migrateCmd.Flags().StringVar(&flags.Image, "image", "registry.opensuse.org/uyuni/server", "Image")
	migrateCmd.Flags().StringVar(&flags.ImageTag, "tag", "latest", "Tag Image")

	return migrateCmd
}
