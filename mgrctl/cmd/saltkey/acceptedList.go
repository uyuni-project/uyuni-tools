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

type acceptedListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func acceptedListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acceptedList",
		Short: "List accepted salt keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags acceptedListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, acceptedList)
		},
	}

	return cmd
}

func acceptedList(globalFlags *types.GlobalFlags, flags *acceptedListFlags, cmd *cobra.Command, args []string) error {

	res, err := saltkey.Saltkey(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
