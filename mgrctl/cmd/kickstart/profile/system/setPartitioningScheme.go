package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setPartitioningSchemeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
}

func setPartitioningSchemeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPartitioningScheme",
		Short: "Set the partitioning scheme for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setPartitioningSchemeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setPartitioningScheme)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of the kickstart profile to update")

	return cmd
}

func setPartitioningScheme(globalFlags *types.GlobalFlags, flags *setPartitioningSchemeFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
