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

type setUpdateTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	UpdateType          string
}

func setUpdateTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setUpdateType",
		Short: "Set the update typefor a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setUpdateTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setUpdateType)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")
	cmd.Flags().String("UpdateType", "", "The new update type to set. Possible values are 'all' and 'none'.")

	return cmd
}

func setUpdateType(globalFlags *types.GlobalFlags, flags *setUpdateTypeFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.UpdateType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

