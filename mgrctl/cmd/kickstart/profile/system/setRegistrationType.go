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

type setRegistrationTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func setRegistrationTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setRegistrationType",
		Short: "Sets the registration type of a given kickstart profile.
 Registration Type can be one of reactivation/deletion/none
 These types determine the behaviour of the re registration when using
 this profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setRegistrationTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setRegistrationType)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func setRegistrationType(globalFlags *types.GlobalFlags, flags *setRegistrationTypeFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

