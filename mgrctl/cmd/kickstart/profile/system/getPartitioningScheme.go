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

type getPartitioningSchemeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
}

func getPartitioningSchemeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPartitioningScheme",
		Short: "Get the partitioning scheme for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPartitioningSchemeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPartitioningScheme)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of a kickstart profile")

	return cmd
}

func getPartitioningScheme(globalFlags *types.GlobalFlags, flags *getPartitioningSchemeFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
