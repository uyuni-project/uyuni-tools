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

type setKickstartTreeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	KstreeLabel          string
}

func setKickstartTreeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setKickstartTree",
		Short: "Set the kickstart tree for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setKickstartTreeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setKickstartTree)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")
	cmd.Flags().String("KstreeLabel", "", "Label of new kickstart tree.")

	return cmd
}

func setKickstartTree(globalFlags *types.GlobalFlags, flags *setKickstartTreeFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.KstreeLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

