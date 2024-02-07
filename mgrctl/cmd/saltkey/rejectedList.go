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

type rejectedListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func rejectedListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rejectedList",
		Short: "List of rejected salt keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags rejectedListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, rejectedList)
		},
	}

	return cmd
}

func rejectedList(globalFlags *types.GlobalFlags, flags *rejectedListFlags, cmd *cobra.Command, args []string) error {

	res, err := saltkey.Saltkey(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
