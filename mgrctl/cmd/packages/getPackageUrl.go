package packages

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getPackageUrlFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid          int
}

func getPackageUrlCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPackageUrl",
		Short: "Retrieve the url that can be used to download a package.
      This will expire after a certain time period.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPackageUrlFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPackageUrl)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func getPackageUrl(globalFlags *types.GlobalFlags, flags *getPackageUrlFlags, cmd *cobra.Command, args []string) error {

res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

