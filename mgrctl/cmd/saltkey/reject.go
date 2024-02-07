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

type rejectFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MinionId              string
}

func rejectCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject",
		Short: "Reject a minion key",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags rejectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, reject)
		},
	}

	cmd.Flags().String("MinionId", "", "")

	return cmd
}

func reject(globalFlags *types.GlobalFlags, flags *rejectFlags, cmd *cobra.Command, args []string) error {

	res, err := saltkey.Saltkey(&flags.ConnectionDetails, flags.MinionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
