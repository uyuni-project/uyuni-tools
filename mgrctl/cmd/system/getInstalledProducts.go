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

type getInstalledProductsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	LoggedInUser          User
	Sid          int
}

func getInstalledProductsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getInstalledProducts",
		Short: "Get a list of installed products for given system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getInstalledProductsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getInstalledProducts)
		},
	}

	cmd.Flags().String("LoggedInUser", "", "")
	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getInstalledProducts(globalFlags *types.GlobalFlags, flags *getInstalledProductsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.LoggedInUser, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

