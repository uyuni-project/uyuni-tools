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

type comparePackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KickstartLabel1          string
	KickstartLabel2          string
}

func comparePackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comparePackages",
		Short: "Returns a list for each kickstart profile; each list will contain
             package names not present on the other profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags comparePackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, comparePackages)
		},
	}

	cmd.Flags().String("KickstartLabel1", "", "")
	cmd.Flags().String("KickstartLabel2", "", "")

	return cmd
}

func comparePackages(globalFlags *types.GlobalFlags, flags *comparePackagesFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KickstartLabel1, flags.KickstartLabel2)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

