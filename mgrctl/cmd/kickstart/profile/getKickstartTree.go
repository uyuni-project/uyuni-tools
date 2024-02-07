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

type getKickstartTreeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getKickstartTreeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getKickstartTree",
		Short: "Get the kickstart tree for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getKickstartTreeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getKickstartTree)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")

	return cmd
}

func getKickstartTree(globalFlags *types.GlobalFlags, flags *getKickstartTreeFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

