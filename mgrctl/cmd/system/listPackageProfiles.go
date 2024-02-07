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

type listPackageProfilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listPackageProfilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackageProfiles",
		Short: "List the package profiles in this organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackageProfilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackageProfiles)
		},
	}

	return cmd
}

func listPackageProfiles(globalFlags *types.GlobalFlags, flags *listPackageProfilesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
