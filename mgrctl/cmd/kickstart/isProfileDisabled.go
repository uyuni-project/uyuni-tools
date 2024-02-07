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

type isProfileDisabledFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProfileLabel          string
}

func isProfileDisabledCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isProfileDisabled",
		Short: "Returns whether a kickstart profile is disabled",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isProfileDisabledFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isProfileDisabled)
		},
	}

	cmd.Flags().String("ProfileLabel", "", "kickstart profile label")

	return cmd
}

func isProfileDisabled(globalFlags *types.GlobalFlags, flags *isProfileDisabledFlags, cmd *cobra.Command, args []string) error {

	res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.ProfileLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
