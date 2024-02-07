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

type createPackageProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ProfileLabel          string
	Description          string
}

func createPackageProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createPackageProfile",
		Short: "Create a new stored Package Profile from a systems
      installed package list.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createPackageProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createPackageProfile)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("ProfileLabel", "", "")
	cmd.Flags().String("Description", "", "")

	return cmd
}

func createPackageProfile(globalFlags *types.GlobalFlags, flags *createPackageProfileFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ProfileLabel, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

