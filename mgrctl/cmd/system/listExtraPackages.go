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

type listExtraPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listExtraPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listExtraPackages",
		Short: "List extra packages for a system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listExtraPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listExtraPackages)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listExtraPackages(globalFlags *types.GlobalFlags, flags *listExtraPackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
