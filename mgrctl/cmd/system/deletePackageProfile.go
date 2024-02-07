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

type deletePackageProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProfileId          int
}

func deletePackageProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deletePackageProfile",
		Short: "Delete a package profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deletePackageProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deletePackageProfile)
		},
	}

	cmd.Flags().String("ProfileId", "", "")

	return cmd
}

func deletePackageProfile(globalFlags *types.GlobalFlags, flags *deletePackageProfileFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.ProfileId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

