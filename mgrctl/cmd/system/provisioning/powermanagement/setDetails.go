package powermanagement

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/provisioning/powermanagement"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	$param.getFlagName()          $param.getType()
	Name          string
}

func setDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setDetails",
		Short: "Get current power management settings of the given system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setDetails)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("Name", "", "")

	return cmd
}

func setDetails(globalFlags *types.GlobalFlags, flags *setDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := powermanagement.Powermanagement(&flags.ConnectionDetails, flags.Sid, flags.$param.getFlagName(), flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

