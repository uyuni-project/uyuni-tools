package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setSoftwareDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	$param.getFlagName()          $param.getType()
}

func setSoftwareDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setSoftwareDetails",
		Short: "Sets kickstart profile software details.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setSoftwareDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setSoftwareDetails)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of the kickstart profile")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setSoftwareDetails(globalFlags *types.GlobalFlags, flags *setSoftwareDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.KsLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

