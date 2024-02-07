package tree

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/tree"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	TreeLabel             string
	BasePath              string
	ChannelLabel          string
	InstallType           string
	InstallType           string
	KernelOptions         string
	PostKernelOptions     string
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Kickstart Tree (Distribution) in #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("TreeLabel", "", "The new kickstart tree label.")
	cmd.Flags().String("BasePath", "", "Path to the base or root of the kickstart tree.")
	cmd.Flags().String("ChannelLabel", "", "Label of channel to associate with the kickstart tree. ")
	cmd.Flags().String("InstallType", "", "Label for KickstartInstallType (rhel_6, rhel_7, rhel_8, rhel_9, fedora_9).")
	cmd.Flags().String("InstallType", "", "Label for KickstartInstallType (rhel_2.1, rhel_3, rhel_4, rhel_5, fedora_9).")
	cmd.Flags().String("KernelOptions", "", "Options to be passed to the kernel when booting for the installation. ")
	cmd.Flags().String("PostKernelOptions", "", "Options to be passed to the kernel when booting for the installation. ")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

	res, err := tree.Tree(&flags.ConnectionDetails, flags.TreeLabel, flags.BasePath, flags.ChannelLabel, flags.InstallType, flags.InstallType, flags.KernelOptions, flags.PostKernelOptions)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
