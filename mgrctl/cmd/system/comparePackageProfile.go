package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type comparePackageProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ProfileLabel          string
}

func comparePackageProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comparePackageProfile",
		Short: "Compare a system's packages against a package profile.  In
 the result returned, 'this_system' represents the server provided as an input
 and 'other_system' represents the profile provided as an input.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags comparePackageProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, comparePackageProfile)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("ProfileLabel", "", "")

	return cmd
}

func comparePackageProfile(globalFlags *types.GlobalFlags, flags *comparePackageProfileFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ProfileLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

