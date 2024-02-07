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

type getRegistrationTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getRegistrationTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRegistrationType",
		Short: "returns the registration type of a given kickstart profile.
 Registration Type can be one of reactivation/deletion/none
 These types determine the behaviour of the registration when using
 this profile for reprovisioning.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRegistrationTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRegistrationType)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func getRegistrationType(globalFlags *types.GlobalFlags, flags *getRegistrationTypeFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

