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

type orderScriptsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func orderScriptsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orderScripts",
		Short: "Change the order that kickstart scripts will run for
 this kickstart profile. Scripts will run in the order they appear
 in the array. There are three arrays, one for all pre scripts, one
 for the post scripts that run before registration and server
 actions happen, and one for post scripts that run after registration
 and server actions. All scripts must be included in one of these
 lists, as appropriate.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags orderScriptsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, orderScripts)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart")

	return cmd
}

func orderScripts(globalFlags *types.GlobalFlags, flags *orderScriptsFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

