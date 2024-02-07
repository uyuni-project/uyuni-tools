package saltkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/saltkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deniedListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func deniedListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deniedList",
		Short: "List of denied salt keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deniedListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deniedList)
		},
	}


	return cmd
}

func deniedList(globalFlags *types.GlobalFlags, flags *deniedListFlags, cmd *cobra.Command, args []string) error {

res, err := saltkey.Saltkey(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

