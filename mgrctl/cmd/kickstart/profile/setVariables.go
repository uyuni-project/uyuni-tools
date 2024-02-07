package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setVariablesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	$param.getFlagName()          $param.getType()
}

func setVariablesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setVariables",
		Short: "Associates list of kickstart variables
                              with the specified kickstart profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setVariablesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setVariables)
		},
	}

	cmd.Flags().String("KsLabel", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setVariables(globalFlags *types.GlobalFlags, flags *setVariablesFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

