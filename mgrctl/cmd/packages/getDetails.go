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

type getDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid                   int
}

func getDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDetails",
		Short: "Retrieve details for the package with the ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDetails)
		},
	}

	cmd.Flags().String("Pid", "", "")

	return cmd
}

func getDetails(globalFlags *types.GlobalFlags, flags *getDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := packages.Packages(&flags.ConnectionDetails, flags.Pid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
