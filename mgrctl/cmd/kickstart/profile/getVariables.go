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

type getVariablesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getVariablesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getVariables",
		Short: "Returns a list of variables
                      associated with the specified kickstart profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getVariablesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getVariables)
		},
	}

	cmd.Flags().String("KsLabel", "", "")

	return cmd
}

func getVariables(globalFlags *types.GlobalFlags, flags *getVariablesFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

