package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type disableProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProfileLabel          string
	Disabled          string
}

func disableProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disableProfile",
		Short: "Enable/Disable a Kickstart Profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disableProfile)
		},
	}

	cmd.Flags().String("ProfileLabel", "", "Label for the kickstart tree you want to en/disable")
	cmd.Flags().String("Disabled", "", "true to disable the profile")

	return cmd
}

func disableProfile(globalFlags *types.GlobalFlags, flags *disableProfileFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.ProfileLabel, flags.Disabled)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

